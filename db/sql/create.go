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
	"github.com/Ready-Stock/Noah/db/sql/plan"
	"github.com/Ready-Stock/Noah/db/system"
	parser "github.com/Ready-Stock/pg_query_go"
	"github.com/Ready-Stock/pg_query_go/nodes"
	"github.com/ahmetb/go-linq"
	"github.com/kataras/go-errors"
	"strings"
)

var (
	ErrTableExists = errors.New("table with name [%s] already exists in the cluster")
)

type CreateStatement struct {
	Statement pg_query.CreateStmt
	IQueryStatement
}

func CreateCreateStatement(stmt pg_query.CreateStmt) *CreateStatement {
	return &CreateStatement{
		Statement: stmt,
	}
}

func (stmt *CreateStatement) Execute(ex *connExecutor, res RestrictedCommandResult) error {
	existingTable, _ := ex.SystemContext.Schema.GetTable(*stmt.Statement.Relation.Relname)
	if existingTable != nil {
		return ErrTableExists.Format(*stmt.Statement.Relation.Relname)
	}

	targetNodes, err := stmt.getTargetNodes(ex)
	if err != nil {
		return err
	}

	plans, err := stmt.compilePlan(ex, targetNodes)
	if err != nil {
		return err
	}

	return ex.ExecutePlans(plans, res)
}

func (stmt *CreateStatement) getTargetNodes(ex *connExecutor) ([]system.NNode, error) {
	writeNodes, err := ex.SystemContext.Nodes.GetNodes()
	if err != nil {
		return nil, err
	}

	allNodes := make([]system.NNode, 0)
	linq.From(writeNodes).WhereT(func(node system.NNode) bool {
		return node.ReplicaOf == 0
	}).ToSlice(&allNodes)

	liveNodes := linq.From(allNodes).CountWithT(func(node system.NNode) bool {
		return node.IsAlive && node.ReplicaOf == 0
	})

	if liveNodes != len(allNodes) {
		return nil, errors.New("not enough nodes available in cluster to create table")
	}

	return allNodes, nil
}

func (stmt *CreateStatement) compilePlan(ex *connExecutor, nodes []system.NNode) ([]plan.NodeExecutionPlan, error) {
	plans := make([]plan.NodeExecutionPlan, len(nodes))

	stmt.handleTableType(ex) // Handle sharding

	// Add handling here for custom column types.

	deparsed, err := parser.Deparse(stmt.Statement)
	if err != nil {
		ex.Error(err.Error())
		return nil, err
	}
	ex.Debug(*deparsed)
	for i := 0; i < len(plans); i++ {
		plans[i] = plan.NodeExecutionPlan{
			CompiledQuery: *deparsed,
			Node:          nodes[i],
			ReadOnly:      false,
		}
	}
	return plans, nil
}

func (stmt *CreateStatement) handleSequences() {
	if stmt.Statement.TableElts.Items != nil && len(stmt.Statement.TableElts.Items) > 0 {
		for _, col := range stmt.Statement.TableElts.Items {
			columnDefinition := col.(pg_query.ColumnDef)
			if columnDefinition.TypeName != nil &&
				columnDefinition.TypeName.Names.Items != nil &&
				len(columnDefinition.TypeName.Names.Items) > 0 {
				columnType := columnDefinition.TypeName.Names.Items[len(columnDefinition.TypeName.Names.Items)].(pg_query.String) // The last type name

				// This switch statement will handle any custom column types that we would like.
				switch strings.ToLower(columnType.Str) {
				case "nserial": // Emulate 32 bit sequence

				case "nbigserial": // Emulate 64 bit sequence

				default:

				}
			}
		}
	}
}

func (stmt *CreateStatement) handleTableType(ex *connExecutor) error {
	if stmt.Statement.Tablespacename != nil {
		switch strings.ToLower(*stmt.Statement.Tablespacename) {
		case "global": // Table has the same data on all shards

		case "account": // Table is sharded by shard column

		default: // Other

		}
	}
	return nil
}
