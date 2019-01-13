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
    "github.com/readystock/golinq"
    "github.com/kataras/go-errors"
    "github.com/readystock/noah/db/sql/plan"
    "github.com/readystock/noah/db/system"
    "github.com/readystock/pg_query_go/nodes"
    "strconv"
    "strings"
)

type InsertStatement struct {
    Statement pg_query.InsertStmt
    IQueryStatement
}

func CreateInsertStatement(stmt pg_query.InsertStmt) *InsertStatement {
    return &InsertStatement{
        Statement: stmt,
    }
}

func (stmt *InsertStatement) Execute(ex *connExecutor, res RestrictedCommandResult, pinfo *plan.PlaceholderInfo) error {
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

func (stmt *InsertStatement) getTargetNodes(ex *connExecutor) ([]system.NNode, error) {
    tableName := *stmt.Statement.Relation.Relname

    table, err := ex.SystemContext.Schema.GetTable(tableName)
    if err != nil {
        return nil, err
    }

    if table == nil {
        return nil, errors.New(fmt.Sprintf("table [%s] does not exist", tableName))
    }

    // If the insert query targets global or account tables then we want to target all nodes.
    if table.TableType == system.NTableType_GLOBAL || table.TableType == system.NTableType_ACCOUNT {
        return ex.SystemContext.Nodes.GetNodes()
    } else {
        accountIds, err := stmt.getAccountIds(*table)
        if err != nil {
            return nil, err
        }

        if len(accountIds) > 1 {
            return nil, errors.New("cannot insert into more than 1 account ID at this time")
        }

        return ex.SystemContext.Accounts.GetNodesForAccount(accountIds[0])
    }
}

func (stmt *InsertStatement) getAccountIds(table system.NTable) ([]uint64, error) {
    // Since we are inserting into a sharded table we need to find the shard key in the insert
    // statement. If the shard key is missing we want to throw an error, if its present we want
    // to take its index and look at its provided value.
    shardKey := table.ShardKey.(*system.NTable_SKey).SKey

    shardKeyIndex := linq.From(stmt.Statement.Cols.Items).IndexOfT(func(col pg_query.ResTarget) bool {
        return strings.ToLower(*col.Name) == shardKey.ColumnName
    })

    // We couldn't find a value for the shard key, throw an error and return the shard column name
    if shardKeyIndex < 0 {
        return nil, errors.New(fmt.Sprintf("insert statement is missing the shard column [%s] value", shardKey.ColumnName))
    }

    if  stmt.Statement.SelectStmt == nil || stmt.Statement.SelectStmt.(pg_query.SelectStmt).ValuesLists == nil {
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

func (stmt *InsertStatement) compilePlan(ex *connExecutor, nodes []system.NNode) ([]plan.NodeExecutionPlan, error) {
    plans := make([]plan.NodeExecutionPlan, len(nodes))
    deparsed, err := pg_query.Deparse(stmt.Statement)
    if err != nil {
        return nil, err
    }
    for i := 0; i < len(plans); i++ {
        plans[i] = plan.NodeExecutionPlan{
            CompiledQuery: *deparsed,
            Node:          nodes[i],
            ReadOnly:      true,
            Type:          stmt.Statement.StatementType(),
        }
    }
    return plans, nil
}
