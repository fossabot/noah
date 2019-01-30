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

package sql

import (
	"context"
	"fmt"
	"github.com/kataras/go-errors"
	"github.com/readystock/noah/db/sql/plan"
	"github.com/readystock/noah/db/system"
	pq "github.com/readystock/pg_query_go/nodes"
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

func (stmt *TransactionStatement) Execute(ctx context.Context, ex *connExecutor, res RestrictedCommandResult, pinfo *plan.PlaceholderInfo) error {
	switch stmt.Statement.Kind {
	case pq.TRANS_STMT_BEGIN, pq.TRANS_STMT_START:
		return ex.BeginTransaction(ctx)
	case pq.TRANS_STMT_COMMIT:
		defer func() {
			ex.TransactionState = TransactionState_NONE
			ex.TransactionMode = TransactionMode_AutoCommit
		}()

		if ex.TransactionState != TransactionState_ENTERED {
			// If the user sent a begin but then hasn't sent anything else, then simply exit the
			// transaction.
			if ex.TransactionState == TransactionState_PRE {
				return nil
			}
			return errors.New("no transaction to commit")
		}
		ex.TransactionState = TransactionState_ENDING

		nodesInTransaction := 0
		for _, inTransaction := range ex.nodeTransactions {
			if inTransaction {
				nodesInTransaction++
			}
		}

		// if there are no nodes currently in a transaction then simply return null
		if nodesInTransaction == 0 {
			return nil
		} else if nodesInTransaction > 1 {
			if err := ex.PrepareTwoPhase(); err != nil {
				return err
			} else {
				return ex.CommitTwoPhase()
			}
		} else {
			return ex.Commit()
		}

	case pq.TRANS_STMT_ROLLBACK:
		defer func() {
			ex.nSync.Lock()
			defer ex.nSync.Unlock()
			ex.TransactionState = TransactionState_NONE
			ex.TransactionMode = TransactionMode_AutoCommit
			ex.nodeTransactions = map[uint64]bool{}
		}()

		if ex.TransactionState != TransactionState_ENTERED {
			// If the user sent a begin but then hasn't sent anything else, then simply exit the
			// transaction.
			if ex.TransactionState == TransactionState_PRE {
				return nil
			}

			return errors.New("no transaction to rollback")
		}

		ex.TransactionState = TransactionState_ENDING

		nodesInTransaction := 0
		for _, inTransaction := range ex.nodeTransactions {
			if inTransaction {
				nodesInTransaction++
			}
		}

		// if there are no nodes currently in a transaction then simply return null
		if nodesInTransaction == 0 {
			return nil
		} else {
			return ex.Rollback()
		}
	}
	return nil
}

func (stmt *TransactionStatement) getTargetNodes(ctx context.Context, ex *connExecutor) ([]system.NNode, error) {
	return nil, nil
}

func (stmt *TransactionStatement) compilePlan(ctx context.Context, ex *connExecutor, nodes []system.NNode) ([]plan.NodeExecutionPlan, error) {
	return nil, nil
}

func (ex *connExecutor) BeginTransaction(ctx context.Context) error {
	if ex.TransactionState != TransactionState_NONE {
		return errors.New(fmt.Sprintf("error cannot begin transaction at this time, current transaction state [%s]", ex.TransactionState.String()))
	}

	if id, err := ex.SystemContext.NewSnowflake(); err != nil {
		return err
	} else {
		ex.TransactionID = id
		ex.TransactionState = TransactionState_PRE
		ex.TransactionMode = TransactionMode_Manual
		return nil
	}
}
