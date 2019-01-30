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
	"github.com/pkg/errors"
	"github.com/readystock/golinq"
	"github.com/readystock/golog"
	"github.com/readystock/noah/db/sql/plan"
	"github.com/readystock/noah/db/system"
	"github.com/readystock/pg_query_go/nodes"
)

type DropStatement struct {
	Statement pg_query.DropStmt
	IQueryStatement
}

func CreateDropStatement(stmt pg_query.DropStmt) *DropStatement {
	return &DropStatement{
		Statement: stmt,
	}
}

func (stmt *DropStatement) Execute(ctx context.Context, ex *connExecutor, res RestrictedCommandResult, pinfo *plan.PlaceholderInfo) error {
	switch stmt.Statement.RemoveType {
	case pg_query.OBJECT_TABLE:
		targetNodes, err := stmt.getTargetNodes(ctx, ex)
		if err != nil {
			return err
		}
		golog.Debugf("Preparing to send query to %d node(s)", len(targetNodes))

		plans, err := stmt.compilePlan(ctx, ex, targetNodes)
		if err != nil {
			return err
		}

		if err := ex.ExecutePlans(ctx, plans, res); err != nil {
			return err
		} else {
			inner := stmt.Statement.Objects.Items[0].(pg_query.List).Items
			tableName := inner[len(inner)-1].(pg_query.String).Str
			return ex.SystemContext.Schema.DropTable(tableName)
		}
	default:
		return errors.New(fmt.Sprintf("cannot drop object [%v], not yet supported", stmt.Statement.RemoveType))
	}
}

func (stmt *DropStatement) compilePlan(ctx context.Context, ex *connExecutor, nodes []system.NNode) ([]plan.NodeExecutionPlan, error) {
	plans := make([]plan.NodeExecutionPlan, len(nodes))
	compiled, err := pg_query.Deparse(stmt.Statement)
	if err != nil {
		golog.Error(err.Error())
		return nil, err
	}
	golog.Debugf("Recompiled query: %s", *compiled)
	for i := 0; i < len(plans); i++ {
		plans[i] = plan.NodeExecutionPlan{
			CompiledQuery: *compiled,
			Node:          nodes[i],
			ReadOnly:      false,
			Type:          stmt.Statement.StatementType(),
		}
	}
	return plans, nil
}

func (stmt *DropStatement) getTargetNodes(ctx context.Context, ex *connExecutor) ([]system.NNode, error) {
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
