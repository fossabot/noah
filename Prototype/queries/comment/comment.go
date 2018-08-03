package comment

import (
	"github.com/Ready-Stock/pg_query_go/nodes"
	"github.com/Ready-Stock/Noah/Prototype/context"
	pgq "github.com/Ready-Stock/pg_query_go"
	"fmt"
)

type CommentStatement struct {
	Statement pg_query.CommentStmt
	Query     string
}

func CreateCommentStatment(stmt pg_query.CommentStmt, tree pgq.ParsetreeList) CommentStatement {
	return CommentStatement{
		Statement: stmt,
		Query:     tree.Query,
	}
}

func (stmt CommentStatement) HandleComment(ctx *context.SessionContext) error {
	fmt.Printf("Preparing Comment Query\n")
	j, _ := stmt.Statement.MarshalJSON()
	fmt.Println(string(j))


	ids := stmt.getAllNodes()

	response := ctx.DistributeQuery(stmt.Query, ids...)
	return ctx.HandleResponse(response)
}