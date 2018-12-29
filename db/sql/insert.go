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
    "github.com/kataras/go-errors"
    "github.com/readystock/noah/db/sql/plan"
    "github.com/readystock/noah/db/system"
    "github.com/readystock/pg_query_go/nodes"
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
    accounts, err := stmt.getAccountIDs()
    if err != nil {
        return nil, err
    }

    if len(accounts) == 1 {
        return ex.SystemContext.Accounts.GetNodesForAccount(accounts[0])
    } else if len(accounts) > 1 {
        return nil, errors.New("multi account queries are not supported at this time.")
        // node_ids := make([]uint64, 0)
        // From(accounts).SelectManyT(func(id uint64) Query {
        //     if ids, err := ex.GetNodesForAccountID(&id); err == nil {
        //         return From(ids)
        //     }
        //     return From(make([]uint64, 0))
        // }).Distinct().ToSlice(&node_ids)
        // if len(node_ids) == 0 {
        //     return nil, errors.New("could not find nodes for account IDs")
        // } else {
        //     return node_ids, nil
        // }
    } else {
        return ex.SystemContext.Nodes.GetNodes()
    }
}

func (stmt *InsertStatement) getAccountIDs() ([]uint64, error) {
    return make([]uint64, 0), nil
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
        }
    }
    return plans, nil
}
