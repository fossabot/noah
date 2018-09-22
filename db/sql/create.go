package sql

import (
	"github.com/Ready-Stock/Noah/db/sql/plan"
	pg_query2 "github.com/Ready-Stock/pg_query_go"
	"github.com/Ready-Stock/pg_query_go/nodes"
)

type CreateStatement struct {
	Statement pg_query.CreateStmt
	IQueryStatement
}

func CreateCreateStatement(stmt pg_query.CreateStmt) *CreateStatement {
	return &CreateStatement{
		Statement: stmt,
	}
}

func (stmt *CreateStatement) Execute(ex *connExecutor, res RestrictedCommandResult) error {
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

func (stmt *CreateStatement) getTargetNodes(ex *connExecutor) ([]uint64, error) {
	// A create statement would want to target all live nodes.

	return []uint64{0, 1, 2, 3}, nil
}

func (stmt *CreateStatement) getAccountIDs() ([]uint64, error) {

	return nil, nil
}

func (stmt *CreateStatement) compilePlan(ex *connExecutor, nodes []uint64) ([]plan.NodeExecutionPlan, error) {
	plans := make([]plan.NodeExecutionPlan, len(nodes))
	deparsed, err := pg_query2.Deparse(stmt.Statement)
	if err != nil {
		ex.Error(err.Error())
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