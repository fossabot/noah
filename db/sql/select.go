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
	"github.com/Ready-Stock/Noah/db/sql/plan"
	"github.com/Ready-Stock/Noah/db/system"
	pg_query2 "github.com/Ready-Stock/pg_query_go"
	"github.com/Ready-Stock/pg_query_go/nodes"
	"github.com/kataras/go-errors"
)

type SelectStatement struct {
	Statement pg_query.SelectStmt
	IQueryStatement
}

func CreateSelectStatement(stmt pg_query.SelectStmt) *SelectStatement {
	return &SelectStatement{
		Statement: stmt,
	}
}

func (stmt *SelectStatement) Execute(ex *connExecutor, res RestrictedCommandResult) error {
	targetNodes, err := stmt.getTargetNodes(ex)
	if err != nil {
		return err
	}
	ex.Debug("Preparing to send query to %d node(s)", len(targetNodes))

	plans, err := stmt.compilePlan(ex, targetNodes)
	if err != nil {
		return err
	}

	return ex.ExecutePlans(plans, res)
}

func (stmt *SelectStatement) getTargetNodes(ex *connExecutor) ([]system.NNode, error) {
	accounts, err := stmt.getAccountIDs()
	if err != nil {
		return nil, err
	}

	if len(accounts) == 1 {
		return ex.SystemContext.GetNodesForAccount(accounts[0])
	} else if len(accounts) > 1 {
		return nil, errors.New("multi account queries are not supported at this time.")
		// node_ids := make([]uint64, 0)
		// From(accounts).SelectManyT(func(id uint64) Query {
		// 	if ids, err := ex.GetNodesForAccountID(&id); err == nil {
		// 		return From(ids)
		// 	}
		// 	return From(make([]uint64, 0))
		// }).Distinct().ToSlice(&node_ids)
		// if len(node_ids) == 0 {
		// 	return nil, errors.New("could not find nodes for account IDs")
		// } else {
		// 	return node_ids, nil
		// }
	} else {
		return ex.SystemContext.GetNodes()
	}
}

func (stmt *SelectStatement) getAccountIDs() ([]uint64, error) {

	return make([]uint64, 0), nil
}

func (stmt *SelectStatement) compilePlan(ex *connExecutor, nodes []system.NNode) ([]plan.NodeExecutionPlan, error) {
	plans := make([]plan.NodeExecutionPlan, len(nodes))
	deparsed, err := pg_query2.Deparse(stmt.Statement)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(plans); i++ {
		plans[i] = plan.NodeExecutionPlan{
			CompiledQuery: *deparsed,
			Node:          nodes[i],
			ReadOnly:      true,
		}
	}
	return plans, nil
}
