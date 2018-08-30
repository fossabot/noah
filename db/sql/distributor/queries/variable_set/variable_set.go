package variable_set

import (
	"github.com/Ready-Stock/Noah/db/sql/context"
	"github.com/Ready-Stock/Noah/db/sql/distributor/queries"
	"github.com/Ready-Stock/pg_query_go/nodes"
)

type VariableSetStatement struct {
	Statement pg_query.VariableSetStmt
	queries.IQueryStatement
}

func CreateVariableSetStatement(stmt pg_query.VariableSetStmt) *VariableSetStatement {
	return &VariableSetStatement{
		Statement: stmt,
	}
}

func (stmt *VariableSetStatement) Execute(ctx *context.NContext) error {
	target_nodes, err := stmt.getTargetNodes(ctx)
	if err != nil {
		return err
	}

	plans, err := stmt.compilePlan(ctx, target_nodes)
	if err != nil {
		return err
	}

	return ctx.ExecutePlans(plans)
}
