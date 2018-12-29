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
    "github.com/readystock/noah/db/sql/plan"
    "github.com/readystock/noah/db/system"
    pg_query2 "github.com/readystock/pg_query_go"
    "github.com/readystock/pg_query_go/nodes"
    "strings"
)

type VariableSetStatement struct {
    Statement pg_query.VariableSetStmt
    IQueryStatement
}

func CreateVariableSetStatement(stmt pg_query.VariableSetStmt) *VariableSetStatement {
    return &VariableSetStatement{
        Statement: stmt,
    }
}

func (stmt *VariableSetStatement) Execute(ex *connExecutor, res RestrictedCommandResult, pinfo *plan.PlaceholderInfo) error {
    if strings.HasPrefix(strings.ToLower(*stmt.Statement.Name), "noah") {
        settingName := strings.Replace(strings.ToLower(*stmt.Statement.Name), "noah.", "", 1)
        settingValue, err := pg_query2.DeparseValue(stmt.Statement.Args.Items[0].(pg_query.A_Const))
        if err != nil {
            return err
        }
        return ex.SystemContext.Settings.SetSetting(settingName, settingValue)
    }

    return nil
    target_nodes, err := stmt.getTargetNodes(ex)
    if err != nil {
        return err
    }

    plans, err := stmt.compilePlan(ex, target_nodes)
    if err != nil {
        return err
    }

    return ex.ExecutePlans(plans, res)
}

func (stmt *VariableSetStatement) getTargetNodes(ex *connExecutor) ([]system.NNode, error) {
    return nil, nil
    // return ex.GetNodesForAccountID(nil)
}

func (stmt *VariableSetStatement) compilePlan(ex *connExecutor, nodes []system.NNode) ([]plan.NodeExecutionPlan, error) {
    plans := make([]plan.NodeExecutionPlan, len(nodes))
    deparsed, err := pg_query.Deparse(stmt.Statement)
    if err != nil {
        return nil, err
    }
    for i := 0; i < len(plans); i++ {
        plans[i] = plan.NodeExecutionPlan{
            CompiledQuery: *deparsed,
            Node:          nodes[i],
        }
    }
    return plans, nil
}
