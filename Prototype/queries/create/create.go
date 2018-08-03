package create

import (
	"github.com/Ready-Stock/pg_query_go/nodes"
	pgq "github.com/Ready-Stock/pg_query_go"
	"github.com/Ready-Stock/Noah/Prototype/context"
	"fmt"
		"github.com/Ready-Stock/Noah/Prototype/cluster"
		)

type CreateStatement struct {
	Statement pg_query.CreateStmt
	Query     string
}

func CreateCreateStatment(stmt pg_query.CreateStmt, tree pgq.ParsetreeList) CreateStatement {
	return CreateStatement{
		Statement: stmt,
		Query:     tree.Query,
	}
}

func (stmt CreateStatement) HandleCreate(ctx *context.SessionContext) error {
	fmt.Printf("Preparing Create Query\n")
	j, _ := stmt.Statement.MarshalJSON()
	fmt.Println(string(j))

	var has_id bool
	var id_index int
	var id_name string

	if has_id, id_index, id_name = stmt.hasIdentityColumn(); has_id {
		fmt.Printf("Found identity column at index: (%d) name: (%s)\n", id_index, id_name)
	}

	ids := stmt.getAllNodes()

	response := ctx.DistributeQuery(stmt.Query, ids...)
	return ctx.HandleResponse(response)
}

func (stmt CreateStatement) hasIdentityColumn() (bool, int, string) {
	identity_index := -1
	for i, telt := range stmt.Statement.TableElts.Items {
		col := telt.(pg_query.ColumnDef)
		if col.Constraints.Items != nil && len(col.Constraints.Items) > 0 {
			for _, c := range col.Constraints.Items {
				constraint := c.(pg_query.Constraint)
				if constraint.Contype == pg_query.CONSTR_PRIMARY {
					identity_index = i
					return true, identity_index, *col.Colname
				}
			}
		}
	}
	return false, identity_index, ""
}

func (stmt CreateStatement) getAllNodes() ([]int) {
	ids := make([]int, 0)
	for _, node := range cluster.Nodes {
		ids = append(ids, node.NodeID)
	}
	return ids
}
