package sql

import (
	"github.com/Ready-Stock/pg_query_go/nodes"
)

type SelectStatement struct {
	Statement pg_query.SelectStmt
	IQueryStatement
}

func CreateSelectStatement(stmt pg_query.SelectStmt) *SelectStatement {
	return &SelectStatement{
		Statement: stmt,
	}
}

func (stmt *SelectStatement) Execute(ex *connExecutor, res RestrictedCommandResult) error {
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
