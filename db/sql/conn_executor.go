/*
 * Copyright (c) 2019 Ready Stock
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package sql

import (
	"context"
	"fmt"
	"github.com/readystock/golog"
	"github.com/readystock/noah/db/sql/driver/npgx"
	"github.com/readystock/noah/db/sql/pgwire/pgerror"
	"github.com/readystock/noah/db/sql/plan"
	"github.com/readystock/noah/db/sql/types"
	"github.com/readystock/noah/db/system"
	"github.com/readystock/noah/db/util"
	"github.com/readystock/noah/db/util/fsm"
	nodes "github.com/readystock/pg_query_go/nodes"
	"io"
	"sync"
)

type Server struct {
}

type connExecutor struct {
	ClientAddress      string
	server             *Server
	stmtBuf            *StmtBuf
	Backlog            []string
	clientComm         ClientComm
	prepStmtsNamespace prepStmtNamespace
	curStmt            *nodes.Stmt
	SystemContext      *system.SContext

	nSync            sync.Mutex
	nodes            map[uint64]*npgx.Conn
	nodeTransactions map[uint64]bool
	TransactionState TransactionState
	TransactionMode  TransactionMode
	TransactionID    uint64

	typeInfo *types.ConnInfo
	changes  uint64
}

func (ex *connExecutor) GetNodesForAccountID(id *uint64) ([]uint64, error) {
	return []uint64{1, 2}, nil
}

func (ex *connExecutor) BacklogQuery(query string) {
	ex.Backlog = append(ex.Backlog, query)
}

func (ex *connExecutor) GetNodeConnection(nodeId uint64, doTransaction bool) (*npgx.Conn, error) {
	ex.nSync.Lock()
	defer ex.nSync.Unlock()
	if ex.nodes == nil {
		ex.nodes = map[uint64]*npgx.Conn{}
	}
	if ex.nodeTransactions == nil {
		ex.nodeTransactions = map[uint64]bool{}
	}

	conn, ok := ex.nodes[nodeId]
	if !ok {
		golog.Debugf("node [%d] is not in the session, acquiring connection", nodeId)
		// A connection has not yet been made to this node. Allocate one.
		if t, err := ex.SystemContext.Pool.Acquire(nodeId); err != nil {
			golog.Errorf(err.Error())
			return nil, err
		} else {
			ex.nodes[nodeId] = t
			golog.Warnf("new connection count is %d node(s) from executor", len(ex.nodes))
			conn = t
		}
	}

	nodeInTransaction, ok := ex.nodeTransactions[nodeId]
	if !ok {
		nodeInTransaction = false
	}

	// We have just acquired a new connection to the node. If we are currently in a transaction state
	// then we will want to begin a transaction on this node.
	if !nodeInTransaction && (doTransaction || ex.TransactionState == TransactionState_PRE || ex.TransactionState == TransactionState_ENTERED) {
		result, err := conn.Query("BEGIN")
		if err != nil {
			return nil, err
		}

		result.Next()

		if result.Err() != nil {
			return nil, result.Err()
		}
		ex.TransactionState = TransactionState_ENTERED
		ex.nodeTransactions[nodeId] = true

		if ex.TransactionID == 0 {
			transactionId, err := ex.SystemContext.NewSnowflake()
			if err != nil {
				return nil, err
			}
			ex.TransactionID = transactionId
		}
	}

	return conn, nil
}

func (ex *connExecutor) ReleaseNodeConnection(nodeId uint64) error {
	ex.nSync.Lock()
	defer ex.nSync.Unlock()
	if ex.nodes == nil {
		return nil
	}

	conn, ok := ex.nodes[nodeId]
	if !ok {
		return nil
	}

	delete(ex.nodes, nodeId)
	delete(ex.nodeTransactions, nodeId)
	return ex.SystemContext.Pool.Release(conn)
}

func (ex *connExecutor) ReleaseAllConnections() error {
	errs := make([]error, 0)
	// return any connections to pool.
	golog.Infof("releasing connection to %d node(s) from executor", len(ex.nodes))
	for nodeId := range ex.nodes {
		golog.Verbosef("releasing connection to node [%d] from executor", nodeId)
		if err := ex.ReleaseNodeConnection(nodeId); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return util.CombineErrors(errs)
	}

	return nil
}

func (ex *connExecutor) SetNodeTransaction(nodeId uint64, conn *npgx.Conn) {
	ex.nSync.Lock()
	defer ex.nSync.Unlock()
	if ex.nodes == nil {
		ex.nodes = map[uint64]*npgx.Conn{}
	}
	ex.nodes[nodeId] = conn
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
	return &Server{}
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
		TransactionState: TransactionState_NONE,
		TransactionMode:  TransactionMode_AutoCommit,
		typeInfo:         types.NewConnInfo(),
	}
	ex.typeInfo.InitializeDataTypes(npgx.NameOIDs)
	return ex
}

func (ex *connExecutor) run() (err error) {
	defer util.CatchPanic(&err)
	for {
		cmd, pos, err := ex.stmtBuf.curCmd()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if cmd == nil {
			golog.Warnf("found null command, advancing 1")
			ex.stmtBuf.advanceOne()
		}
		// var ev fsm.Event
		var res ResultBase
		var payload fsm.EventPayload

		switch tcmd := cmd.(type) {
		case ExecStmt:
			if tcmd.Stmt == nil {
				res = ex.clientComm.CreateEmptyQueryResult(pos)
				break
			}
			ex.curStmt = &tcmd.Stmt

			stmtRes := ex.clientComm.CreateStatementResult(tcmd.Stmt, NeedRowDesc, pos, nil /* formatCodes */)
			res = stmtRes

			// ex.phaseTimes[sessionQueryReceived] = tcmd.TimeReceived
			// ex.phaseTimes[sessionStartParse] = tcmd.ParseStart
			// ex.phaseTimes[sessionEndParse] = tcmd.ParseEnd

			err = ex.execStmt(*ex.curStmt, stmtRes, pos, nil)
		case ExecPortal:
			portal, ok := ex.prepStmtsNamespace.portals[tcmd.Name]
			if !ok {
				err = pgerror.NewErrorf(pgerror.CodeInvalidCursorNameError, "unknown portal %q", tcmd.Name)
				payload = eventNonRetriableErrPayload{err: err}
				res = ex.clientComm.CreateErrorResult(pos)
				break
			}
			ex.curStmt = portal.Stmt.Statement

			pinfo := &plan.PlaceholderInfo{
				TypeHints: portal.Stmt.TypeHints,
				Types:     portal.Stmt.Types,
				Values:    portal.Qargs,
			}

			if portal.Stmt.Statement == nil {
				res = ex.clientComm.CreateEmptyQueryResult(pos)
				break
			}

			stmtRes := ex.clientComm.CreateStatementResult(
				*portal.Stmt.Statement,
				// The client is using the extended protocol, so no row description is
				// needed.
				NeedRowDesc, pos, portal.OutFormats)
			res = stmtRes
			err = ex.execStmt(*ex.curStmt, stmtRes, pos, pinfo)
		case PrepareStmt:
			res = ex.clientComm.CreatePrepareResult(pos)
			err = ex.execPrepare(tcmd)
		case DescribeStmt:
			descRes := ex.clientComm.CreateDescribeResult(pos)
			res = descRes
			err = ex.execDescribe(context.Background(), tcmd, descRes)
		case BindStmt:
			res = ex.clientComm.CreateBindResult(pos)
			err = ex.execBind(tcmd)
		case DeletePreparedStmt:
			res = ex.clientComm.CreateDeleteResult(pos)
		case SendError:
			res = ex.clientComm.CreateErrorResult(pos)
			// ev = eventNonRetriableErr{IsCommit: fsm.False}
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
			// TODO (elliotcourant) Add proper draining handling in the connection executor.
		case Flush:
			// Closing the res will flush the connection's buffer.
			res = ex.clientComm.CreateFlushResult(pos)
		default:
			panic(fmt.Sprintf("unsupported command type: %T", cmd))
		}

		// if ev != nil {
		//     var err error
		//     advInfo, err = ex.txnStateTransitionsApplyWrapper(ev, payload, res, pos)
		//     if err != nil {
		//         return err
		//     }
		// } else {
		//     // If no event was generated synthesize an advance code.
		//     advInfo = advanceInfo{
		//         code: advanceOne,
		//     }
		// }

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

// stmtHasNoData returns true if describing a result of the input statement
// type should return NoData.
func stmtHasNoData(stmt nodes.Stmt) bool {
	return stmt == nil || stmt.StatementType() != nodes.Rows
}
