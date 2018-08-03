package _select

import (
	"github.com/Ready-Stock/pg_query_go/nodes"
	"fmt"
	"github.com/Ready-Stock/Noah/Prototype/cluster"
	"errors"
	pgq "github.com/Ready-Stock/pg_query_go"
	data "github.com/Ready-Stock/Noah/Prototype/datums"
	"strconv"
	"github.com/Ready-Stock/Noah/Prototype/context"
	"database/sql"
)

type SelectStatement struct {
	Statement pg_query.SelectStmt
	Query     string
}

func CreateSelectStatement(stmt pg_query.SelectStmt, tree pgq.ParsetreeList) SelectStatement {
	return SelectStatement{
		Statement: stmt,
		Query:     tree.Query,
	}
}

func (stmt SelectStatement) HandleSelect(ctx *context.SessionContext) (*sql.Rows, error) {
	fmt.Printf("Preparing Select Query\n")
	j, _ := stmt.Statement.MarshalJSON()
	fmt.Println(string(j))
	if node, err := getTargetNodeForSelect(stmt.Statement); err != nil {
		return nil, err
	} else if node == nil {
		return nil, errors.New("no node targeted for select query")
	} else {
		fmt.Printf("Sending Query To Node ID (%d)\n", node.NodeID)
		response := ctx.DistributeQuery(stmt.Query, node.NodeID)
		if response.Errors
	}

	return nil, nil
}

func getTargetNodeForSelect(stmt pg_query.SelectStmt) (*data.Node, error) {
	var node data.Node
	if len(stmt.FromClause.Items) == 0 {
		switch stmt.TargetList.Items[0].(pg_query.ResTarget).Val.(type) {
		case pg_query.FuncCall:
			return nil, errors.New("function calls are not supported without a from clause")
		default:
			fmt.Printf("Query does not have a from clause!\n")
			node = cluster.GetRandomNode()
			return &node, nil
		}
	} else {
		if is_global, err := getTableIsGlobal(stmt); err != nil {
			return nil, err
		} else if is_global {
			node = cluster.GetRandomNode()
			return &node, nil
		} else {
			if account_id, err := getAccountIDForQuery(&stmt.WhereClause); err != nil {
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



func getAccountIDForQuery(exp *pg_query.Node) (*int, error) {
	if exp == nil {
		return nil, nil
	} else {
		return parseExpression(*exp)
	}
}

func parseExpression(filter pg_query.Node) (*int, error) {
	switch where := filter.(type) {
	case pg_query.A_Expr:
		return parseAExpr(where)
	case pg_query.BoolExpr:
		for _, expr := range where.Args.Items {
			switch filterT := expr.(type) {
			case pg_query.A_Expr:
				if aExprIsAccountID(filterT) {
					return parseAExpr(filterT)
				}
			case pg_query.BoolExpr:
				if id, err := parseExpression(filterT); err != nil {
					return nil, err
				} else if id != nil {
					return id, nil
				}
			default:
				return nil, errors.New("unknown expression type")
			}
		}
	default:
		return nil, errors.New("could not parse where filter")
	}
	return nil, nil
}



func aExprIsAccountID(expr pg_query.A_Expr) (bool) {
	return expr.Lexpr.(pg_query.ColumnRef).Fields.Items[len(expr.Lexpr.(pg_query.ColumnRef).Fields.Items)-1].(pg_query.String).Str == "account_id"
}

func parseAExpr(expr pg_query.A_Expr) (*int, error) {
	if expr.Lexpr.(pg_query.ColumnRef).Fields.Items[len(expr.Lexpr.(pg_query.ColumnRef).Fields.Items)-1].(pg_query.String).Str != "account_id" {
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

func getTableIsGlobal(stmt pg_query.SelectStmt) (bool, error) {
	switch from := stmt.FromClause.Items[0].(type) {
	case pg_query.RangeVar:
		if _, ok := cluster.Tables[*from.Relname]; !ok {
			return false, errors.New("cannot find table (" + *from.Relname + ")")
		} else {
			return true, nil
		}
	case pg_query.JoinExpr:
		if tables, err := getTablesFromJoin(from); err != nil {
			return false, err
		} else {
			for _, table := range tables {
				if tablet, ok := cluster.Tables[table]; !ok {
					return false, errors.New("table (" + table + ") is not in metadata.")
				} else {
					if !tablet.IsGlobal {
						return false, nil
					}
				}
			}
			return false, nil // If it gets this far then that means that all of the tables in the query are global tables and can be served by any node
		}
	default:
		return false, errors.New("unrecognized from clause type")
	}
}

func getTablesFromJoin(join pg_query.JoinExpr) ([]string, error) {
	tables := make([]string, 0)
	parseJoinArg := func(node pg_query.Node) ([]string, error) {
		switch arg := node.(type) {
		case pg_query.RangeVar:
			return []string{ *arg.Relname }, nil
		case pg_query.JoinExpr:
			return getTablesFromJoin(arg)
		default:
			return nil, errors.New("invalid node type")
		}
	}
	if ltab, err := parseJoinArg(join.Larg); err != nil {
		return nil, err
	} else {
		tables = append(tables, ltab...)
	}
	if rtab, err := parseJoinArg(join.Rarg); err != nil {
		return nil, err
	} else {
		tables = append(tables, rtab...)
	}
	return tables, nil
}
