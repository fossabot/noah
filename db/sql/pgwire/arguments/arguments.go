package arguments

import (
	"github.com/readystock/pg_query_go/nodes"
)

func GetArguments(stmt pg_query.Stmt) (numArgs int) {
	switch query := stmt.(type) {
	case pg_query.SelectStmt:
		return getArgumentsFromSelect(query)
	}
	return 0
}
