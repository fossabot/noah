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
	"github.com/ahmetb/go-linq"
	"github.com/kataras/go-errors"
	"github.com/readystock/golog"
	"github.com/readystock/noah/db/sql/plan"
	"github.com/readystock/noah/db/system"
	"github.com/readystock/noah/db/util/queryutil"
	"github.com/readystock/pg_query_go/nodes"
	"math/rand"
	"time"
)

type SelectStatement struct {
	Statement pg_query.SelectStmt
	tables    []system.NTable
	IQueryStatement
}

func CreateSelectStatement(stmt pg_query.SelectStmt) *SelectStatement {
	return &SelectStatement{
		Statement: stmt,
	}
}

func (stmt *SelectStatement) Execute(ctx context.Context, ex *connExecutor, res RestrictedCommandResult, pinfo *plan.PlaceholderInfo) error {
	targetNodes, err := stmt.getTargetNodes(ctx, ex)
	if err != nil {
		return err
	}
	golog.Debugf("Preparing to send query to %d node(s)", len(targetNodes))

	if err := stmt.replaceParameters(ctx, ex, pinfo); err != nil {
		return err
	}

	plans, err := stmt.compilePlan(ctx, ex, targetNodes)
	if err != nil {
		return err
	}

	return ex.ExecutePlans(ctx, plans, res)
}

func (stmt *SelectStatement) replaceParameters(ctx context.Context, ex *connExecutor, pinfo *plan.PlaceholderInfo) error {
	newStmt := queryutil.ReplaceArguments(stmt.Statement, pinfo.Values)
	stmt.Statement = newStmt.(pg_query.SelectStmt)
	return nil
}

func (stmt *SelectStatement) getTargetNodes(ctx context.Context, ex *connExecutor) ([]system.NNode, error) {
	// If there is no from clause, this query can be sent to any node, if we already have a
	// connection to a node use that one, if not get a new one.
	tables, err := stmt.getTables(ctx, ex.SystemContext)
	if err != nil {
		return nil, err
	}
	stmt.tables = tables

	// If the query doesn't target any tables then the query should be able to be served by any node
	// Or if the query targets only global tables.
	if len(tables) == 0 || linq.From(tables).AllT(func(table system.NTable) bool {
		return table.TableType == system.NTableType_GLOBAL || table.TableType == system.NTableType_ACCOUNT
	}) {
		if node, err := ex.SystemContext.Nodes.GetFirstNode(system.AllNodes); err != nil {
			return nil, errors.New("no nodes available to serve this query.")
		} else {
			return []system.NNode{*node}, nil
		}
	}

	// Check to see if all of the tables are global or shard. If its a global table query it can be
	// directed to any node in the cluster.
	// tables := stmt.getTables()

	accounts, err := stmt.getAccountIDs(ctx, ex)
	if err != nil {
		return nil, err
	}

	if len(accounts) == 1 {
		nodes, err := ex.SystemContext.Accounts.GetNodesForAccount(accounts[0])
		if err != nil {
			return nil, err
		}
		rand.Seed(time.Now().Unix())
		return []system.NNode{nodes[rand.Intn(len(nodes))]}, nil
	} else if len(accounts) > 1 {
		return nil, errors.New("multi account queries are not supported at this time.")
	} else {
		// Realistically, we wouldn't automatically target every node in the cluster.
		// If This query targets a global or accounts table, then we can target a
		// single node. But if it targets a sharded table then we would want to gather
		// a complete list of account IDs and then generate a list of queries and nodes
		// where we could query there data for all of those accounts from the smallest
		// number of nodes possible.
		return ex.SystemContext.Nodes.GetNodes()
	}
}

func (stmt *SelectStatement) getTables(ctx context.Context, sctx *system.SContext) ([]system.NTable, error) {
	tableKeys := queryutil.GetTables(stmt.Statement)
	tables := make([]system.NTable, len(tableKeys))
	for i, tableName := range tableKeys {
		if table, err := sctx.Schema.GetTable(tableName); err != nil {
			return nil, err
		} else {
			if table == nil {
				return nil, errors.New(fmt.Sprintf("table [%s] does not exist", tableName))
			}
			tables[i] = *table
		}
	}
	return tables, nil
}

func (stmt *SelectStatement) getAccountIDs(ctx context.Context, ex *connExecutor) ([]uint64, error) {
	accountIds := make([]uint64, 0)
	for _, table := range stmt.tables {
		if table.TableType != system.NTableType_SHARD {
			continue
		}

		shardColumn := table.ShardKey.(*system.NTable_SKey).SKey.ColumnName
		ids := queryutil.FindAccountIds(stmt.Statement, shardColumn)
		accountIds = append(accountIds, ids...)
	}
	return accountIds, nil
}

func (stmt *SelectStatement) compilePlan(ctx context.Context, ex *connExecutor, nodes []system.NNode) ([]plan.NodeExecutionPlan, error) {
	if err := stmt.handleNoahFunctions(ctx, ex); err != nil {
		return nil, err
	}
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

func (stmt *SelectStatement) handleNoahFunctions(ctx context.Context, ex *connExecutor) error {
	newStmt, err := queryutil.ReplaceFunctionCalls(stmt.Statement, ex.GetDropInFunctions())
	if err != nil {
		return err
	}
	stmt.Statement = newStmt.(pg_query.SelectStmt)
	return nil
}
