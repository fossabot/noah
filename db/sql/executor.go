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
 */

package sql

import (
	"github.com/Ready-Stock/Noah/db/sql/driver/npgx"
	"github.com/Ready-Stock/Noah/db/sql/plan"
	"github.com/Ready-Stock/Noah/db/sql/types"
)

type executeResponse struct {
	Error   error
	Rows    *npgx.Rows
	NodeID  uint64
}

func (ex *connExecutor) ExecutePlans(plans []plan.NodeExecutionPlan) error {
	responses := make(chan *executeResponse, len(plans))
	for _, p := range plans {
		go func(ex *connExecutor, pln plan.NodeExecutionPlan) {
			exResponse := executeResponse{
				NodeID: pln.Node.NodeID,
			}
			defer func() { responses <- &exResponse }()

			if !pln.Node.Alive {
				ex.Warn("Deferring query: `%s` for node [%d]", pln.CompiledQuery, pln.Node.NodeID)
				return
			}

			ex.Info("Executing query: `%s` on node [%d]", pln.CompiledQuery, pln.Node.NodeID)
			tx, ok := ex.GetNodeTransaction(pln.Node.NodeID)
			if !ok {
				// A connection has not yet been made to this node. Allocate one.
				if t, err := ex.SystemContext.Pool.AcquireTransaction(pln.Node.NodeID); err != nil {
					ex.Error(err.Error())
					exResponse.Error = err
					return
				} else {
					ex.SetNodeTransaction(pln.Node.NodeID, t)
					tx = t
				}
			}

			if ex.TransactionStatus == NTXNoTransaction {
				ex.TransactionStatus = NTXInProgress
			}

			if pln.ReadOnly {
				rows, err := tx.Query(pln.CompiledQuery)
				if err != nil {
					ex.Error(err.Error())
					exResponse.Error = err
					return
				}
				exResponse.Rows = rows
			} else {

			}
		}(ex, p)
	}
	columns := make([]npgx.FieldDescription, 0)
	result := make([][]types.Value, 0)
	for i := 0; i < len(plans); i++ {
		response := <-responses
		if response.Error != nil {
			return response.Error
		}
		if response.Rows != nil {
			for response.Rows.Next() {
				if len(columns) == 0 {
					columns = response.Rows.FieldDescriptions()
				}
				row := make([]types.Value, len(columns))
				if values, err := response.Rows.Values(); err != nil {

				} else {
					for c, v := range values {
						if tValue, ok := v.(types.Value); ok {
							row[c] = tValue
						}
					}
				}
				result = append(result, row)
			}
		}
	}
	return nil
}

func (ex *connExecutor) PrepareCommit() error {
	responses := make(chan *executeResponse, len(ex.nodes))
	for nodeId, tx := range ex.nodes {
		go func(tx *npgx.Transaction) {
			exResponse := executeResponse{
				NodeID: nodeId,
			}
			defer func() { responses <- &exResponse }()
			node, err := ex.SystemContext.GetNode(nodeId)
			if err != nil {
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
	return nil
}
