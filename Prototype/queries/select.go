package queries

import (
	"github.com/Ready-Stock/pg_query_go/nodes"
	"fmt"
	"github.com/Ready-Stock/Noah/Prototype/cluster"
	"errors"
	data "github.com/Ready-Stock/Noah/Prototype/datums"
	"strconv"
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
				fmt.Printf("Query targets account id: %s \n", strconv.Itoa(*account_id))
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
		return parseAExpr(where)
	case pg_query.BoolExpr:
		for _, expr := range where.Args.Items {
			if aExprIsAccountID(expr.(pg_query.A_Expr)) {
				return parseAExpr(expr.(pg_query.A_Expr))
			}
		}
	}
	return nil, nil
}

func aExprIsAccountID(expr pg_query.A_Expr) (bool) {
	return expr.Lexpr.(pg_query.ColumnRef).Fields.Items[len(expr.Lexpr.(pg_query.ColumnRef).Fields.Items) - 1].(pg_query.String).Str == "account_id"
}

func parseAExpr(expr pg_query.A_Expr) (*int, error) {
	if expr.Lexpr.(pg_query.ColumnRef).Fields.Items[len(expr.Lexpr.(pg_query.ColumnRef).Fields.Items) - 1].(pg_query.String).Str != "account_id" {
		return nil, nil
	}
	if expr.Name.Items[0].(pg_query.String).Str != "=" {
		return nil, errors.New("tenant filtering only supports equal comparisons")
	}
	switch v := expr.Rexpr.(pg_query.A_Const).Val.(type) {
	case pg_query.Integer:
		in := strconv.FormatInt(v.Ival, 10)
		i, _ := strconv.Atoi(in)
		return &i, nil
	case pg_query.String:
		if i, err := strconv.Atoi(v.Str); err != nil {
			return nil, errors.New("cannot parse account id, not a valid integer")
		} else {
			return &i, nil
		}
	}
	return nil, nil
}


