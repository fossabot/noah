package queries

import (
	"github.com/Ready-Stock/pg_query_go/nodes"
	"github.com/Ready-Stock/Noah/Prototype/context"
	"fmt"
	data "github.com/Ready-Stock/Noah/Prototype/datums"
	"github.com/Ready-Stock/Noah/Prototype/cluster"
	"github.com/kataras/go-errors"
	"strconv"
)

var (
	errorInsertWithoutTransaction = errors.New("inserts can only be performed from within a transaction")
	errorCouldNotFindTable        = errors.New("could not find table (%s) in metadata")
	errorRelationIsNull			  = errors.New("relation is null")
	errorNoAccountIDColumn		  = errors.New("could not find column designating tenant_id")
	errorAccountIDInvalid 		  = errors.New("tenant_id (%d) is not valid")
)

func HandleInsert(ctx *context.SessionContext, stmt pg_query.InsertStmt) error {
	fmt.Printf("Preparing Insert Query\n")
	j, _ := stmt.MarshalJSON()
	fmt.Println(string(j))
	if ctx.TransactionState != context.StateInTxn {
		return errorInsertWithoutTransaction
	}

	if nodes, err := getTargetNodesForInsert(stmt); err != nil {
		return err
	} else {
		for _, node := range nodes {
			if cachedNode, ok := ctx.Nodes[node.NodeID]; !ok {
				ctx.Nodes[node.NodeID] = context.NodeContext{
					TransactionState: context.StateNoTxn,
					NodeID:           node.NodeID,
				}
				fmt.Printf("Node (%d) was has been added to this context \n", node.NodeID)
				if err := continueOrStartTransaction(ctx, ctx.Nodes[node.NodeID]); err != nil {
					return err
				}
			} else {
				fmt.Printf("Node (%d) is already in the current context \n", node.NodeID)
				if err := continueOrStartTransaction(ctx, cachedNode); err != nil {
					return err
				}
			}

		}
		for _, node := range nodes {
			fmt.Printf("Sending insert to node (%d) \n", node.NodeID)
		}
	}

	return nil
}

func getTargetNodesForInsert(stmt pg_query.InsertStmt) ([]data.Node, error) {
	if global, err := getTargetTableIsGlobal(stmt); err != nil {
		return nil, err
	} else if global {
		nodes := make([]data.Node, len(cluster.Nodes))
		for i, n := range cluster.Nodes {
			nodes[i] = n
		}
		return nodes, nil
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

func getInsertNodesForAccountID(account_id int) ([]data.Node, error) {
	if account, ok := cluster.Accounts[account_id]; !ok {
		return nil, errorAccountIDInvalid.Format(account_id)
	} else {
		nodes := make([]data.Node, len(account.NodeIDs))
		for i, nid := range account.NodeIDs {
			nodes[i] = cluster.Nodes[nid]
		}
		return nodes, nil
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

func continueOrStartTransaction(ctx *context.SessionContext, node context.NodeContext) error {
	if node.TransactionState != context.StateInTxn {
		fmt.Printf("Sent (BEGIN) to node (%d)\n", node.NodeID)
		node.TransactionState = context.StateInTxn
		ctx.Nodes[node.NodeID] = node
		return nil
	} else {
		fmt.Printf("Node (%d) already in transaction\n", node.NodeID)
		return nil
	}
}
