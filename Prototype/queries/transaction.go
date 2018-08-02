package queries

import (
	"github.com/Ready-Stock/pg_query_go/nodes"
	"fmt"
	"github.com/Ready-Stock/Noah/Prototype/context"
)

func HandleTransaction(ctx *context.SessionContext, stmt pg_query.TransactionStmt) error {
	fmt.Printf("Preparing Transaction Query\n")
	j, _ := stmt.MarshalJSON()
	fmt.Println(string(j))
	switch stmt.Kind {
	case pg_query.TRANS_STMT_BEGIN, pg_query.TRANS_STMT_START:
		ctx.TransactionState = context.StateInTxn
	case pg_query.TRANS_STMT_COMMIT:
	case pg_query.TRANS_STMT_ROLLBACK:
	case pg_query.TRANS_STMT_SAVEPOINT:
	case pg_query.TRANS_STMT_RELEASE:
		ctx.TransactionState = context.StateNoTxn
	case pg_query.TRANS_STMT_ROLLBACK_TO:
	case pg_query.TRANS_STMT_PREPARE:
	case pg_query.TRANS_STMT_COMMIT_PREPARED:
	case pg_query.TRANS_STMT_ROLLBACK_PREPARED:
	}
	return nil
}