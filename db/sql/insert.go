/*
 * Copyright (c) 2019 Ready Stock
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
	"context"
	"fmt"
	"github.com/gogo/protobuf/sortkeys"
	"strconv"
	"strings"

	"github.com/kataras/go-errors"
	"github.com/readystock/golinq"
	"github.com/readystock/golog"
	"github.com/readystock/noah/db/sql/plan"
	"github.com/readystock/noah/db/system"
	"github.com/readystock/pg_query_go/nodes"
)

type InsertStatement struct {
	Statement  pg_query.InsertStmt
	table      system.NTable
	accountIds []uint64
	IQueryStatement
}

func CreateInsertStatement(stmt pg_query.InsertStmt) *InsertStatement {
	return &InsertStatement{
		Statement: stmt,
	}
}

func (stmt *InsertStatement) Execute(ctx context.Context, ex *connExecutor, res RestrictedCommandResult, pinfo *plan.PlaceholderInfo) error {
	targetNodes, err := stmt.getTargetNodes(ctx, ex)
	if err != nil {
		return err
	}
	golog.Debugf("Preparing to send query to %d node(s)", len(targetNodes))

	plans, err := stmt.compilePlan(ctx, ex, targetNodes)
	if err != nil {
		return err
	}

	return ex.ExecutePlans(ctx, plans, res)
}

func (stmt *InsertStatement) getTargetNodes(ctx context.Context, ex *connExecutor) ([]system.NNode, error) {
	tableName := *stmt.Statement.Relation.Relname

	table, err := ex.SystemContext.Schema.GetTable(tableName)
	if err != nil {
		return nil, err
	}

	if table == nil {
		return nil, errors.New(fmt.Sprintf("table [%s] does not exist", tableName))
	}

	stmt.table = *table

	// If the insert query targets global or account tables then we want to target all nodes.
	if table.TableType == system.NTableType_GLOBAL || table.TableType == system.NTableType_ACCOUNT {
		return ex.SystemContext.Nodes.GetNodes()
	} else {
		accountIds, err := stmt.getAccountIds(ctx, *table)
		if err != nil {
			return nil, err
		}

		stmt.accountIds = accountIds

		if len(accountIds) > 1 {
			return nil, errors.New("cannot insert into more than 1 account ID at this time")
		}
		return ex.SystemContext.Accounts.GetNodesForAccount(accountIds[0])
	}
}

func (stmt *InsertStatement) getColumnIndex(ctx context.Context, colName string) int {
	return linq.From(stmt.Statement.Cols.Items).IndexOfT(func(col pg_query.ResTarget) bool {
		return strings.ToLower(*col.Name) == strings.ToLower(colName)
	})
}

func (stmt *InsertStatement) getAccountIds(ctx context.Context, table system.NTable) ([]uint64, error) {
	// Since we are inserting into a sharded table we need to find the shard key in the insert
	// statement. If the shard key is missing we want to throw an error, if its present we want
	// to take its index and look at its provided value.
	shardKey := table.ShardKey.(*system.NTable_SKey).SKey

	shardKeyIndex := stmt.getColumnIndex(ctx, shardKey.ColumnName)

	// We couldn't find a value for the shard key, throw an error and return the shard column name
	if shardKeyIndex < 0 {
		return nil, fmt.Errorf("insert statement is missing the shard column [%s] value", shardKey.ColumnName)
	}

	if stmt.Statement.SelectStmt == nil || stmt.Statement.SelectStmt.(pg_query.SelectStmt).ValuesLists == nil {
		return nil, errors.New("value list was not provided, these types of inserts are not yet supported")
	}

	valuesList := stmt.Statement.SelectStmt.(pg_query.SelectStmt).ValuesLists

	accountIds := make([]uint64, 0)

	linq.From(valuesList).SelectT(func(nodes []pg_query.Node) uint64 {
		column := nodes[shardKeyIndex]
		val, err := column.Deparse(pg_query.Context_None)
		if err != nil {
			return 0
		}
		iVal, err := strconv.ParseUint(*val, 10, 64)
		if err != nil {
			return 0
		}
		return iVal
	}).Distinct().ToSlice(&accountIds)

	if linq.From(accountIds).AnyWithT(func(id uint64) bool {
		return id == 0
	}) {
		return nil, errors.New("there are invalid account IDs provided in insert values")
	}

	return accountIds, nil
}

func (stmt *InsertStatement) compilePlan(ctx context.Context, ex *connExecutor, nodes []system.NNode) ([]plan.NodeExecutionPlan, error) {
	plans := make([]plan.NodeExecutionPlan, len(nodes))

	for _, column := range stmt.table.Columns {
		columnIndex := stmt.getColumnIndex(ctx, column.ColumnName)

		if columnIndex < 0 {
			if column.IsSequence {
				// If the index is -1 then the column is not specified in the insert and we need to
				// add it.
				stmt.Statement.Cols.Items = append(stmt.Statement.Cols.Items, pg_query.ResTarget{
					Name: &column.ColumnName,
				})
				values := stmt.Statement.SelectStmt.(pg_query.SelectStmt)

				for i, valueList := range values.ValuesLists {
					nextId := uint64(0)
					// If this is an account table that we are inserting into then we want to create
					// new account records
					if stmt.table.TableType == system.NTableType_ACCOUNT {
						account, _, err := ex.SystemContext.Accounts.CreateAccount()
						if err != nil {
							return nil, err
						}
						nextId = account.AccountId
					} else {
						next, err := ex.SystemContext.Sequences.GetNextValueForSequence(fmt.Sprintf("%s.%s", stmt.table.TableName, column.ColumnName))
						if err != nil {
							return nil, err
						}
						nextId = *next
					}

					valueList = append(valueList, pg_query.A_Const{
						Val: pg_query.Integer{
							Ival: int64(nextId),
						},
					})

					values.ValuesLists[i] = valueList
				}

				stmt.Statement.SelectStmt = values
			}
		} else {
			if column.IsSequence {
				return nil, fmt.Errorf("cannot specify value when inserting into serialized column [%s]", column.ColumnName)
			}
			// If this is a sharded table and we have a value in this column, check to see if this is
			// the shard column. If it is we want to make sure that it is distinct.
			if stmt.table.TableType == system.NTableType_SHARD && (stmt.table.ShardKey.(*system.NTable_SKey)).SKey.ColumnName == column.ColumnName {
				values := stmt.Statement.SelectStmt.(pg_query.SelectStmt)

				shardKeys := make([]uint64, 0)
				linq.From(values.ValuesLists).SelectT(func(val []pg_query.Node) uint64 {
					aConst, ok := val[columnIndex].(pg_query.A_Const)
					if !ok {
						return 0
					}

					aInteger, ok := aConst.Val.(pg_query.Integer)
					if !ok {
						return 0
					}

					return uint64(aInteger.Ival)
				}).Distinct().ToSlice(&shardKeys)

				sortkeys.Uint64s(shardKeys)

				if shardKeys[0] == 0 {
					return nil, fmt.Errorf("invalid shard key was provided for column [%s]", column.ColumnName)
				}

				if len(shardKeys) > 1 {
					// Split the insert stmt into multiple inserts
					return nil, fmt.Errorf("multi-shard inserts are not supported at this time")
				}
			}
		}
	}

	deparsed, err := pg_query.Deparse(stmt.Statement)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(plans); i++ {
		// If we are generating a plan for more than one node, then we only want to have a returning
		// clause on one of the queries executed to make sure that duplicate data isn't returned.
		if i > 0 && len(stmt.Statement.ReturningList.Items) > 0 {
			stmt.Statement.ReturningList.Items = make([]pg_query.Node, 0)
			deparsed, err = pg_query.Deparse(stmt.Statement)
			if err != nil {
				return nil, err
			}
		}

		plans[i] = plan.NodeExecutionPlan{
			CompiledQuery: *deparsed,
			Node:          nodes[i],
			ReadOnly:      false,
			Type:          stmt.Statement.StatementType(),
		}
	}
	return plans, nil
}
