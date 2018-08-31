package sql

import (
	"github.com/Ready-Stock/Noah/db/sql/plan"
	"github.com/Ready-Stock/pg_query_go"
)

func (stmt *VariableShowStatement) getTargetNodes(ex *connExecutor) ([]int, error) {
	return ex.GetNodesForAccountID(nil)
}

func (stmt *VariableShowStatement) compilePlan(ex *connExecutor, nodes []int) ([]plan.NodeExecutionPlan, error) {
	plans := make([]plan.NodeExecutionPlan, len(nodes))
	deparsed, err := pg_query.Deparse(stmt.Statement)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(plans); i++ {
		plans[i] = plan.NodeExecutionPlan{
			CompiledQuery: *deparsed,
			NodeID: nodes[i],
		}
	}
	return plans, nil
}
