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
			golog.Errorf("done executing non-transaction statement, releasing connections")
			if err := ex.ReleaseAllConnections(); err != nil {
				errs = append(errs, err)
			}
		}
	}()

	// defer util.CatchPanic(&err)
	if len(plans) == 0 {
		golog.Errorf("no plans were provided, nothing will be executed")
		return errors.New("no plans were provided")
	}

	// If none of the plans are read only then increment the number of changes at the end of the
	// execution.
	if linq.From(plans).AllT(func(plan plan.NodeExecutionPlan) bool {
		return !plan.ReadOnly
	}) {
		defer func() {
			ex.changes++
		}()
	}

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

			golog.Infof("executing query: `%s` on database node [%d]", pln.CompiledQuery, pln.Node.NodeId)
			conn, err := ex.GetNodeConnection(pln.Node.NodeId)
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

					golog.Debugf("retrieved %d column(s)", len(columns))

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

				response.Rows.Close()
			} else {
				golog.Debugf("no rows returned for query `%s`", plans[0].CompiledQuery)
			}
		case pg_query.DDL:
			if response.Rows != nil {
				response.Rows.Next()

				if err := response.Rows.Err(); err != nil {
					golog.Errorf("received error from node [%d]: %s", response.NodeID, err)
					errs = append(errs, err)
				}

				response.Rows.Close()
			}
		default:
			return errors.Errorf("cannot handle statement type %d", response.Type)
		}
	}
	golog.Debugf("%d row(s) compiled for query `%s`", len(result), plans[0].CompiledQuery)

	if ex.TransactionMode == TransactionMode_AutoCommit {
		if linq.From(plans).AnyWithT(func(plan plan.NodeExecutionPlan) bool {
			return !plan.ReadOnly
		}) {
			// If there were writes in the current query plan.
			if len(errs) == 0 {
				// If we are targeting more than 1 node then we want to handle that.
				if len(ex.nodes) > 1 {
					if err := ex.Commit(); err != nil {
						return util.CombineErrors(append(errs, err))
					}
				}
			} else {
				return util.CombineErrors(errs)
			}
		}
	}
	return util.CombineErrors(errs)
}

func (ex *connExecutor) PrepareTwoPhase() error {
	if ex.nodes == nil {
		ex.nodes = map[uint64]*npgx.Conn{}
		return nil // There were never any transactions to release
	}

	golog.Debugf("preparing two-phase transaction for changes on %d node(s), transaction ID [%d]", len(ex.nodes), ex.TransactionID)

	return ex.doTransaction(func(nodeID uint64) string {
		return fmt.Sprintf("PREPARE TRANSACTION %d_%d", nodeID, ex.TransactionID)
	})
}

func (ex *connExecutor) CommitTwoPhase() error {
	if ex.nodes == nil {
		ex.nodes = map[uint64]*npgx.Conn{}
		return nil // There were never any transactions to release
	}

	golog.Debugf("preparing two-phase commit for changes on %d node(s), transaction ID [%d]", len(ex.nodes), ex.TransactionID)

	return ex.doTransaction(func(nodeID uint64) string {
		return fmt.Sprintf("COMMIT TRANSACTION %d_%d", nodeID, ex.TransactionID)
	})
}

func (ex *connExecutor) RollbackTwoPhase() error {
	if ex.nodes == nil {
		ex.nodes = map[uint64]*npgx.Conn{}
		return nil // There were never any transactions to release
	}

	golog.Debugf("preparing two-phase rollback for changes on %d node(s), transaction ID [%d]", len(ex.nodes), ex.TransactionID)

	return ex.doTransaction(func(nodeID uint64) string {
		return fmt.Sprintf("ROLLBACK TRANSACTION %d_%d", nodeID, ex.TransactionID)
	})
}

func (ex *connExecutor) Commit() error {
	if ex.nodes == nil {
		ex.nodes = map[uint64]*npgx.Conn{}
		return nil // There were never any transactions to release
	}

	golog.Debugf("committing changes on %d node(s), transaction ID [%d]", len(ex.nodes), ex.TransactionID)

	return ex.doTransaction(func(nodeID uint64) string {
		return "COMMIT"
	})
}

func (ex *connExecutor) Rollback() error {
	if ex.nodes == nil {
		ex.nodes = map[uint64]*npgx.Conn{}
		return nil // There were never any transactions to release
	}

	golog.Debugf("rolling back changes on %d node(s), transaction ID [%d]", len(ex.nodes), ex.TransactionID)

	return ex.doTransaction(func(nodeID uint64) string {
		return "ROLLBACK"
	})
}

func (ex *connExecutor) doTransaction(query func(nodeID uint64) string) error {
	if ex.nodes == nil {
		ex.nodes = map[uint64]*npgx.Conn{}
		return nil // There were never any transactions to release
	}

	// golog.Debugf("preparing two-phase commit for changes on %d node(s), transaction ID [%d]", len(ex.nodes), transactionID)

	count := len(ex.nodes)
	responses := make(chan executeResponse, count)
	for id, conn := range ex.nodes {
		func(conn *npgx.Conn) {
			response := executeResponse{
				NodeID: id,
			}
			defer func() { responses <- response }()
			result, err := conn.Query(query(id))
			if err != nil {
				response.Error = err
				return
			}

			result.Next()

			if err := result.Err(); err != nil {
				response.Error = err
			}
		}(conn)
	}

	errs := make([]error, 0)
	for i := 0; i < count; i++ {
		response := <-responses
		if response.Error != nil {
			errs = append(errs, response.Error)
		}
	}

	golog.Verbosef("%d error(s) from execution %v", len(errs), errs)
	return util.CombineErrors(errs)

}
