/*
 * Copyright (c) 2018 Ready Stock
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package transaction

import (
	"fmt"
	"github.com/kataras/go-errors"
	"github.com/readystock/noah/Prototype/context"
	"github.com/readystock/pg_query_go/nodes"
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
		nodesInTran := make([]int, 0)
		for _, n := range ctx.Nodes {
			if n.TransactionState == context.StateInTxn {
				nodesInTran = append(nodesInTran, n.NodeID)
			}
		}
		if len(nodesInTran) == 0 {
			return errors.New("cannot rollback transaction, no nodes have a pending transaction")
		} else {
			response := ctx.DoTxnOnAllNodes(context.TxnRollback, nodesInTran...)
			return ctx.HandleResponse(response)
		}
	}
}

func commit(ctx *context.SessionContext) error {
	if ctx.TransactionState != context.StateInTxn {
		return errors.New("cannot commit transaction, no transaction has been created")
	} else {
		nodesInTran := make([]int, 0)
		for _, n := range ctx.Nodes {
			if n.TransactionState == context.StateInTxn {
				nodesInTran = append(nodesInTran, n.NodeID)
			}
		}
		if len(nodesInTran) == 0 {
			return errors.New("cannot commit transaction, no nodes have a pending transaction")
		} else {
			response := ctx.DoTxnOnAllNodes(context.TxnCommit, nodesInTran...)
			return ctx.HandleResponse(response)
		}
	}
}
