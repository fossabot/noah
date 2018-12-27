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
	"github.com/ahmetb/go-linq"
	"github.com/juju/errors"
	"github.com/readystock/golog"
	"github.com/readystock/noah/db/sql/plan"
	"github.com/readystock/noah/db/system"
	"github.com/readystock/pg_query_go/nodes"
	"strings"
)

var (
	ErrNotEnoughNodesAvailable = errors.New("not enough nodes available in cluster to create table")
	ErrTablespaceNotSpecified  = errors.New("tablespace must be specified when creating a table")
)

type CreateStatement struct {
	Statement pg_query.CreateStmt
	table     system.NTable
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
	//     if err := ex.execStmt(pg_query.TransactionStmt{Kind: pg_query.TRANS_STMT_BEGIN}, nil, 0); err != nil {
	//         return err
	//     }
	//     ex.TransactionStatus = NTXInProgress
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

	// Schema changes can only be made when all (non-replica) nodes are alive, if any nodes are
	// unavailable then the schema change will be rejected to ensure consistency.
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

	stmt.table = system.NTable{
		TableName: *stmt.Statement.Relation.Relname,
		Schema:    "default",
		Columns:   make([]*system.NColumn, len(stmt.Statement.TableElts.Items)),
	}

	// Determine the distribution of a table in the cluster
	if err := stmt.handleTableType(ex, &stmt.table); err != nil { // Handle sharding
		return nil, err
	}

	// Add handling here for custom column types.
	if err := stmt.handleColumns(ex, &stmt.table); err != nil { // Handle sharding
		return nil, err
	}

	// // Create the table in the coordinator cluster
	// if err := ex.SystemContext.Schema.CreateTable(table); err != nil {
	// 	return nil, err
	// }

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
	if table.TableType == system.NTableType_ACCOUNT {
		if accountsTable, err := ex.SystemContext.Schema.GetAccountsTable(); err != nil {
			return err
		} else if accountsTable != nil {
			return errors.Errorf("an accounts table named [%s] already exists in this cluster", accountsTable.TableName)
		}
	}

	return nil
}

func (stmt *CreateStatement) handleColumns(ex *connExecutor, table *system.NTable) error {
	if stmt.Statement.TableElts.Items != nil && len(stmt.Statement.TableElts.Items) > 0 {
		for i, col := range stmt.Statement.TableElts.Items {
			switch tableItem := col.(type) {
			case pg_query.ColumnDef:
				noahColumn := &system.NColumn{ColumnName: *tableItem.Colname}

				// There are a few types that are handled by noah as a middle man; such as ID generation
				// because of this we want to replace serial columns with their base types since noah
				// will rewrite incoming queries to include an ID when performing inserts.
				if tableItem.TypeName != nil &&
					tableItem.TypeName.Names.Items != nil &&
					len(tableItem.TypeName.Names.Items) > 0 {
					columnType := tableItem.TypeName.Names.Items[len(tableItem.TypeName.Names.Items)-1].(pg_query.String) // The last type name
					golog.Verbosef("Processing column [%s] type [%s]", *tableItem.Colname, strings.ToLower(columnType.Str))
					// This switch statement will handle any custom column types that we would like.
					switch strings.ToLower(columnType.Str) {
					case "serial": // Emulate 32 bit sequence
						columnType.Str = "int"
						noahColumn.IsSequence = true
					case "snowflake": // Snowflakes are generated using twitters id system
						noahColumn.IsSnowflake = true
						fallthrough
					case "bigserial": // Emulate 64 bit sequence
						columnType.Str = "bigint"
						noahColumn.IsSequence = true
					default:
						// Other column types wont be handled.
					}
					noahColumn.ColumnTypeName = strings.ToLower(columnType.Str)
					tableItem.TypeName.Names.Items = []pg_query.Node{columnType}
					stmt.Statement.TableElts.Items[i] = tableItem
				}

				// Check to see if this column is the primary key, primary keys will be used for tables
				// like account tables. If someone tries to create a table without a primary key an
				// error will be returned at this time.
				if tableItem.Constraints.Items != nil && len(tableItem.Constraints.Items) > 0 {
					noahColumn.IsPrimaryKey = linq.From(tableItem.Constraints.Items).AnyWithT(func(constraint pg_query.Constraint) bool {
						return constraint.Contype == pg_query.CONSTR_PRIMARY
					})

					if noahColumn.IsPrimaryKey {
						if table.PrimaryKey != nil {
							return errors.New("cannot define more than 1 primary key on a single table")
						}

						switch noahColumn.ColumnTypeName {
						case "bigint", "int", "tinyint":
						default:
							// At the moment noah only supports integer column sharding.
							return errors.Errorf("column [%s] cannot be a primary key, a primary key must be an integer column", noahColumn.ColumnName)
						}

						table.PrimaryKey = &system.NTable_PKey{
							PKey: noahColumn,
						}
					}
				}

				table.Columns[i] = noahColumn
			case pg_query.Constraint:
				// Its possible for primary keys, foreign keys and identities to be defined
				// somewhere other than the column line itself, if this happens we still want to
				// handle it gracefully.
				switch tableItem.Contype {
				case pg_query.CONSTR_PRIMARY:
				case pg_query.CONSTR_FOREIGN:
				case pg_query.CONSTR_IDENTITY:
				}
			}
		}
	}

	if table.PrimaryKey == nil && table.TableType == system.NTableType_ACCOUNT {
		return errors.New("cannot create an account table without a primary key")
	}

	return nil
}

func (stmt *CreateStatement) handleTableType(ex *connExecutor, table *system.NTable) error {
	if stmt.Statement.Tablespacename != nil {
		switch strings.ToLower(*stmt.Statement.Tablespacename) {
		case "noah.global": // Table has the same data on all shards
			table.TableType = system.NTableType_GLOBAL
		case "noah.shard": // Table is sharded by shard column
			table.TableType = system.NTableType_SHARD
		case "noah.account": // Table contains all of the records of accounts for cluster
			table.TableType = system.NTableType_ACCOUNT
		default: // Other
			return ErrTablespaceNotSpecified
		}
	}
	return nil
}
