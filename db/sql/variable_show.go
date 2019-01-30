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
	"context"
	"github.com/readystock/noah/db/sql/pgwire/pgproto"
	"github.com/readystock/noah/db/sql/plan"
	"github.com/readystock/noah/db/sql/types"
	"github.com/readystock/noah/db/system"
	"github.com/readystock/pg_query_go/nodes"
	"strings"
)

type VariableShowStatement struct {
	Statement pg_query.VariableShowStmt
	IQueryStatement
}

func CreateVariableShowStatement(stmt pg_query.VariableShowStmt) *VariableShowStatement {
	return &VariableShowStatement{
		Statement: stmt,
	}
}

func (stmt *VariableShowStatement) Execute(ctx context.Context, ex *connExecutor, res RestrictedCommandResult, pinfo *plan.PlaceholderInfo) error {
	name := strings.ToLower(*stmt.Statement.Name)
	if strings.HasPrefix(name, "noah") {
		settingName := strings.Replace(name, "noah.", "", 1)
		value, err := ex.SystemContext.Settings.GetSetting(system.NoahSetting(settingName))
		if err != nil {
			return err
		}
		columns := []pgproto.FieldDescription{
			{
				Name:                 settingName,
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          25,
				DataTypeSize:         int16(len(*value)),
				TypeModifier:         0,
				Format:               0,
			},
		}
		res.SetColumns(columns)
		values := []types.Value{
			&types.Text{
				String: *value,
				Status: types.Present,
			},
		}
		res.AddRow(values)
		return nil
	}

	targetNodes, err := stmt.getTargetNodes(ctx, ex)
	if err != nil {
		return err
	}

	plans, err := stmt.compilePlan(ctx, ex, targetNodes)
	if err != nil {
		return err
	}

	return ex.ExecutePlans(ctx, plans, res)
}

func (stmt *VariableShowStatement) getTargetNodes(ctx context.Context, ex *connExecutor) ([]system.NNode, error) {
	return nil, nil
	// return ex.GetNodesForAccountID(nil)
}

func (stmt *VariableShowStatement) compilePlan(ctx context.Context, ex *connExecutor, nodes []system.NNode) ([]plan.NodeExecutionPlan, error) {
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
