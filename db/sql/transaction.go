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

func (stmt *TransactionStatement) Execute(ex *connExecutor, res RestrictedCommandResult, pinfo *plan.PlaceholderInfo) error {
	switch stmt.Statement.Kind {
	case pq.TRANS_STMT_BEGIN, pq.TRANS_STMT_START:
		return ex.BeginTransaction()
	case pq.TRANS_STMT_COMMIT:
		if err := ex.PrepareTwoPhase(); err != nil {
			return err
		} else {
			return ex.CommitTwoPhase()
		}
	case pq.TRANS_STMT_ROLLBACK:
		if ex.TransactionState != TransactionState_ENTERED {
			return errors.New("no transaction to rollback")
		}
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
	if ex.TransactionState != TransactionState_NONE {
		return errors.New("error cannot begin transaction at this time")
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
