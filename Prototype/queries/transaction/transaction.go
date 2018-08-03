package transaction

import (
	"github.com/Ready-Stock/pg_query_go/nodes"
	"fmt"
	"github.com/Ready-Stock/Noah/Prototype/context"
	"github.com/kataras/go-errors"
	"github.com/Ready-Stock/Noah/Prototype/distributor"
	"github.com/Ready-Stock/Noah/Prototype/datums"
	"github.com/Ready-Stock/Noah/Prototype/cluster"
)

func HandleTransaction(ctx *context.SessionContext, stmt pg_query.TransactionStmt) error {
	fmt.Printf("Preparing Transaction Query\n")
	j, _ := stmt.MarshalJSON()
	fmt.Println(string(j))
	switch stmt.Kind {
	case pg_query.TRANS_STMT_BEGIN, pg_query.TRANS_STMT_START:
		ctx.TransactionState = context.StateInTxn
	case pg_query.TRANS_STMT_COMMIT:
		return commit(ctx)
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
		nodesInTran := make([]datums.Node, 0)
		for _, n := range ctx.Nodes {
			if n.TransactionState == context.StateInTxn {
				nodesInTran = append(nodesInTran, cluster.Nodes[n.NodeID])
			}
		}
		if len(nodesInTran) == 0 {
			return errors.New("cannot rollback transaction, no nodes have a pending transaction")
		} else {
			responses := distributor.DistributeQuery("ROLLBACK;", nodesInTran...)
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
				return err
			}
		}
	}
}

func commit(ctx *context.SessionContext) error  {
	if ctx.TransactionState != context.StateInTxn {
		return errors.New("cannot commit transaction, no transaction has been created")
	} else {
		nodesInTran := make([]datums.Node, 0)
		for _, n := range ctx.Nodes {
			if n.TransactionState == context.StateInTxn {
				nodesInTran = append(nodesInTran, cluster.Nodes[n.NodeID])
			}
		}
		if len(nodesInTran) == 0 {
			return errors.New("cannot commit transaction, no nodes have a pending transaction")
		} else {
			responses := distributor.DistributeQuery("COMMIT;", nodesInTran...)
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
				return err
			}
		}
	}
}