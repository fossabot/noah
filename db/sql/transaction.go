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

package sql

import (
	"github.com/Ready-Stock/Noah/db/sql/plan"
	"github.com/Ready-Stock/Noah/db/store"
	"github.com/Ready-Stock/Noah/db/system"
	pq "github.com/Ready-Stock/pg_query_go/nodes"
	"github.com/kataras/go-errors"
)

type TransactionStatement struct {
	Statement pq.TransactionStmt
	IQueryStatement
}

func CreateTransactionStatement(stmt pq.TransactionStmt) *TransactionStatement {
	return &TransactionStatement{
		Statement: stmt,
	}
}

func (stmt *TransactionStatement) Execute(ex *connExecutor, res RestrictedCommandResult) error {
	if stmt.Statement.Kind == pq.TRANS_STMT_BEGIN || stmt.Statement.Kind == pq.TRANS_STMT_START {
		ex.storeTransaction = store.BeginStoreUpdateTransaction(ex.SystemContext)
		return ex.BeginTransaction()
	}
	switch stmt.Statement.Kind {
	case pq.TRANS_STMT_COMMIT:
		if err := ex.PrepareTwoPhase(); err != nil {
			return err
		} else {
			if err := ex.storeTransaction.Commit(); err != nil {
				return err
			}
			return ex.CommitTwoPhase()
		}
	case pq.TRANS_STMT_ROLLBACK:
		if ex.TransactionStatus != NTXPreparedSuccess {
			return errors.New("no transaction to rollback")
		}
		ex.storeTransaction.Rollback()
		return ex.RollbackTwoPhase()
	}
	return nil
}

func (stmt *TransactionStatement) getTargetNodes(ex *connExecutor) ([]system.NNode, error) {
	return nil, nil
}

func (stmt *TransactionStatement) compilePlan(ex *connExecutor, nodes []system.NNode) ([]plan.NodeExecutionPlan, error) {
	return nil, nil
}


func (ex *connExecutor) BeginTransaction() error {
	if ex.TransactionStatus != NTXNoTransaction {
		return errors.New("error cannot begin transaction at this time")
	}
	if id, err := ex.SystemContext.Snowflake.NextID(); err != nil {
		return err
	} else {
		ex.TransactionID = id
		ex.TransactionStatus = NTXNotStarted
		return nil
	}
}