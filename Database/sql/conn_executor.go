package sql

import (
	"io"
	"fmt"
	"github.com/Ready-Stock/Noah/Database/sql/pgwire/pgerror"
	"github.com/Ready-Stock/Noah/Database/sql/sem/tree"
	"github.com/Ready-Stock/pg_query_go"
	"github.com/Ready-Stock/Noah/Database/sql/context"
)

type Server struct {
}

type connExecutor struct {
	server             *Server
	stmtBuf            *StmtBuf
	context			   context.Context
	prepStmtsNamespace prepStmtNamespace
	curStmt            *pg_query.ParsetreeList
}

type prepStmtNamespace struct {
	// prepStmts contains the prepared statements currently available on the
	// session.
	prepStmts map[string]prepStmtEntry
	// portals contains the portals currently available on the session.
	portals map[string]portalEntry
}

type prepStmtEntry struct {
	*PreparedStatement
	portals map[string]struct{}
}

func (pe *prepStmtEntry) copy() prepStmtEntry {
	cpy := prepStmtEntry{}
	cpy.PreparedStatement = pe.PreparedStatement
	cpy.portals = make(map[string]struct{})
	for pname := range pe.portals {
		cpy.portals[pname] = struct{}{}
	}
	return cpy
}

type portalEntry struct {
	*PreparedPortal
	psName string
}

// resetTo resets a namespace to equate another one (`to`). Prep stmts and portals
// that are present in ns but not in to are deallocated.
//
// A (pointer to) empty `to` can be passed in to deallocate everything.
func (ns *prepStmtNamespace) resetTo(to *prepStmtNamespace) {
	for name, ps := range ns.prepStmts {
		bps, ok := to.prepStmts[name]
		// If the prepared statement didn't exist before (including if a statement
		// with the same name existed, but it was different), close it.
		if !ok || bps.PreparedStatement != ps.PreparedStatement {
			ps.close()
		}
	}
	for name, p := range ns.portals {
		bp, ok := to.portals[name]
		// If the prepared statement didn't exist before (including if a statement
		// with the same name existed, but it was different), close it.
		if !ok || bp.PreparedPortal != p.PreparedPortal {
			p.close()
		}
	}
	*ns = to.copy()
}

func (ns *prepStmtNamespace) copy() prepStmtNamespace {
	var cpy prepStmtNamespace
	cpy.prepStmts = make(map[string]prepStmtEntry)
	for name, psEntry := range ns.prepStmts {
		cpy.prepStmts[name] = psEntry.copy()
	}
	cpy.portals = make(map[string]portalEntry)
	for name, p := range ns.portals {
		cpy.portals[name] = p
	}
	return cpy
}

func NewServer() *Server {
	return &Server{

	}
}

func (s *Server) ServeConn(stmtBuf *StmtBuf, ctx context.Context) error {
	ex := s.newConnExecutor(stmtBuf, ctx)
	return ex.run()
}

func (s *Server) newConnExecutor(stmtBuf *StmtBuf, ctx context.Context) *connExecutor {
	ex := &connExecutor{
		server:     s,
		stmtBuf:    stmtBuf,
		context: 	ctx,
		prepStmtsNamespace: prepStmtNamespace{
			prepStmts: make(map[string]prepStmtEntry),
			portals:   make(map[string]portalEntry),
		},
	}
	return ex
}

func (ex *connExecutor) run() error {
	for {
		cmd, pos, err := ex.stmtBuf.curCmd()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		fmt.Printf("[pos:%d] executing %s \n", pos, cmd)
		var res ResultBase
		switch tcmd := cmd.(type) {
		case ExecStmt:
			fmt.Println("TYPE: ExecStmt")
		case ExecPortal:
			// ExecPortal is handled like ExecStmt, except that the placeholder info
			// is taken from the portal.
			fmt.Println("TYPE: ExecPortal")
			portal, ok := ex.prepStmtsNamespace.portals[tcmd.Name]
			if !ok {
				err = pgerror.NewErrorf(pgerror.CodeInvalidCursorNameError, "unknown portal %q", tcmd.Name)
				// ev = eventNonRetriableErr{IsCommit: fsm.False}
				// payload = eventNonRetriableErrPayload{err: err}
				res = ex.context.ClientComm.CreateErrorResult(pos)
				break
			}
			ex.curStmt = portal.Stmt.Statement

			pinfo := &tree.PlaceholderInfo{
				TypeHints: portal.Stmt.TypeHints,
				Types:     portal.Stmt.Types,
				Values:    portal.Qargs,
			}

			if portal.Stmt.Statement == nil {
				res = ex.context.ClientComm.CreateEmptyQueryResult(pos)
				break
			}

			stmtRes := ex.context.ClientComm.CreateStatementResult(
				*portal.Stmt.Statement,
				// The client is using the extended protocol, so no row description is
				// needed.
				DontNeedRowDesc,
				pos, portal.OutFormats,
				ex.sessionData.Location, ex.sessionData.BytesEncodeFormat)
			stmtRes.SetLimit(tcmd.Limit)

		case PrepareStmt:
			fmt.Println("TYPE: PrepareStmt")
			fmt.Println("Len:", tcmd.PGQuery.Query)
			res = ex.context.ClientComm.CreatePrepareResult(pos)
		case DescribeStmt:
			fmt.Println("TYPE: DescribeStmt")
		case BindStmt:
			fmt.Println("TYPE: BindStmt")
			res = ex.context.ClientComm.CreateBindResult(pos)
		case DeletePreparedStmt:
			fmt.Println("TYPE: DeletePreparedStmt")
			res = ex.context.ClientComm.CreateDeleteResult(pos)
		case SendError:
			fmt.Println("TYPE: SendError")
			res = ex.context.ClientComm.CreateErrorResult(pos)
		case Sync:
			fmt.Println("TYPE: Sync")
			// Note that the Sync result will flush results to the network connection.
			res = ex.context.ClientComm.CreateSyncResult(pos)
		case CopyIn:
			fmt.Println("TYPE: CopyIn")
			res = ex.context.ClientComm.CreateCopyInResult(pos)
		case DrainRequest:
			fmt.Println("TYPE: DrainRequest")
			// We received a drain request. We terminate immediately if we're not in a
			// transaction. If we are in a transaction, we'll finish as soon as a Sync
			// command (i.e. the end of a batch) is processed outside of a
			// transaction.

		case Flush:
			fmt.Println("TYPE: Flush")
			// Closing the res will flush the connection's buffer.
			res = ex.context.ClientComm.CreateFlushResult(pos)
		default:
			panic(fmt.Sprintf("unsupported command type: %T", cmd))
		}

		ex.stmtBuf.advanceOne()
		if res != nil {
			res.Close(IdleTxnBlock)
		}
	}
}
