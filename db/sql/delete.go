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
	"fmt"
	"github.com/kataras/go-errors"
	"github.com/readystock/noah/db/sql/plan"
	"github.com/readystock/noah/db/system"
	"github.com/readystock/pg_query_go/nodes"
)

type DeleteStatement struct {
	Statement pg_query.DeleteStmt
	IQueryStatement
}

func CreateDeleteStatement(stmt pg_query.DeleteStmt) *DeleteStatement {
	return &DeleteStatement{
		Statement: stmt,
	}
}

func (stmt *DeleteStatement) Execute(ex *connExecutor, res RestrictedCommandResult, pinfo *plan.PlaceholderInfo) error {

	return nil
}

func (stmt *DeleteStatement) getTargetNodes(ex *connExecutor) ([]system.NNode, error) {
	tableName := *stmt.Statement.Relation.Relname
	table, err := ex.SystemContext.Schema.GetTable(tableName)
	if err != nil {
		return nil, err
	}

	if table == nil {
		return nil, errors.New(fmt.Sprintf("table [%s] does not exist", tableName))
	}

	if stmt.Statement.WhereClause == nil {
		// If there is no where clause then we are deleting everything from the table.
		// In the future we will want to add something to restrict this and allow
		// users to configure this.
		_, err := ex.SystemContext.Nodes.GetNodes()
		if err != nil {
			return nil, err
		}

	}

	return nil, nil
}

func (stmt *DeleteStatement) compilePlan(ex *connExecutor, nodes []system.NNode) ([]plan.NodeExecutionPlan, error) {
	plans := make([]plan.NodeExecutionPlan, len(nodes))
	deparsed, err := stmt.Statement.Deparse(pg_query.Context_None)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(plans); i++ {
		plans[i] = plan.NodeExecutionPlan{
			CompiledQuery: *deparsed,
			Node:          nodes[i],
			ReadOnly:      false,
			Type:          stmt.Statement.StatementType(),
		}
	}
	return plans, nil
}
