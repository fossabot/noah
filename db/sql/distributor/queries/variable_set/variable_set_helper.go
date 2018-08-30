package variable_set

import (
	"github.com/Ready-Stock/Noah/db/sql/context"
	"github.com/Ready-Stock/pg_query_go"
)

func (stmt *VariableSetStatement) getTargetNodes(ctx *context.NContext) ([]int, error) {
	return ctx.GetNodesForAccountID(nil)
}

func (stmt *VariableSetStatement) compilePlan(ctx *context.NContext, nodes []int) ([]context.NodeExecutionPlan, error) {
	plans := make([]context.NodeExecutionPlan, len(nodes))
	deparsed, err := pg_query.Deparse(stmt.Statement)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(plans); i++ {
		plans[i] = context.NodeExecutionPlan{
			CompiledQuery: *deparsed,
			NodeID: nodes[i],
		}
	}
	return plans, nil
}
