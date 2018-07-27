package pglex

import (
	node "github.com/Ready-Stock/pg_query_go/nodes"
	"fmt"
)

func HandleRawStmt(stmt node.RawStmt) error {
	switch tstmt := stmt.Stmt.(type) {
	case node.VariableSetStmt:
		fmt.Sprintf("Changing Setting (%s)\n", tstmt.Name)
	}
	return nil
}