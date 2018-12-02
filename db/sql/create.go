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
    "github.com/Ready-Stock/noah/db/sql/plan"
    "github.com/Ready-Stock/noah/db/system"
	"github.com/Ready-Stock/pg_query_go/nodes"
	"github.com/ahmetb/go-linq"
	"github.com/juju/errors"
	"strings"
)

var (
	ErrNotEnoughNodesAvailable = errors.New("not enough nodes available in cluster to create table")
	ErrTablespaceNotSpecified  = errors.New("tablespace must be specified when creating a table")
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
		return errors.Errorf("table with name [%s] already exists in the cluster", *stmt.Statement.Relation.Relname)
	}

	targetNodes, err := stmt.getTargetNodes(ex)
	if err != nil {
		return err
	}

	plans, err := stmt.compilePlan(ex, targetNodes)
	if err != nil {
		return err
	}

	// if ex.TransactionState == NTXNoTransaction { // When creating a table try to create a new transaction if we are not in one
	// 	if err := ex.execStmt(pg_query.TransactionStmt{Kind: pg_query.TRANS_STMT_BEGIN}, nil, 0); err != nil {
	// 		return err
	// 	}
	// 	ex.TransactionStatus = NTXInProgress
	// }

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
		return nil, ErrNotEnoughNodesAvailable
	}

	if liveNodes == 0 {
		return nil, errors.Errorf("no live nodes, ddl cannot be processed at this time")
	}

	return allNodes, nil
}

func (stmt *CreateStatement) compilePlan(ex *connExecutor, nodes []system.NNode) ([]plan.NodeExecutionPlan, error) {
	plans := make([]plan.NodeExecutionPlan, len(nodes))

	table := system.NTable{
		TableName: *stmt.Statement.Relation.Relname,
		Schema:    "default",
		Columns:   make([]*system.NColumn, len(stmt.Statement.TableElts.Items)),
	}

	// Determine the distribution of a table in the cluster
	if err := stmt.handleTableType(ex, &table); err != nil { // Handle sharding
		return nil, err
	}

	// Add handling here for custom column types.
	if err := stmt.handleColumns(ex, &table); err != nil { // Handle sharding
		return nil, err
	}

	// Create the table in the coordinator cluster
	if err := ex.SystemContext.Schema.CreateTable(table); err != nil {
		return nil, err
	}

	deparsed, err := pg_query.Deparse(stmt.Statement)
	if err != nil {
		ex.Error(err.Error())
		return nil, err
	}
	ex.Debug("Recompiled query: %s", *deparsed)
	for i := 0; i < len(plans); i++ {
		plans[i] = plan.NodeExecutionPlan{
			CompiledQuery: *deparsed,
			Node:          nodes[i],
			ReadOnly:      false,
		}
	}
	return plans, nil
}

func (stmt *CreateStatement) handleValidation(ex *connExecutor, table *system.NTable) error {
	if table.TableType == system.NTableType_Account {
		if accountsTable, err := ex.SystemContext.Schema.GetAccountsTable(); err != nil {
			return err
		} else if accountsTable != nil {
			return errors.Errorf("an accounts table named [%s] already exists in this cluster", accountsTable.TableName)
		}
	}

	return nil
}

func (stmt *CreateStatement) handleColumns(ex *connExecutor, table *system.NTable) error {
	primaryKeyFound := false
	if stmt.Statement.TableElts.Items != nil && len(stmt.Statement.TableElts.Items) > 0 {
		for i, col := range stmt.Statement.TableElts.Items {
			columnDefinition := col.(pg_query.ColumnDef)
			table.Columns[i] = &system.NColumn{ColumnName: *columnDefinition.Colname}

			if columnDefinition.Constraints.Items != nil && len(columnDefinition.Constraints.Items) > 0 {
				table.Columns[i].IsPrimaryKey = linq.From(columnDefinition.Constraints.Items).AnyWithT(func(constraint pg_query.Constraint) bool {
					return constraint.Contype == pg_query.CONSTR_PRIMARY
				})
				if table.Columns[i].IsPrimaryKey {
					if primaryKeyFound {
						return errors.New("cannot define more than 1 primary key on a single table")
					}
					primaryKeyFound = true
				}
			}

			if columnDefinition.TypeName != nil &&
				columnDefinition.TypeName.Names.Items != nil &&
				len(columnDefinition.TypeName.Names.Items) > 0 {
				columnType := columnDefinition.TypeName.Names.Items[len(columnDefinition.TypeName.Names.Items)-1].(pg_query.String) // The last type name
				ex.Debug("Processing column [%s] type [%s]", *columnDefinition.Colname, strings.ToLower(columnType.Str))
				// This switch statement will handle any custom column types that we would like.
				switch strings.ToLower(columnType.Str) {
				case "serial": // Emulate 32 bit sequence
					columnType.Str = "int"
					table.Columns[i].IsSequence = true
				case "snowflake": // Snowflakes are generated using twitters id system
					table.Columns[i].IsSnowflake = true
					fallthrough
				case "bigserial": // Emulate 64 bit sequence
					columnType.Str = "bigint"
					table.Columns[i].IsSequence = true
				default:
					// Other column types wont be handled.
				}
				table.Columns[i].ColumnTypeName = strings.ToLower(columnType.Str)
				columnDefinition.TypeName.Names.Items = []pg_query.Node{columnType}
				stmt.Statement.TableElts.Items[i] = columnDefinition
			}
		}
	}

	if !primaryKeyFound && table.TableType == system.NTableType_Account {
		return errors.New("cannot create an account table without a primary key")
	}

	return nil
}

func (stmt *CreateStatement) handleTableType(ex *connExecutor, table *system.NTable) error {
	if stmt.Statement.Tablespacename != nil {
		switch strings.ToLower(*stmt.Statement.Tablespacename) {
		case "global": // Table has the same data on all shards
			table.TableType = system.NTableType_Global
		case "shard": // Table is sharded by shard column
			table.TableType = system.NTableType_Shard
		case "account": // Table contains all of the records of accounts for cluster
			table.TableType = system.NTableType_Account
		default: // Other
			return ErrTablespaceNotSpecified
		}
	}
	return nil
}
