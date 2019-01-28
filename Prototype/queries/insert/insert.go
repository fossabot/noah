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

package insert

import (
	"fmt"
	"github.com/kataras/go-errors"
	"github.com/readystock/noah/prototype/cluster"
	"github.com/readystock/noah/prototype/context"
	pgq "github.com/readystock/pg_query_go"
	"github.com/readystock/pg_query_go/nodes"
	"strconv"
)

var (
	errorInsertWithoutTransaction = errors.New("inserts can only be performed from within a transaction")
	errorCouldNotFindTable        = errors.New("could not find table (%s) in metadata")
	errorRelationIsNull           = errors.New("relation is null")
	errorNoAccountIDColumn        = errors.New("could not find column designating tenant_id")
	errorAccountIDInvalid         = errors.New("tenant_id (%d) is not valid")
)

type InsertStatement struct {
	Statement pg_query.InsertStmt
	Query     string
}

func CreateInsertStatment(stmt pg_query.InsertStmt, tree pgq.ParsetreeList) InsertStatement {
	return InsertStatement{
		Statement: stmt,
		Query:     tree.Query,
	}
}

func (stmt InsertStatement) HandleInsert(ctx *context.SessionContext) error {
	fmt.Printf("Preparing Insert Query\n")
	j, _ := stmt.Statement.MarshalJSON()
	fmt.Println(string(j))
	if ctx.TransactionState != context.StateInTxn {
		return errorInsertWithoutTransaction
	}

	if nodes, err := getTargetNodesForInsert(stmt.Statement); err != nil {
		return err
	} else {
		response := ctx.DistributeQuery(stmt.Query, nodes...)
		return ctx.HandleResponse(response)
	}
}

func getTargetNodesForInsert(stmt pg_query.InsertStmt) ([]int, error) {
	if global, err := getTargetTableIsGlobal(stmt); err != nil {
		return nil, err
	} else if global {
		if tenant, err := getTargetTableIsTenantTable(stmt); err != nil {
			return nil, err
		} else {
			if tenant {
				fmt.Printf("CREATING NEW TENANT!\n")
			}
			nodes := make([]int, 0)
			for _, n := range cluster.Nodes {
				nodes = append(nodes, n.NodeID)
			}
			return nodes, nil
		}
	} else {
		account_id_index := -1
		for i, res := range stmt.Cols.Items {
			col := res.(pg_query.ResTarget)
			if *col.Name == "account_id" {
				account_id_index = i
				break
			}
		}
		if account_id_index == -1 {
			return nil, errorNoAccountIDColumn
		} else {
			slct := stmt.SelectStmt.(pg_query.SelectStmt)
			if len(slct.ValuesLists[0]) == 0 {
				return nil, errors.New("unsupported insert values")
			} else {
				val := slct.ValuesLists[0][account_id_index].(pg_query.A_Const)
				idstr := ""
				switch valt := val.Val.(type) {
				case pg_query.Integer:
					idstr = strconv.FormatInt(valt.Ival, 10)
				case pg_query.String:
					idstr = valt.Str
				default:
					return nil, errors.New("unsupported value type for tenant_id")
				}
				if id, err := strconv.Atoi(idstr); err != nil {
					return nil, err
				} else {
					return getInsertNodesForAccountID(id)
				}
			}
		}
	}
}

func getInsertNodesForAccountID(account_id int) ([]int, error) {
	if account, ok := cluster.Accounts[account_id]; !ok {
		return nil, errorAccountIDInvalid.Format(account_id)
	} else {
		return account.NodeIDs, nil
	}
}

func getTargetTableIsGlobal(stmt pg_query.InsertStmt) (bool, error) {
	if stmt.Relation != nil && stmt.Relation.Relname != nil {
		if table, ok := cluster.Tables[*stmt.Relation.Relname]; !ok {
			return false, errorCouldNotFindTable.Format(*stmt.Relation.Relname)
		} else {
			return table.IsGlobal, nil
		}
	} else {
		return false, errorRelationIsNull
	}
}

func getTargetTableIsTenantTable(stmt pg_query.InsertStmt) (bool, error) {
	if stmt.Relation != nil && stmt.Relation.Relname != nil {
		if table, ok := cluster.Tables[*stmt.Relation.Relname]; !ok {
			return false, errorCouldNotFindTable.Format(*stmt.Relation.Relname)
		} else {
			return table.IsTenantTable, nil
		}
	} else {
		return false, errorRelationIsNull
	}
}
