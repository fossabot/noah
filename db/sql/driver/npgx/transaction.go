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
 */

package npgx

import (
	"bytes"
	"context"
	"fmt"
	"github.com/kataras/go-errors"
	"time"
)

// ErrTxCommitRollback occurs when an error has occurred in a transaction and
// Commit() is called. PostgreSQL accepts COMMIT on aborted transactions, but
// it is treated as ROLLBACK.
var ErrTxCommitRollback = errors.New("commit unexpectedly resulted in rollback")
var ErrTxClosed = errors.New("tx is closed")
var ErrTxInFailure = errors.New("tx failed")

type TransactionIsoLevel string

// Transaction isolation levels
const (
	Serializable    = TransactionIsoLevel("serializable")
	RepeatableRead  = TransactionIsoLevel("repeatable read")
	ReadCommitted   = TransactionIsoLevel("read committed")
	ReadUncommitted = TransactionIsoLevel("read uncommitted")
)

type TransactionAccessMode string

// Transaction access modes
const (
	ReadWrite = TransactionAccessMode("read write")
	ReadOnly  = TransactionAccessMode("read only")
)

type TransactionDeferrableMode string

// Transaction deferrable modes
const (
	Deferrable    = TransactionDeferrableMode("deferrable")
	NotDeferrable = TransactionDeferrableMode("not deferrable")
)

const (
	TransactionStatusInProgress              = 0
	TransactionStatusCommitFailure           = -1
	TransactionStatusRollbackFailure         = -2
	TransactionStatusInFailure               = -3
	TransactionStatusCommitSuccess           = 1
	TransactionStatusRollbackSuccess         = 2
	TransactionStatusTwoPhasePrepared        = 3
	TransactionStatusTwoPhasePreparedFailure = 4
)

// Tx represents a database transaction.
//
// All Tx methods return ErrTxClosed if Commit or Rollback has already been
// called on the Tx.
type Transaction struct {
	conn                  *Conn
	connPool              *ConnPool
	err                   error
	status                int8
	twoPhaseTransactionId uint64
}

type TransactionOptions struct {
	IsoLevel       TransactionIsoLevel
	AccessMode     TransactionAccessMode
	DeferrableMode TransactionDeferrableMode
}

// Begin starts a transaction with the default transaction mode for the
// current connection. To use a specific transaction mode see BeginEx.
func (c *Conn) Begin() (*Transaction, error) {
	return c.BeginEx(context.Background(), nil)
}

// BeginEx starts a transaction with txOptions determining the transaction
// mode. Unlike database/sql, the context only affects the begin command. i.e.
// there is no auto-rollback on context cancelation.
func (c *Conn) BeginEx(ctx context.Context, txOptions *TransactionOptions) (*Transaction, error) {
	_, err := c.ExecEx(ctx, txOptions.beginSQL(), nil)
	if err != nil {
		// begin should never fail unless there is an underlying connection issue or
		// a context timeout. In either case, the connection is possibly broken.
		c.die(errors.New("failed to begin transaction"))
		return nil, err
	}

	return &Transaction{conn: c}, nil
}

func (txOptions *TransactionOptions) beginSQL() string {
	if txOptions == nil {
		return "begin"
	}

	buf := &bytes.Buffer{}
	buf.WriteString("begin")
	if txOptions.IsoLevel != "" {
		fmt.Fprintf(buf, " isolation level %s", txOptions.IsoLevel)
	}
	if txOptions.AccessMode != "" {
		fmt.Fprintf(buf, " %s", txOptions.AccessMode)
	}
	if txOptions.DeferrableMode != "" {
		fmt.Fprintf(buf, " %s", txOptions.DeferrableMode)
	}

	return buf.String()
}

func (tx *Transaction) PrepareTwoPhase(ctx context.Context, transactionId uint64) error {
	if tx.status != TransactionStatusInProgress {
		return ErrTxClosed
	}
	commandTag, err := tx.conn.ExecEx(ctx, fmt.Sprintf("prepare transaction '%d';", transactionId), nil)

}

// Commit commits the transaction
func (tx *Transaction) Commit() error {
	return tx.CommitEx(context.Background())
}

// CommitEx commits the transaction with a context.
func (tx *Transaction) CommitEx(ctx context.Context) error {
	if tx.status != TransactionStatusInProgress {
		return ErrTxClosed
	}

	commandTag, err := tx.conn.ExecEx(ctx, "commit", nil)
	if err == nil && commandTag == "COMMIT" {
		tx.status = TransactionStatusCommitSuccess
	} else if err == nil && commandTag == "ROLLBACK" {
		tx.status = TransactionStatusCommitFailure
		tx.err = ErrTxCommitRollback
	} else {
		tx.status = TransactionStatusCommitFailure
		tx.err = err
		// A commit failure leaves the connection in an undefined state
		tx.conn.die(errors.New("commit failed"))
	}

	if tx.connPool != nil {
		tx.connPool.Release(tx.conn)
	}

	return tx.err
}

// Rollback rolls back the transaction. Rollback will return ErrTxClosed if the
// Tx is already closed, but is otherwise safe to call multiple times. Hence, a
// defer tx.Rollback() is safe even if tx.Commit() will be called first in a
// non-error condition.
func (tx *Transaction) Rollback() error {
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	return tx.RollbackEx(ctx)
}

// RollbackEx is the context version of Rollback
func (tx *Transaction) RollbackEx(ctx context.Context) error {
	if tx.status != TransactionStatusInProgress {
		return ErrTxClosed
	}

	_, tx.err = tx.conn.ExecEx(ctx, "rollback", nil)
	if tx.err == nil {
		tx.status = TransactionStatusRollbackSuccess
	} else {
		tx.status = TransactionStatusRollbackFailure
		// A rollback failure leaves the connection in an undefined state
		tx.conn.die(errors.New("rollback failed"))
	}

	if tx.connPool != nil {
		tx.connPool.Release(tx.conn)
	}

	return tx.err
}
