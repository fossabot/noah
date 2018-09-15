package sql

import (
	"github.com/Ready-Stock/pg_query_go/nodes"
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

func (stmt *VariableSetStatement) Execute(ex *connExecutor, res RestrictedCommandResult) error {
	return nil
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
