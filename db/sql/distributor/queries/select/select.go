package _select

import (
	"github.com/Ready-Stock/Noah/db/sql/context"
	"github.com/Ready-Stock/Noah/db/sql/distributor/queries"
	"github.com/Ready-Stock/pg_query_go/nodes"
)

type SelectStatement struct {
	Statement pg_query.SelectStmt
	queries.IQueryStatement
}

func CreateSelectStatement(stmt pg_query.SelectStmt) *SelectStatement {
	return &SelectStatement{
		Statement: stmt,
	}
}

func (stmt *SelectStatement) Execute(ctx *context.NContext) error {
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
