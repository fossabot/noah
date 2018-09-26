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

package driver

import (
	"github.com/kataras/go-errors"
)

type TransactionIsolationLevel string

// Transaction isolation levels
const (
	Serializable    = TransactionIsolationLevel("serializable")
	RepeatableRead  = TransactionIsolationLevel("repeatable read")
	ReadCommitted   = TransactionIsolationLevel("read committed")
	ReadUncommitted = TransactionIsolationLevel("read uncommitted")
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
	TransactionStatusInProgress      = 0
	TransactionStatusCommitFailure   = -1
	TransactionStatusRollbackFailure = -2
	TransactionStatusInFailure       = -3
	TransactionStatusCommitSuccess   = 1
	TransactionStatusRollbackSuccess = 2
)

type TransactionOptions struct {
	IsoLevel       TransactionIsolationLevel
	AccessMode     TransactionAccessMode
	DeferrableMode TransactionDeferrableMode
}

var (
	ErrorTransactionClosed = errors.New("transaction is closed")
	ErrorTransactionFailure = errors.New("transaction has failed")

	// ErrorTransactionCommitRollback occurs when an error has occurred in a transaction and
	// Commit() is called. PostgreSQL accepts COMMIT on aborted transactions, but
	// it is treated as ROLLBACK.
	ErrorTransactionCommitRollback = errors.New("commit unexpectedly resulted in rollback")
)

