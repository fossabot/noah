package queries

import (
	"github.com/Ready-Stock/pg_query_go/nodes"
	"fmt"
	"github.com/Ready-Stock/Noah/Prototype/context"
	"github.com/kataras/go-errors"
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
		return rollback(ctx)
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

func rollback(ctx *context.SessionContext) error {
	if ctx.TransactionState != context.StateInTxn {
		return errors.New("cannot rollback transaction, no transaction has been created")
	} else {
		nodesInTran := make([]context.NodeContext, 0)
		for _, n := range ctx.Nodes {
			if n.TransactionState == context.StateInTxn {
				nodesInTran = append(nodesInTran, n)
			}
		}
		if len(nodesInTran) == 0 {
			return errors.New("cannot rollback transaction, no nodes have a pending transaction")
		} else {
			for _, node := range nodesInTran {
				fmt.Printf("Sent (ROLLBACK) to node (%d)\n", node.NodeID)
			}
			return nil
		}
	}
}

func commit(ctx *context.SessionContext) error  {
	
}