/*
 * Copyright (c) 2018 Ready Stock
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package create

import (
	"fmt"
	"github.com/readystock/noah/prototype/cluster"
	"github.com/readystock/noah/prototype/context"
	"github.com/readystock/noah/prototype/datums"
	pgq "github.com/readystock/pg_query_go"
	"github.com/readystock/pg_query_go/nodes"
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

	ids := ctx.GetAllNodes()

	response := ctx.DistributeQuery(stmt.Query, ids...)

	if response.Success {
		cluster.Tables[*stmt.Statement.Relation.Relname] = datums.Table{
			TableName:      *stmt.Statement.Relation.Relname,
			IsGlobal:       false,
			IsTenantTable:  false,
			IdentityColumn: id_name,
		}
	}

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
