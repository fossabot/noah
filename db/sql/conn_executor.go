package sql

import (
	"fmt"
	"github.com/Ready-Stock/Noah/db/sql/pgwire/pgerror"
	"github.com/Ready-Stock/Noah/db/system"
	"github.com/Ready-Stock/Noah/db/util/fsm"
	nodes "github.com/Ready-Stock/pg_query_go/nodes"
	"github.com/Ready-Stock/pgx"
	"io"
)

type Server struct {
}

type connExecutor struct {
	ClientAddress	   string
	server             *Server
	stmtBuf            *StmtBuf
	Backlog			   []string
	clientComm         ClientComm
	prepStmtsNamespace prepStmtNamespace
	curStmt            *nodes.Stmt
	SystemContext      *system.SContext
	Nodes              map[uint64]*pgx.Tx
}

func (ex *connExecutor) GetNodesForAccountID(id *uint64) ([]uint64, error) {
	return []uint64{1, 2}, nil
}

func (ex *connExecutor) BacklogQuery(query string) {
	ex.Backlog = append(ex.Backlog, query)
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

func (s *Server) ServeConn(stmtBuf *StmtBuf, clientComm ClientComm, clientAddress string, sctx *system.SContext) error {
	ex := s.newConnExecutor(stmtBuf, clientComm)
	ex.ClientAddress = clientAddress
	ex.SystemContext = sctx
	return ex.run()
}

func (s *Server) newConnExecutor(stmtBuf *StmtBuf, clientComm ClientComm) *connExecutor {
	ex := &connExecutor{
		server:     s,
		stmtBuf:    stmtBuf,
		clientComm: clientComm,
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
		var res ResultBase
		var payload fsm.EventPayload
		switch tcmd := cmd.(type) {
		case ExecStmt:
		case ExecPortal:
			portal, ok := ex.prepStmtsNamespace.portals[tcmd.Name]
			if !ok {
				err = pgerror.NewErrorf(pgerror.CodeInvalidCursorNameError, "unknown portal %q", tcmd.Name)
				payload = eventNonRetriableErrPayload{err: err}
				res = ex.clientComm.CreateErrorResult(pos)
				break
			}
			// fmt.Printf("portal resolved to: %s", portal.Stmt.Str)
			ex.curStmt = portal.Stmt.Statement

			if portal.Stmt.Statement == nil {
				res = ex.clientComm.CreateEmptyQueryResult(pos)
				break
			}

			stmtRes := ex.clientComm.CreateStatementResult(
				*portal.Stmt.Statement,
				// The client is using the extended protocol, so no row description is
				// needed.
				DontNeedRowDesc, pos, portal.OutFormats)
			res = stmtRes
			err = ex.execStmt(*ex.curStmt, stmtRes, pos)
		case PrepareStmt:
			res = ex.clientComm.CreatePrepareResult(pos)
			err = ex.execPrepare(tcmd)
		case DescribeStmt:
			descRes := ex.clientComm.CreateDescribeResult(pos)
			res = descRes
		case BindStmt:
			res = ex.clientComm.CreateBindResult(pos)
			err = ex.execBind(tcmd)
		case DeletePreparedStmt:
			res = ex.clientComm.CreateDeleteResult(pos)
		case SendError:
			res = ex.clientComm.CreateErrorResult(pos)
			payload = eventNonRetriableErrPayload{err: tcmd.Err}
		case Sync:
			// Note that the Sync result will flush results to the network connection.
			res = ex.clientComm.CreateSyncResult(pos)
		case CopyIn:
			res = ex.clientComm.CreateCopyInResult(pos)
		case DrainRequest:
			// We received a drain request. We terminate immediately if we're not in a
			// transaction. If we are in a transaction, we'll finish as soon as a Sync
			// command (i.e. the end of a batch) is processed outside of a
			// transaction.
		case Flush:
			// Closing the res will flush the connection's buffer.
			res = ex.clientComm.CreateFlushResult(pos)
		default:
			panic(fmt.Sprintf("unsupported command type: %T", cmd))
		}

		if err != nil {
			res.CloseWithErr(err)
			if err := ex.stmtBuf.seekToNextBatch(); err != nil {
				return err
			}
		} else {
			resErr := res.Err()
			pe, ok := payload.(payloadWithError)
			if !ok {
				// error stuff here
			}
			if resErr == nil && ok {
				// ex.server.recordError(pe.errorCause())
				// Depending on whether the result has the error already or not, we have
				// to call either Close or CloseWithErr.
				res.CloseWithErr(pe.errorCause())
				if err := ex.stmtBuf.seekToNextBatch(); err != nil {
					return err
				}
			} else {
				// ex.server.recordError(resErr)
				res.Close(IdleTxnBlock)
				ex.stmtBuf.advanceOne()
			}
		}
	}
}


