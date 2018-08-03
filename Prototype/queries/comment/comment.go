package comment

import (
	"github.com/Ready-Stock/pg_query_go/nodes"
	"github.com/Ready-Stock/Noah/Prototype/context"
	pgq "github.com/Ready-Stock/pg_query_go"
	"fmt"
	"github.com/Ready-Stock/Noah/Prototype/cluster"
	"github.com/kataras/go-errors"
)

var (
	errorTableNotFound = errors.New("table (%s) cannot be found in metadata")
	errorTableTypeNotValid = errors.New("table type (%s) is not valid")
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

	if stmt.Statement.Object != nil && len(stmt.Statement.Object.(pg_query.List).Items) > 0 {
		lst := stmt.Statement.Object.(pg_query.List)
		table := lst.Items[len(lst.Items) - 1].(pg_query.String).Str
		if tblmeta, ok := cluster.Tables[table]; !ok {
			return errorTableNotFound.Format(table)
		} else {
			switch *stmt.Statement.Comment {
			case "GLOBAL":
				tblmeta.IsGlobal = true
				tblmeta.IsTenantTable = false
			case "TENANTS":
				tblmeta.IsGlobal = true
				tblmeta.IsTenantTable = true
			case "LOCAL":
				tblmeta.IsGlobal = false
				tblmeta.IsTenantTable = false
			default:
				return errorTableTypeNotValid.Format(*stmt.Statement.Comment)
			}
			cluster.Tables[table] = tblmeta
		}
	}
	ids := ctx.GetAllNodes()
	response := ctx.DistributeQuery(stmt.Query, ids...)
	return ctx.HandleResponse(response)
}