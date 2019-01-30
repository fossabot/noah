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
 */

package sql

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/juju/errors"
	"github.com/readystock/golog"
	"github.com/readystock/noah/db/sql/driver/npgx"
	"github.com/readystock/noah/db/sql/pgwire/pgproto"
	"github.com/readystock/noah/db/sql/plan"
	"github.com/readystock/noah/db/sql/types"
	"github.com/readystock/noah/db/util"
	"github.com/readystock/pg_query_go/nodes"
)

type executeResponse struct {
	Error  error
	Rows   *npgx.Rows
	NodeID uint64
	Type   pg_query.StmtType
}

func (ex *connExecutor) ExecutePlans(plans []plan.NodeExecutionPlan, res RestrictedCommandResult) error {
	errs := make([]error, 0)
	defer func() {
		if ex.TransactionState == TransactionState_NONE {
			golog.Infof("done executing non-transaction statement, releasing connections")
			if err := ex.ReleaseAllConnections(); err != nil {
				errs = append(errs, err)
			}
		}
	}()

	// defer util.CatchPanic(&err)
	if len(plans) == 0 {
		golog.Warnf("no plans were provided, nothing will be executed")
		return errors.New("no plans were provided")
	}

	// If none of the plans are read only then increment the number of changes at the end of the
	// execution.
	defer func(plans []plan.NodeExecutionPlan) {
		if linq.From(plans).AllT(func(plan plan.NodeExecutionPlan) bool {
			return !plan.ReadOnly
		}) {
			ex.changes++
		}
	}(plans)

	responses := make(chan *executeResponse, len(plans))
	for _, p := range plans {
		go func(ex *connExecutor, pln plan.NodeExecutionPlan) {
			// defer util.CatchPanic(&err)
			exResponse := executeResponse{
				NodeID: pln.Node.NodeId,
				Type:   pln.Type,
			}

			defer func() {
				responses <- &exResponse
			}()

			golog.Infof("executing query: `%s` on database node [%d] | transaction: %v", pln.CompiledQuery, pln.Node.NodeId, !pln.ReadOnly && len(plans) > 1)
			conn, err := ex.GetNodeConnection(pln.Node.NodeId, !pln.ReadOnly && len(plans) > 1)
			if err != nil {
				golog.Errorf(err.Error())
				exResponse.Error = err
				return
			}

			rows, err := conn.Query(pln.CompiledQuery)
			if err != nil {
				golog.Errorf(err.Error())
				exResponse.Error = err
				return
			}

			golog.Debugf("received rows response from node [%d]", pln.Node.NodeId)
			exResponse.Rows = rows
		}(ex, p)
	}

	columns := make([]pgproto.FieldDescription, 0)
	result := make([][]types.Value, 0)
	for i := 0; i < len(plans); i++ {
		response := <-responses
		golog.Debugf("handling response from node [%d]", response.NodeID)
		if response.Error != nil {
			golog.Errorf("received error from node [%d]: %s", response.NodeID, response.Error.Error())
			errs = append(errs, response.Error)
			continue
		}

		switch response.Type {
		case pg_query.Rows:
			if response.Rows != nil {
				for response.Rows.Next() {
					if len(columns) == 0 {
						columns = response.Rows.PgFieldDescriptions()
						res.SetColumns(columns)
					}

					row := make([]types.Value, len(columns))
					if values, err := response.Rows.PgValues(); err != nil {
						golog.Errorf("reading values from wire: %v", err.Error())
						errs = append(errs, err)
						continue
					} else {
						for c, v := range values {
							row[c] = v
						}
					}

					res.AddRow(row)
					result = append(result, row)
				}

				if err := response.Rows.Err(); err != nil {
					golog.Errorf("received error from node [%d]: %s", response.NodeID, err)
					errs = append(errs, err)
				}

				if len(columns) == 0 {
					columns = response.Rows.PgFieldDescriptions()
					res.SetColumns(columns)
				}

				response.Rows.Close()
			} else {
				golog.Debugf("no rows returned for query `%s`", plans[0].CompiledQuery)
			}
		case pg_query.DDL, pg_query.RowsAffected:
			if response.Rows != nil {
				response.Rows.Next()

				if err := response.Rows.Err(); err != nil {
					golog.Errorf("received error from node [%d]: %s", response.NodeID, err)
					errs = append(errs, err)
				}

				response.Rows.Close()
			}
		default:
			return errors.Errorf("cannot handle statement type [%s]", response.Type.String())
		}
	}
	golog.Debugf("%d row(s) compiled for query `%s`", len(result), plans[0].CompiledQuery)

	switch ex.TransactionMode {
	case TransactionMode_AutoCommit:
		// Transaction handling.
		switch ex.TransactionState {
		case TransactionState_NONE, TransactionState_ENTERED:
			// If we are not in a transaction then we want to treat each write statement
			// as a transaction.
			if linq.From(plans).AnyWithT(func(plan plan.NodeExecutionPlan) bool {
				return !plan.ReadOnly
			}) {
				// If there were writes in the current query plan.
				if len(errs) == 0 {
					// If we are targeting more than 1 node then we want to handle that.
					if len(ex.nodes) > 1 {
						if err := ex.PrepareTwoPhase(); err != nil {
							return util.CombineErrors(append(errs, err))
						}

						if err := ex.CommitTwoPhase(); err != nil {
							return util.CombineErrors(append(errs, err))
						}
						return util.CombineErrors(errs)
					} else {
						return util.CombineErrors(errs)
					}
				}
			}
		case TransactionState_PRE:
			// Technically this should not be able to happen.
		case TransactionState_ENDING:

		default:

		}
	case TransactionMode_Manual:

	default:

	}
	return util.CombineErrors(errs)
}

func (ex *connExecutor) PrepareTwoPhase() error {
	if ex.nodes == nil {
		ex.nodes = map[uint64]*npgx.Conn{}
		return nil // There were never any transactions to release
	}

	golog.Debugf("preparing two-phase transaction for changes on %d node(s), transaction ID [%d]", len(ex.nodes), ex.TransactionID)

	if err := ex.doTransaction(func(nodeID uint64) string {
		return fmt.Sprintf(`PREPARE TRANSACTION '%d%d'`, nodeID, ex.TransactionID)
	}); err != nil {
		return err
	}
	return nil
}

func (ex *connExecutor) CommitTwoPhase() error {
	defer func() {
		transactionId, _ := ex.SystemContext.NewSnowflake()
		ex.TransactionID = transactionId
	}()
	if ex.nodes == nil {
		ex.nodes = map[uint64]*npgx.Conn{}
		return nil // There were never any transactions to release
	}

	golog.Debugf("preparing two-phase commit for changes on %d node(s), transaction ID [%d]", len(ex.nodes), ex.TransactionID)
	ex.TransactionState = TransactionState_NONE
	return ex.doTransaction(func(nodeID uint64) string {
		return fmt.Sprintf(`COMMIT PREPARED '%d%d'`, nodeID, ex.TransactionID)
	})
}

func (ex *connExecutor) RollbackTwoPhase() error {
	defer func() {
		transactionId, _ := ex.SystemContext.NewSnowflake()
		ex.TransactionID = transactionId
	}()
	if ex.nodes == nil {
		ex.nodes = map[uint64]*npgx.Conn{}
		return nil // There were never any transactions to release
	}

	golog.Debugf("preparing two-phase rollback for changes on %d node(s), transaction ID [%d]", len(ex.nodes), ex.TransactionID)
	ex.TransactionState = TransactionState_NONE
	return ex.doTransaction(func(nodeID uint64) string {
		return fmt.Sprintf(`ROLLBACK PREPARED '%d%d'`, nodeID, ex.TransactionID)
	})
}

func (ex *connExecutor) Commit() error {
	defer func() {
		transactionId, _ := ex.SystemContext.NewSnowflake()
		ex.TransactionID = transactionId
	}()
	if ex.nodes == nil {
		ex.nodes = map[uint64]*npgx.Conn{}
		return nil // There were never any transactions to release
	}

	golog.Debugf("committing changes on %d node(s), transaction ID [%d]", len(ex.nodes), ex.TransactionID)
	ex.TransactionState = TransactionState_NONE
	return ex.doTransaction(func(nodeID uint64) string {
		return "COMMIT"
	})
}

func (ex *connExecutor) Rollback() error {
	defer func() {
		transactionId, _ := ex.SystemContext.NewSnowflake()
		ex.TransactionID = transactionId
	}()
	if ex.nodes == nil {
		ex.nodes = map[uint64]*npgx.Conn{}
		return nil // There were never any transactions to release
	}

	golog.Debugf("rolling back changes on %d node(s), transaction ID [%d]", len(ex.nodes), ex.TransactionID)
	ex.TransactionState = TransactionState_NONE
	return ex.doTransaction(func(nodeID uint64) string {
		return "ROLLBACK"
	})
}

func (ex *connExecutor) doTransaction(query func(nodeID uint64) string) error {
	ex.nSync.Lock()
	defer ex.nSync.Unlock()
	if ex.nodes == nil {
		ex.nodes = map[uint64]*npgx.Conn{}
		return nil // There were never any transactions to release
	}

	if ex.nodeTransactions == nil {
		ex.nodeTransactions = map[uint64]bool{}
	}

	// golog.Debugf("preparing two-phase commit for changes on %d node(s), transaction ID [%d]", len(ex.nodes), transactionID)

	count := len(ex.nodes)
	responses := make(chan executeResponse, count)
	for id, conn := range ex.nodes {
		go func(conn *npgx.Conn, nodeId uint64) {
			response := executeResponse{
				NodeID: nodeId,
			}
			defer func() { responses <- response }()

			nodeInTransaction, ok := ex.nodeTransactions[nodeId]
			if !ok {
				nodeInTransaction = false
			}

			// if the node is not in a transaction then don't send any transaction nodes
			if !nodeInTransaction {
				return
			}

			qry := query(nodeId)
			golog.Verbosef("issuing transaction query `%s` to node [%d]", qry, nodeId)
			result, err := conn.Query(qry)
			if err != nil {
				response.Error = err
				return
			}

			result.Next()

			if err := result.Err(); err != nil {
				response.Error = err
			}
		}(conn, id)
	}

	errs := make([]error, 0)
	for i := 0; i < count; i++ {
		response := <-responses
		golog.Verbosef("finished issuing transaction to node [%d]", response.NodeID)
		if response.Error != nil {
			golog.Errorf("error from transaction on node [%d]: %s", response.NodeID, response.Error)
			errs = append(errs, response.Error)
		} else {
			golog.Verbosef("transaction on node [%d] was successful", response.NodeID)
		}
	}

	golog.Verbosef("%d error(s) from execution", len(errs))
	return util.CombineErrors(errs)
}
