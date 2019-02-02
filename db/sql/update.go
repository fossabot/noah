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

	"github.com/readystock/noah/db/sql/plan"
	"github.com/readystock/noah/db/system"
	"github.com/readystock/noah/db/util/queryutil"
	pg_query "github.com/readystock/pg_query_go/nodes"
)

type UpdateStatement struct {
	Statement pg_query.UpdateStmt
	tables    []system.NTable
	IQueryStatement
}

func CreateUpdateStatement(stmt pg_query.UpdateStmt) *UpdateStatement {
	return &UpdateStatement{
		Statement: stmt,
	}
}

func (stmt *UpdateStatement) Execute(ctx context.Context, ex *connExecutor, res RestrictedCommandResult, pinfo *plan.PlaceholderInfo) error {
	nodes, err := stmt.getTargetNodes(ctx, ex)
	if err != nil {
		return err
	}

	plans, err := stmt.compilePlan(ctx, ex, nodes)
	if err != nil {
		return err
	}

	return ex.ExecutePlans(ctx, plans, res)
}

func (stmt *UpdateStatement) compilePlan(ctx context.Context, ex *connExecutor, nodes []system.NNode) ([]plan.NodeExecutionPlan, error) {

	return nil, nil
}

func (stmt *UpdateStatement) getTargetNodes(ctx context.Context, ex *connExecutor) ([]system.NNode, error) {
	// tables, err := stmt.getTables(ctx, ex)
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

func (stmt *UpdateStatement) getTables(ctx context.Context, ex *connExecutor) ([]system.NTable, error) {
	tableKeys := queryutil.GetTables(stmt.Statement)
	tables := make([]system.NTable, len(tableKeys))
	for i, tableName := range tableKeys {
		if table, err := ex.SystemContext.Schema.GetTable(tableName); err != nil {
			return nil, err
		} else {
			if table == nil {
				return nil, fmt.Errorf("table [%s] does not exist", tableName)
			}
			tables[i] = *table
		}
	}
	return tables, nil
}
