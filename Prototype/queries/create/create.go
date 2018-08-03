package create

import (
	"github.com/Ready-Stock/pg_query_go/nodes"
	pgq "github.com/Ready-Stock/pg_query_go"
	"github.com/Ready-Stock/Noah/Prototype/context"
	"fmt"
	"github.com/Ready-Stock/Noah/Prototype/datums"
	"github.com/Ready-Stock/Noah/Prototype/cluster"
	"github.com/Ready-Stock/Noah/Prototype/distributor"
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

	nodes := stmt.getAllNodes()

	for _, node := range nodes {
		if cachedNode, ok := ctx.Nodes[node.NodeID]; !ok {
			ctx.Nodes[node.NodeID] = context.NodeContext{
				TransactionState: context.StateNoTxn,
				NodeID:           node.NodeID,
			}
			fmt.Printf("Node (%d) was has been added to this context \n", node.NodeID)
			if err := stmt.continueOrStartTransaction(ctx, ctx.Nodes[node.NodeID]); err != nil {
				return err
			}
		} else {
			fmt.Printf("Node (%d) is already in the current context \n", node.NodeID)
			if err := stmt.continueOrStartTransaction(ctx, cachedNode); err != nil {
				return err
			}
		}
	}

	responses := distributor.DistributeQuery(stmt.Query, nodes...)
	return stmt.handleQueryResponse(responses, nodes...)
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

func (stmt CreateStatement) getAllNodes() ([]datums.Node) {
	nodes := make([]datums.Node, 0)
	for _, node := range cluster.Nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

func (stmt CreateStatement) continueOrStartTransaction(ctx *context.SessionContext, node context.NodeContext) error {
	if node.TransactionState != context.StateInTxn {
		responses := distributor.DistributeQuery("BEGIN;", cluster.Nodes[node.NodeID])
		for _, response := range responses {
			if response.Error != nil {
				return response.Error
			}
		}
		fmt.Printf("Sent (BEGIN) to node (%d) successfully\n", node.NodeID)
		node.TransactionState = context.StateInTxn
		ctx.Nodes[node.NodeID] = node
		return nil
	} else {
		fmt.Printf("Node (%d) already in transaction\n", node.NodeID)
		return nil
	}
}

func (stmt CreateStatement) handleQueryResponse(responses []distributor.QueryResult, nodes ...datums.Node) error {
	success := true
	var err error
	for _, response := range responses {
		if response.Error != nil {
			success = false
			err = response.Error
			break
		}
	}
	if success {
		return nil
	} else {
		fmt.Printf("SENDING ROLLBACK TO NODES \n")
		finalize := distributor.DistributeQuery("ROLLBACK;", nodes...)
		for _, response := range finalize {
			if response.Error != nil {
				fmt.Printf("FAILED TO ROLLBACK ON NODE (%d)\n", response.NodeID)
				return response.Error
			}
		}
		fmt.Printf("SUCCESSFULLY ROLLED BACK ANY CHANGES ON NODES \n")
		return err
	}
}