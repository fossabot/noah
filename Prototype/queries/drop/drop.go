package drop

import (
	"github.com/Ready-Stock/pg_query_go/nodes"
	"github.com/Ready-Stock/Noah/Prototype/context"
	pgq "github.com/Ready-Stock/pg_query_go"
	"fmt"
)

type DropStatement struct {
	Statement pg_query.DropStmt
	Query     string
}

func CreateDropStatment(stmt pg_query.DropStmt, tree pgq.ParsetreeList) DropStatement {
	return DropStatement{
		Statement: stmt,
		Query:     tree.Query,
	}
}

func (stmt DropStatement) HandleComment(ctx *context.SessionContext) error {
	fmt.Printf("Preparing Drop Query\n")
	j, _ := stmt.Statement.MarshalJSON()
	fmt.Println(string(j))
	ids := ctx.GetAllNodes()
	response := ctx.DistributeQuery(stmt.Query, ids...)
	return ctx.HandleResponse(response)
}
