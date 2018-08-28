package _select

import (
	"github.com/Ready-Stock/pg_query_go/nodes"
	pgq "github.com/Ready-Stock/pg_query_go"
	"github.com/Ready-Stock/Noah/db/sql/context"
	"github.com/Ready-Stock/Noah/db/sql/distributor/queries"
)

type SelectStatement struct {
	Statement pg_query.SelectStmt
	Query     string
	queries.IQueryStatement
}

func CreateSelectStatement(stmt pg_query.SelectStmt, tree pgq.ParsetreeList) *SelectStatement {
	return &SelectStatement{
		Statement: stmt,
		Query:     tree.Query,
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
