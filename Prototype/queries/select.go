package queries

import (
	"github.com/Ready-Stock/pg_query_go/nodes"
	"fmt"
	"github.com/Ready-Stock/Noah/Prototype/cluster"
	"errors"
	data "github.com/Ready-Stock/Noah/Prototype/datums"
)

func HandleSelect(stmt pg_query.SelectStmt) error {
	fmt.Printf("Preparing Select Query\n")
	j, _ := stmt.MarshalJSON()
	fmt.Println(string(j))
	if node, err := getTargetNodeForQuery(stmt); err != nil {
		return err
	} else if node == nil {
		return errors.New("no node targeted for select query")
	} else {
		fmt.Printf("Sending Query To Node ID (%d)", node.NodeID)
	}

	return nil
}

func getTargetNodeForQuery(stmt pg_query.SelectStmt) (*data.Node, error) {
	var node data.Node
	if len(stmt.FromClause.Items) == 0 {
		fmt.Printf("Query does not have a from clause!\n")
		node = cluster.GetRandomNode()
		return &node, nil
	} else {
		table := stmt.FromClause.Items[0].(pg_query.RangeVar)
		if _, ok := cluster.Tables[*table.Relname]; !ok {
			return nil, errors.New("cannot find table (" + *table.Relname + ")")
		} else {
			if account_id, err := getAccountIDForQuery(stmt); err != nil {
				return nil, err
			} else if account_id == nil {
				return nil, errors.New("select query targets all accounts, query is not safe")
			} else {
				node = cluster.GetNodeForAccount(*account_id)
				return &node, nil
			}
		}
	}
}

func getAccountIDForQuery(stmt pg_query.SelectStmt) (*int, error) {
	if stmt.WhereClause == nil {
		return nil, nil
	}
	switch where := stmt.WhereClause.(type) {
	case pg_query.A_Expr:

	case pg_query.BoolExpr:

	}
	return nil, nil
}


