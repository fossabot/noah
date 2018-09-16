package sql

import (
	"github.com/Ready-Stock/Noah/db/sql/plan"
	pg_query2 "github.com/Ready-Stock/pg_query_go"
	"github.com/Ready-Stock/pg_query_go/nodes"
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

func (stmt *VariableShowStatement) Execute(ex *connExecutor, res RestrictedCommandResult) error {
	target_nodes, err := stmt.getTargetNodes(ex)
	if err != nil {
		return err
	}

	plans, err := stmt.compilePlan(ex, target_nodes)
	if err != nil {
		return err
	}

	return ex.ExecutePlans(plans)
}

func (stmt *VariableShowStatement) getTargetNodes(ex *connExecutor) ([]uint64, error) {
	return ex.GetNodesForAccountID(nil)
}

func (stmt *VariableShowStatement) compilePlan(ex *connExecutor, nodes []uint64) ([]plan.NodeExecutionPlan, error) {
	plans := make([]plan.NodeExecutionPlan, len(nodes))
	deparsed, err := pg_query2.Deparse(stmt.Statement)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(plans); i++ {
		plans[i] = plan.NodeExecutionPlan{
			CompiledQuery: *deparsed,
			NodeID:        nodes[i],
		}
	}
	return plans, nil
}
