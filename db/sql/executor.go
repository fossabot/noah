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
 *
 * Project: pq github.com/lib/pq
 * Copyright 2018  'pq' Contributors Portions Copyright (C) 2011 Blake Mizerany
 * License https://github.com/lib/pq/blob/master/LICENSE.md
 */

package sql

import (
	"fmt"
	"github.com/Ready-Stock/noah/db/sql/driver/npgx"
	"github.com/Ready-Stock/noah/db/sql/pgwire/pgproto"
	"github.com/Ready-Stock/noah/db/sql/plan"
	"github.com/Ready-Stock/noah/db/sql/types"
	"github.com/Ready-Stock/noah/db/util"
	"github.com/ahmetb/go-linq"
	"github.com/juju/errors"
)

type executeResponse struct {
	Error  error
	Rows   *npgx.Rows
	NodeID uint64
}

func (ex *connExecutor) ExecutePlans(plans []plan.NodeExecutionPlan, res RestrictedCommandResult) (err error) {
	// defer util.CatchPanic(&err)
	if len(plans) == 0 {
		ex.Error("no plans were provided, nothing will be executed")
		return errors.New("no plans were provided")
	}
	responses := make(chan *executeResponse, len(plans))
	for _, p := range plans {
		go func(ex *connExecutor, pln plan.NodeExecutionPlan) {
			// defer util.CatchPanic(&err)
			exResponse := executeResponse{
				NodeID: pln.Node.NodeId,
			}

			defer func() {
				exResponse.Error = err
				responses <- &exResponse
			}()

			ex.Info("Executing query: `%s` on node [%d]", pln.CompiledQuery, pln.Node.NodeId)
			tx, err := ex.GetNodeTransaction(pln.Node.NodeId)
			if err != nil {
				ex.Error(err.Error())
				exResponse.Error = err
				return
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
	errs := make([]error, 0)
	for i := 0; i < len(plans); i++ {
		response := <-responses
		ex.Debug("handling response from node [%d]", response.NodeID)
		if response.Error != nil {
			ex.Error("received error from node [%d]: %s", response.NodeID, response.Error.Error())
			errs = append(errs, response.Error)
			continue
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
					errs = append(errs, err)
					continue
				} else {
					for c, v := range values {
						row[c] = v
					}
				}
				fmt.Printf("Row: %+v\n", row)
				res.AddRow(row)
				result = append(result, row)
			}

			if err := response.Rows.Err(); err != nil {
				ex.Error("received error from node [%d]: %s", response.NodeID, err)
				errs = append(errs, err)
				continue
			}

			response.Rows.Close()
		} else {
			ex.Debug("no rows returned for query `%s`", plans[0].CompiledQuery)
		}
	}
	ex.Debug("%d row(s) compiled for query `%s`", len(result), plans[0].CompiledQuery)

	if ex.TransactionMode == TransactionMode_AutoCommit && linq.From(plans).AnyWithT(func(plan plan.NodeExecutionPlan) bool {
		return !plan.ReadOnly
	}) { // If we are auto-committing this stuff and there are no errors
		if len(errs) == 0 {
			if err := ex.PrepareTwoPhase(); err != nil {
				return errors.Wrap(err, ex.Rollback()) // If the prepare two phase failed (which is really rare) then rollback the current changes in autocommit and return an error
			} else {
				return ex.CommitTwoPhase()
			}
		} else {
			return ex.Rollback()
		}
	}
	return util.CombineErrors(errs)
}

func (ex *connExecutor) PrepareTwoPhase() error {
	if ex.nodes == nil {
		ex.nodes = map[uint64]*npgx.Transaction{}
		return nil // There were never any transactions to release
	}

	transactionId := uint64(0)
	if ex.TransactionMode == TransactionMode_AutoCommit { // If we are auto-committing and we are targeting multiple nodes
		id, err := ex.SystemContext.NewSnowflake()
		if err != nil {
			return err
		}
		transactionId = id
		ex.TransactionID = id
	}

	transactionId = ex.TransactionID

	count := len(ex.nodes)
	responses := make(chan executeResponse, count)
	for id, tx := range ex.nodes {
		func(tx *npgx.Transaction) {
			response := executeResponse{
				NodeID: id,
			}
			defer func() { responses <- response }()
			response.Error = tx.PrepareTwoPhase(transactionId)
		}(tx)
	}
	errs := make([]error, 0)
	for i := 0; i < count; i++ {
		response := <- responses
		if response.Error != nil {
			errs = append(errs, response.Error)
		}
	}
	return util.CombineErrors(errs)
}

func (ex *connExecutor) CommitTwoPhase() error {
	if ex.nodes == nil {
		ex.nodes = map[uint64]*npgx.Transaction{}
		return nil // There were never any transactions to commit
	}

	count := len(ex.nodes)
	responses := make(chan executeResponse, count)
	for id, tx := range ex.nodes {
		func(tx *npgx.Transaction) {
			response := executeResponse{
				NodeID: id,
			}
			defer func() { responses <- response }()
			response.Error = tx.CommitTwoPhase()
		}(tx)
	}

	ex.nSync.Lock()
	defer ex.nSync.Unlock()
	errs := make([]error, 0)
	for i := 0; i < count; i++ {
		response := <- responses
		if response.Error != nil {
			errs = append(errs, response.Error)
		}
		delete(ex.nodes, response.NodeID)
	}
	return util.CombineErrors(errs)
}

func (ex *connExecutor) RollbackTwoPhase() error {
	if ex.nodes == nil {
		ex.nodes = map[uint64]*npgx.Transaction{}
		return nil // There were never any transactions to commit
	}

	count := len(ex.nodes)
	responses := make(chan executeResponse, count)
	for id, tx := range ex.nodes {
		func(tx *npgx.Transaction) {
			response := executeResponse{
				NodeID: id,
			}
			defer func() { responses <- response }()
			response.Error = tx.RollbackTwoPhase()
		}(tx)
	}

	ex.nSync.Lock()
	defer ex.nSync.Unlock()
	errs := make([]error, 0)
	for i := 0; i < count; i++ {
		response := <- responses
		if response.Error != nil {
			errs = append(errs, response.Error)
		}
		delete(ex.nodes, response.NodeID)
	}
	return util.CombineErrors(errs)
}

func (ex *connExecutor) Commit() error {
	if ex.nodes == nil {
		ex.nodes = map[uint64]*npgx.Transaction{}
		return nil // There were never any transactions to commit
	}

	count := len(ex.nodes)
	responses := make(chan executeResponse, count)
	for id, tx := range ex.nodes {
		func(tx *npgx.Transaction) {
			response := executeResponse{
				NodeID: id,
			}
			defer func() { responses <- response }()
			response.Error = tx.Commit()
		}(tx)
	}

	ex.nSync.Lock()
	defer ex.nSync.Unlock()
	errs := make([]error, 0)
	for i := 0; i < count; i++ {
		response := <- responses
		if response.Error != nil {
			errs = append(errs, response.Error)
		}
		delete(ex.nodes, response.NodeID)
	}
	return util.CombineErrors(errs)
}

func (ex *connExecutor) Rollback() error {
	if ex.nodes == nil {
		ex.nodes = map[uint64]*npgx.Transaction{}
		return nil // There were never any transactions to commit
	}

	count := len(ex.nodes)
	responses := make(chan executeResponse, count)
	for id, tx := range ex.nodes {
		func(tx *npgx.Transaction) {
			response := executeResponse{
				NodeID: id,
			}
			defer func() { responses <- response }()
			response.Error = tx.Rollback()
		}(tx)
	}

	ex.nSync.Lock()
	defer ex.nSync.Unlock()
	errs := make([]error, 0)
	for i := 0; i < count; i++ {
		response := <- responses
		if response.Error != nil {
			errs = append(errs, response.Error)
		}
		delete(ex.nodes, response.NodeID)
	}
	return util.CombineErrors(errs)
}
