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
 *
 * Project: pq github.com/lib/pq
 * Copyright 2018  'pq' Contributors Portions Copyright (C) 2011 Blake Mizerany
 * License https://github.com/lib/pq/blob/master/LICENSE.md
 */

package tree

import (
	"fmt"
	"strings"

    "github.com/Ready-Stock/noah/db/sql/pgwire/pgerror"
)

// IsolationLevel holds the isolation level for a transaction.
type IsolationLevel int

// IsolationLevel values
const (
	UnspecifiedIsolation IsolationLevel = iota
	SerializableIsolation
)

var isolationLevelNames = [...]string{
	UnspecifiedIsolation:  "UNSPECIFIED",
	SerializableIsolation: "SERIALIZABLE",
}

// IsolationLevelMap is a map from string isolation level name to isolation
// level, in the lowercase format that set isolation_level supports.
var IsolationLevelMap = map[string]IsolationLevel{
	"serializable": SerializableIsolation,
	"snapshot":     SerializableIsolation,
}

func (i IsolationLevel) String() string {
	if i < 0 || i > IsolationLevel(len(isolationLevelNames)-1) {
		return fmt.Sprintf("IsolationLevel(%d)", i)
	}
	return isolationLevelNames[i]
}

// UserPriority holds the user priority for a transaction.
type UserPriority int

// UserPriority values
const (
	UnspecifiedUserPriority UserPriority = iota
	Low
	Normal
	High
)

var userPriorityNames = [...]string{
	UnspecifiedUserPriority: "UNSPECIFIED",
	Low:    "LOW",
	Normal: "NORMAL",
	High:   "HIGH",
}

func (up UserPriority) String() string {
	if up < 0 || up > UserPriority(len(userPriorityNames)-1) {
		return fmt.Sprintf("UserPriority(%d)", up)
	}
	return userPriorityNames[up]
}

// ReadWriteMode holds the read write mode for a transaction.
type ReadWriteMode int

// ReadWriteMode values
const (
	UnspecifiedReadWriteMode ReadWriteMode = iota
	ReadOnly
	ReadWrite
)

var readWriteModeNames = [...]string{
	UnspecifiedReadWriteMode: "UNSPECIFIED",
	ReadOnly:                 "ONLY",
	ReadWrite:                "WRITE",
}

func (ro ReadWriteMode) String() string {
	if ro < 0 || ro > ReadWriteMode(len(readWriteModeNames)-1) {
		return fmt.Sprintf("ReadWriteMode(%d)", ro)
	}
	return readWriteModeNames[ro]
}

// TransactionModes holds the transaction modes for a transaction.
type TransactionModes struct {
	Isolation     IsolationLevel
	UserPriority  UserPriority
	ReadWriteMode ReadWriteMode
}

// Format implements the NodeFormatter interface.
func (node *TransactionModes) Format(ctx *FmtCtx) {
	var sep string
	if node.Isolation != UnspecifiedIsolation {
		ctx.Printf(" ISOLATION LEVEL %s", node.Isolation)
		sep = ","
	}
	if node.UserPriority != UnspecifiedUserPriority {
		ctx.Printf("%s PRIORITY %s", sep, node.UserPriority)
		sep = ","
	}
	if node.ReadWriteMode != UnspecifiedReadWriteMode {
		ctx.Printf("%s READ %s", sep, node.ReadWriteMode)
	}
}

var (
	errIsolationLevelSpecifiedMultipleTimes = pgerror.NewError(pgerror.CodeSyntaxError, "isolation level specified multiple times")
	errUserPrioritySpecifiedMultipleTimes   = pgerror.NewError(pgerror.CodeSyntaxError, "user priority specified multiple times")
	errReadModeSpecifiedMultipleTimes       = pgerror.NewError(pgerror.CodeSyntaxError, "read mode specified multiple times")
)

// Merge groups two sets of transaction modes together.
// Used in the parser.
func (node *TransactionModes) Merge(other TransactionModes) error {
	if other.Isolation != UnspecifiedIsolation {
		if node.Isolation != UnspecifiedIsolation {
			return errIsolationLevelSpecifiedMultipleTimes
		}
		node.Isolation = other.Isolation
	}
	if other.UserPriority != UnspecifiedUserPriority {
		if node.UserPriority != UnspecifiedUserPriority {
			return errUserPrioritySpecifiedMultipleTimes
		}
		node.UserPriority = other.UserPriority
	}
	if other.ReadWriteMode != UnspecifiedReadWriteMode {
		if node.ReadWriteMode != UnspecifiedReadWriteMode {
			return errReadModeSpecifiedMultipleTimes
		}
		node.ReadWriteMode = other.ReadWriteMode
	}
	return nil
}

// BeginTransaction represents a BEGIN statement
type BeginTransaction struct {
	Modes TransactionModes
}

// Format implements the NodeFormatter interface.
func (node *BeginTransaction) Format(ctx *FmtCtx) {
	ctx.WriteString("BEGIN TRANSACTION")
	node.Modes.Format(ctx)
}

// CommitTransaction represents a COMMIT statement.
type CommitTransaction struct{}

// Format implements the NodeFormatter interface.
func (node *CommitTransaction) Format(ctx *FmtCtx) {
	ctx.WriteString("COMMIT TRANSACTION")
}

// RollbackTransaction represents a ROLLBACK statement.
type RollbackTransaction struct{}

// Format implements the NodeFormatter interface.
func (node *RollbackTransaction) Format(ctx *FmtCtx) {
	ctx.WriteString("ROLLBACK TRANSACTION")
}

// RestartSavepointName is the only savepoint name that we accept, modulo
// capitalization.
const RestartSavepointName string = "COCKROACH_RESTART"

// ValidateRestartCheckpoint checks that a checkpoint name is our magic restart
// value.
// We accept everything with the desired prefix because at least the C++ libpqxx
// appends sequence numbers to the savepoint name specified by the user.
func ValidateRestartCheckpoint(savepoint string) error {
	if !strings.HasPrefix(strings.ToUpper(savepoint), RestartSavepointName) {
		return pgerror.NewErrorf(pgerror.CodeFeatureNotSupportedError, "SAVEPOINT not supported except for %s", RestartSavepointName)
	}
	return nil
}

// Savepoint represents a SAVEPOINT <name> statement.
type Savepoint struct {
	Name string
}

// Format implements the NodeFormatter interface.
func (node *Savepoint) Format(ctx *FmtCtx) {
	ctx.WriteString("SAVEPOINT ")
	ctx.WriteString(node.Name)
}

// ReleaseSavepoint represents a RELEASE SAVEPOINT <name> statement.
type ReleaseSavepoint struct {
	Savepoint string
}

// Format implements the NodeFormatter interface.
func (node *ReleaseSavepoint) Format(ctx *FmtCtx) {
	ctx.WriteString("RELEASE SAVEPOINT ")
	ctx.WriteString(node.Savepoint)
}

// RollbackToSavepoint represents a ROLLBACK TO SAVEPOINT <name> statement.
type RollbackToSavepoint struct {
	Savepoint string
}

// Format implements the NodeFormatter interface.
func (node *RollbackToSavepoint) Format(ctx *FmtCtx) {
	ctx.WriteString("ROLLBACK TRANSACTION TO SAVEPOINT ")
	ctx.WriteString(node.Savepoint)
}
