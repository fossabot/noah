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
	"github.com/Ready-Stock/Noah/db/sql/driver/npgx"
	"github.com/Ready-Stock/Noah/db/sql/pgwire/pgproto"
	"github.com/Ready-Stock/Noah/db/sql/plan"
	"github.com/Ready-Stock/Noah/db/sql/types"
	"github.com/kataras/go-errors"
)

type executeResponse struct {
	Error   error
	Rows    *npgx.Rows
	NodeID  uint64
}

func (ex *connExecutor) ExecutePlans(plans []plan.NodeExecutionPlan, res RestrictedCommandResult) (err error) {
	//defer util.CatchPanic(&err)
	if len(plans) == 0 {
		ex.Error("no plans were provided, nothing will be executed")
		return errors.New("no plans were provided")
	}
	responses := make(chan *executeResponse, len(plans))
	for _, p := range plans {
		go func(ex *connExecutor, pln plan.NodeExecutionPlan) {
			//defer util.CatchPanic(&err)
			exResponse := executeResponse{
				NodeID: pln.Node.NodeId,
			}

			defer func() {
				exResponse.Error = err
				responses <- &exResponse
			}()

			// if !pln.Node.Alive {
			// 	ex.Warn("Deferring query: `%s` for node [%d]", pln.CompiledQuery, pln.Node.NodeID)
			// 	return
			// }

			ex.Info("Executing query: `%s` on node [%d]", pln.CompiledQuery, pln.Node.NodeId)
			tx, ok := ex.GetNodeTransaction(pln.Node.NodeId)
			if !ok {
				ex.Debug("node [%d] is not in the session, acquiring connection", pln.Node.NodeId)
				// A connection has not yet been made to this node. Allocate one.
				if t, err := ex.SystemContext.Pool.AcquireTransaction(pln.Node.NodeId); err != nil {
					ex.Error(err.Error())
					exResponse.Error = err
					return
				} else {
					ex.SetNodeTransaction(pln.Node.NodeId, t)
					tx = t
				}
			}

			if ex.TransactionStatus == NTXNoTransaction {
				ex.TransactionStatus = NTXInProgress
			}
			rows, err := tx.Query(pln.CompiledQuery)
			if err != nil {
				ex.Error(err.Error())
				exResponse.Error = err
				return
			}
			ex.Debug("received rows response from node [%d]", pln.Node.NodeId)
			exResponse.Rows = rows
		}(ex, p)
	}
	columns := make([]pgproto.FieldDescription, 0)
	result := make([][]types.Value, 0)
	for i := 0; i < len(plans); i++ {
		response := <- responses
		ex.Debug("handling response from node [%d]", response.NodeID)
		if response.Error != nil {
			return response.Error
		}
		if response.Rows != nil {
			for response.Rows.Next() {
				if len(columns) == 0 {
					columns = response.Rows.PgFieldDescriptions()
					ex.Debug("retrieved %d column(s)", len(columns))
					res.SetColumns(columns)
				}
				row := make([]types.Value, len(columns))
				if values, err := response.Rows.PgValues(); err != nil {
					ex.Error(err.Error())
					return err
				} else {
					for c, v := range values {
						row[c] = v
					}
				}
				res.AddRow(row)
				result = append(result, row)
			}
		} else {
			ex.Debug("no rows returned for query `%s`", plans[0].CompiledQuery)
		}
	}
	ex.Debug("%d row(s) compiled for query `%s`", len(result), plans[0].CompiledQuery)
	return nil
}

func (ex *connExecutor) PrepareTwoPhase() error {
	responses := make(chan *executeResponse, len(ex.nodes))
	for nodeId, tx := range ex.nodes {
		go func(tx *npgx.Transaction) {
			exResponse := executeResponse{
				NodeID: nodeId,
			}
			defer func() { responses <- &exResponse }()
			node, err := ex.SystemContext.Nodes.GetNode(nodeId)
			if err != nil {
				ex.Error(err.Error())
				exResponse.Error = err
				return
			}

			if !node.IsAlive {
				return // TODO (elliotcourant) Add handling for a dead node.
			}

			if err := tx.PrepareTwoPhase(ex.TransactionID); err != nil {
				ex.Error(err.Error())
				exResponse.Error = err
				return
			}
		}(tx)
	}
	for i := 0; i < len(ex.nodes); i++ {
		response := <- responses
		if response.Error != nil {
			return response.Error
		}
	}
	ex.TransactionStatus = NTXPreparedSuccess
	return nil
}

func (ex *connExecutor) CommitTwoPhase() error {
	return nil
}

func (ex *connExecutor) RollbackTwoPhase() error {
	return nil
}

func (ex *connExecutor) Commit() error {
	return nil
}

func (ex *connExecutor) Rollback() error {
	return nil
}