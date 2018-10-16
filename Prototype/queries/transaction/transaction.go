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
 *
 * This application uses Open Source components. You can find the 
 * source code of their open source projects along with license 
 * information below. We acknowledge and are grateful to these 
 * developers for their contributions to open source.
 *
 * Project: CockroachDB https://github.com/cockroachdb/cockroach
 * Copyright 2018 The Cockroach Authors.
 * License (Apache License 2.0) https://github.com/cockroachdb/cockroach/blob/master/LICENSE
 *
 * Project: Vitess https://github.com/vitessio/vitess
 * Copyright 2018 Google Inc.
 * License (Apache License 2.0) https://github.com/vitessio/vitess/blob/master/LICENSE
 *
 * Project: Citus https://github.com/citusdata/citus
 * Copyright 2018 Citus Data, Inc.
 * License (GNU Affero General Public License v3.0) https://github.com/citusdata/citus/blob/master/LICENSE
 *
 * Project: pg_query_go https://github.com/lfittl/pg_query_go
 * Copyright 2018 Lukas Fittl
 * License (3-Clause BSD) https://github.com/lfittl/pg_query_go/blob/master/LICENSE
 *
 * Project: pgx https://github.com/jackc/pgx
 * Copyright 2018 Jack Christensen
 * License (MIT) https://github.com/jackc/pgx/blob/master/LICENSE
 *
 * Project: BadgerDB https://github.com/dgraph-io/badger
 * Copyright 2018 Dgraph Labs, Inc. and Contributors
 * License (MIT) https://github.com/dgraph-io/badger/blob/master/LICENSE
 *
 * Project: Sonyflake https://github.com/sony/sonyflake
 * Copyright 2018 Sony Corporation
 * License (MIT) https://github.com/sony/sonyflake/blob/master/LICENSE
 *
 * Project: Raft https://github.com/hashicorp/raft
 * Copyright 2018 HashiCorp
 * License (MPL-2.0) https://github.com/hashicorp/raft/blob/master/LICENSE
 */

package transaction

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

func commit(ctx *context.SessionContext) error  {
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