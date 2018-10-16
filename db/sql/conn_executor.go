/*
 * Copyright (c) 2018 Ready Stock
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
 *
 * This application uses Open Source components. You can find the
 * source code of their open source projects along with license
 * information below. We acknowledge and are grateful to these
 * developers for their contributions to open source.
 *
 * Project: CockroachDB https://github.com/cockroachdb/cockroach
 * Copyright 2018 The Cockroach Authors.
 * License (Apache License 2.0) https://github.com/cockroachdb/cockroach/blob/master/LICENSE
 *
 * Project: Vitess https://github.com/vitessio/vitess
 * Copyright 2018 Google Inc.
 * License (Apache License 2.0) https://github.com/vitessio/vitess/blob/master/LICENSE
 *
 * Project: Citus https://github.com/citusdata/citus
 * Copyright 2018 Citus Data, Inc.
 * License (GNU Affero General Public License v3.0) https://github.com/citusdata/citus/blob/master/LICENSE
 *
 * Project: pg_query_go https://github.com/lfittl/pg_query_go
 * Copyright 2018 Lukas Fittl
 * License (3-Clause BSD) https://github.com/lfittl/pg_query_go/blob/master/LICENSE
 *
 * Project: pgx https://github.com/jackc/pgx
 * Copyright 2018 Jack Christensen
 * License (MIT) https://github.com/jackc/pgx/blob/master/LICENSE
 *
 * Project: BadgerDB https://github.com/dgraph-io/badger
 * Copyright 2018 Dgraph Labs, Inc. and Contributors
 * License (MIT) https://github.com/dgraph-io/badger/blob/master/LICENSE
 *
 * Project: Sonyflake https://github.com/sony/sonyflake
 * Copyright 2018 Sony Corporation
 * License (MIT) https://github.com/sony/sonyflake/blob/master/LICENSE
 *
 * Project: Raft https://github.com/hashicorp/raft
 * Copyright 2018 HashiCorp
 * License (MPL-2.0) https://github.com/hashicorp/raft/blob/master/LICENSE
 */

package sql

import (
	"fmt"
	"github.com/Ready-Stock/Noah/db/sql/driver/npgx"
	"github.com/Ready-Stock/Noah/db/sql/pgwire/pgerror"
	"github.com/Ready-Stock/Noah/db/system"
	"github.com/Ready-Stock/Noah/db/util"
	"github.com/Ready-Stock/Noah/db/util/fsm"
	nodes "github.com/Ready-Stock/pg_query_go/nodes"
	"sync"

	// "github.com/Ready-Stock/pgx"
	"io"
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

	nSync             sync.Mutex
	nodes             map[uint64]*npgx.Transaction
	TransactionStatus NTXStatus
	TransactionID     uint64
}

func (ex *connExecutor) GetNodesForAccountID(id *uint64) ([]uint64, error) {
	return []uint64{1, 2}, nil
}

func (ex *connExecutor) BacklogQuery(query string) {
	ex.Backlog = append(ex.Backlog, query)
}

func (ex *connExecutor) GetNodeTransaction(nodeId uint64) (*npgx.Transaction, bool) {
	ex.nSync.Lock()
	defer ex.nSync.Unlock()
	tx, ok := ex.nodes[nodeId]
	return tx, ok
}

func (ex *connExecutor) SetNodeTransaction(nodeId uint64, tx *npgx.Transaction) {
	ex.nSync.Lock()
	defer ex.nSync.Unlock()
	if ex.nodes == nil {
		ex.nodes = map[uint64]*npgx.Transaction{}
	}
	ex.nodes[nodeId] = tx
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
		TransactionStatus: NTXNoTransaction,
	}
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
			ex.Warn("found null command, advancing 1")
			ex.stmtBuf.advanceOne()
		}
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

			err = ex.execStmt(*ex.curStmt, stmtRes, pos)
			if err != nil {
				return err
			}
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
				NeedRowDesc, pos, portal.OutFormats)
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
			// TODO (elliotcourant) Add proper draining handling in the connection executor.
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
