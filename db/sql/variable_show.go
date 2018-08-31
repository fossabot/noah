package sql

import (
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