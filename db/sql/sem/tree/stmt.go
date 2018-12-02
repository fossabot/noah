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

// This code was derived from https://github.com/youtube/vitess.

package tree

import (
	"fmt"
		)

// Instructions for creating new types: If a type needs to satisfy an
// interface, declare that function along with that interface. This
// will help users identify the list of types to which they can assert
// those interfaces. If the member of a type has a string with a
// predefined list of values, declare those values as const following
// the type. For interfaces that define dummy functions to
// consolidate a set of types, define the function as typeName().
// This will help avoid name collisions.

// StatementType is the enumerated type for Statement return styles on
// the wire.
type StatementType int

//go:generate stringer -type=StatementType
const (
	// Ack indicates that the statement does not have a meaningful
	// return. Examples include SET, BEGIN, COMMIT.
	Ack StatementType = iota
	// DDL indicates that the statement mutates the database schema.
	//
	// Note: this is the type indicated back to the client; it is not a
	// sufficient test for schema mutation for planning purposes. There
	// are schema-modifying statements (e.g. CREATE TABLE AS) which
	// report RowsAffected to the client, not DDL.
	// Use CanModifySchema() below instead.
	DDL
	// RowsAffected indicates that the statement returns the count of
	// affected rows.
	RowsAffected
	// Rows indicates that the statement returns the affected rows after
	// the statement was applied.
	Rows
	// CopyIn indicates a COPY FROM statement.
	CopyIn
	// Unknown indicates that the statement does not have a known
	// return style at the time of parsing. This is not first in the
	// enumeration because it is more convenient to have Ack as a zero
	// value, and because the use of Unknown should be an explicit choice.
	// The primary example of this statement type is EXECUTE, where the
	// statement type depends on the statement type of the prepared statement
	// being executed.
	Unknown
)

// func GetStatementType(stmt pg_query.ParsetreeList) (*StatementType, error) {
// 	if len(stmt.Statements) == 0 {
// 		return nil, pgerror.NewErrorf(pgerror.CodeSyntaxError, "no statements provided")
// 	} else {
// 		switch stmt.Statements[0].(type) {
//
// 		}
// 	}
// }

// Statement represents a statement.
type Statement interface {
	fmt.Stringer
	NodeFormatter
	StatementType() StatementType
	// StatementTag is a short string identifying the type of statement
	// (usually a single verb). This is different than the Stringer output,
	// which is the actual statement (including args).
	// TODO(dt): Currently tags are always pg-compatible in the future it
	// might make sense to pass a tag format specifier.
	StatementTag() string
}

// canModifySchema is to be implemented by statements that can modify
// the database schema but may have StatementType() != DDL.
// See CanModifySchema() below.
type canModifySchema interface {
	modifiesSchema() bool
}

// CanModifySchema returns true if the statement can modify
// the database schema.
func CanModifySchema(stmt Statement) bool {
	if stmt.StatementType() == DDL {
		return true
	}
	scm, ok := stmt.(canModifySchema)
	return ok && scm.modifiesSchema()
}

// CanWriteData returns true if the statement can modify data.
func CanWriteData(stmt Statement) bool {
	switch stmt.(type) {
	// Normal write operations.
	case *Insert, *Delete, *Update, *Truncate:
		return true
	// Import operations.
	case *CopyFrom, *Import, *Restore:
		return true
	// CockroachDB extensions.
	case *Split, *Relocate, *Scatter:
		return true
	}
	return false
}

// HiddenFromStats is a pseudo-interface to be implemented
// by statements that should not show up in per-app statistics.
type HiddenFromStats interface {
	hiddenFromStats()
}

// HiddenFromShowQueries is a pseudo-interface to be implemented
// by statements that should not show up in SHOW QUERIES (and are hence
// not cancellable using CANCEL QUERIES either). Usually implemented by
// statements that spawn jobs.
type HiddenFromShowQueries interface {
	hiddenFromShowQueries()
}

// IndependentFromParallelizedPriors is a pseudo-interface to be implemented
// by statements which do not force parallel statement execution synchronization
// when they run.
// NB: Only statements that don't send any requests using the current
// transaction can implement this. Otherwise, the statement will fail if any of
// the parallel statements has encoutered a KV error (which toasts the txn).
type IndependentFromParallelizedPriors interface {
	independentFromParallelizedPriors()
}

// StatementList is a list of statements.
type StatementList []Statement

// Format implements the NodeFormatter interface.
func (l *StatementList) Format(ctx *FmtCtx) {
	for i, s := range *l {
		if i > 0 {
			ctx.WriteString("; ")
		}
		ctx.FormatNode(s)
	}
}

// ObserverStatement is a marker interface for statements which are allowed to
// run regardless of the current transaction state: statements other than
// rollback are generally rejected if the session is in a failed transaction
// state, but it's convenient to allow some statements (e.g. "show syntax; set
// tracing").
// Such statements are not expected to modify the database, the transaction or
// session state (other than special cases such as enabling/disabling tracing).
//
// These statements short-circuit the regular execution - they don't get planned
// (there are no corresponding planNodes). The connExecutor recognizes them and
// handles them.
type ObserverStatement interface {
	observerStatement()
}

// StatementType implements the Statement interface.
func (*AlterIndex) StatementType() StatementType { return DDL }

// StatementTag returns a short string identifying the type of statement.
func (*AlterIndex) StatementTag() string { return "ALTER INDEX" }

func (*AlterIndex) hiddenFromShowQueries() {}

// StatementType implements the Statement interface.
func (*AlterTable) StatementType() StatementType { return DDL }

// StatementTag returns a short string identifying the type of statement.
func (*AlterTable) StatementTag() string { return "ALTER TABLE" }

func (*AlterTable) hiddenFromShowQueries() {}

// StatementType implements the Statement interface.
func (*AlterSequence) StatementType() StatementType { return DDL }

// StatementTag returns a short string identifying the type of statement.
func (*AlterSequence) StatementTag() string { return "ALTER SEQUENCE" }

// StatementType implements the Statement interface.
func (*AlterUserSetPassword) StatementType() StatementType { return RowsAffected }

// StatementTag returns a short string identifying the type of statement.
func (*AlterUserSetPassword) StatementTag() string { return "ALTER USER" }

func (*AlterUserSetPassword) hiddenFromShowQueries() {}

// StatementType implements the Statement interface.
func (*Backup) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*Backup) StatementTag() string { return "BACKUP" }

func (*Backup) hiddenFromShowQueries() {}

// StatementType implements the Statement interface.
func (*BeginTransaction) StatementType() StatementType { return Ack }

// StatementTag returns a short string identifying the type of statement.
func (*BeginTransaction) StatementTag() string { return "BEGIN" }

func (*BeginTransaction) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*ControlJobs) StatementType() StatementType { return RowsAffected }

// StatementTag returns a short string identifying the type of statement.
func (n *ControlJobs) StatementTag() string {
	return fmt.Sprintf("%s JOBS", JobCommandToStatement[n.Command])
}

func (*ControlJobs) independentFromParallelizedPriors() {}

// StatementType implements the Statement interface.
func (*CancelQueries) StatementType() StatementType { return RowsAffected }

// StatementTag returns a short string identifying the type of statement.
func (*CancelQueries) StatementTag() string { return "CANCEL QUERIES" }

func (*CancelQueries) independentFromParallelizedPriors() {}

// StatementType implements the Statement interface.
func (*CancelSessions) StatementType() StatementType { return RowsAffected }

// StatementTag returns a short string identifying the type of statement.
func (*CancelSessions) StatementTag() string { return "CANCEL SESSIONS" }

func (*CancelSessions) independentFromParallelizedPriors() {}

// StatementType implements the Statement interface.
func (*CommitTransaction) StatementType() StatementType { return Ack }

// StatementTag returns a short string identifying the type of statement.
func (*CommitTransaction) StatementTag() string { return "COMMIT" }

func (*CommitTransaction) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*CopyFrom) StatementType() StatementType { return CopyIn }

// StatementTag returns a short string identifying the type of statement.
func (*CopyFrom) StatementTag() string { return "COPY" }

// StatementType implements the Statement interface.
func (*CreateChangefeed) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*CreateChangefeed) StatementTag() string { return "CREATE CHANGEFEED" }

// StatementType implements the Statement interface.
func (*CreateDatabase) StatementType() StatementType { return DDL }

// StatementTag returns a short string identifying the type of statement.
func (*CreateDatabase) StatementTag() string { return "CREATE DATABASE" }

// StatementType implements the Statement interface.
func (*CreateIndex) StatementType() StatementType { return DDL }

// StatementTag returns a short string identifying the type of statement.
func (*CreateIndex) StatementTag() string { return "CREATE INDEX" }

// StatementType implements the Statement interface.
func (n *CreateTable) StatementType() StatementType {
	if n.As() {
		return RowsAffected
	}
	return DDL
}

// StatementTag returns a short string identifying the type of statement.
func (n *CreateTable) StatementTag() string {
	if n.As() {
		return "SELECT"
	}
	return "CREATE TABLE"
}

// modifiesSchema implements the canModifySchema interface.
func (*CreateTable) modifiesSchema() bool { return true }

// StatementType implements the Statement interface.
func (*CreateUser) StatementType() StatementType { return RowsAffected }

// StatementTag returns a short string identifying the type of statement.
func (*CreateUser) StatementTag() string { return "CREATE USER" }

func (*CreateUser) hiddenFromShowQueries() {}

// StatementType implements the Statement interface.
func (*CreateRole) StatementType() StatementType { return RowsAffected }

// StatementTag returns a short string identifying the type of statement.
func (*CreateRole) StatementTag() string { return "CREATE ROLE" }

func (*CreateRole) hiddenFromShowQueries() {}

// StatementType implements the Statement interface.
func (*CreateView) StatementType() StatementType { return DDL }

// StatementTag returns a short string identifying the type of statement.
func (*CreateView) StatementTag() string { return "CREATE VIEW" }

// StatementType implements the Statement interface.
func (*CreateSequence) StatementType() StatementType { return DDL }

// StatementTag returns a short string identifying the type of statement.
func (*CreateSequence) StatementTag() string { return "CREATE SEQUENCE" }

// StatementType implements the Statement interface.
func (*CreateStats) StatementType() StatementType { return DDL }

// StatementTag returns a short string identifying the type of statement.
func (*CreateStats) StatementTag() string { return "CREATE STATISTICS" }

// StatementType implements the Statement interface.
func (*Deallocate) StatementType() StatementType { return Ack }

// StatementTag returns a short string identifying the type of statement.
func (n *Deallocate) StatementTag() string {
	// Postgres distinguishes the command tags for these two cases of Deallocate statements.
	if n.Name == "" {
		return "DEALLOCATE ALL"
	}
	return "DEALLOCATE"
}

func (*Deallocate) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*Discard) StatementType() StatementType { return Ack }

// StatementTag returns a short string identifying the type of statement.
func (*Discard) StatementTag() string { return "DISCARD" }

// StatementType implements the Statement interface.
func (n *Delete) StatementType() StatementType { return n.Returning.statementType() }

// StatementTag returns a short string identifying the type of statement.
func (*Delete) StatementTag() string { return "DELETE" }

// StatementType implements the Statement interface.
func (*DropDatabase) StatementType() StatementType { return DDL }

// StatementTag returns a short string identifying the type of statement.
func (*DropDatabase) StatementTag() string { return "DROP DATABASE" }

// StatementType implements the Statement interface.
func (*DropIndex) StatementType() StatementType { return DDL }

// StatementTag returns a short string identifying the type of statement.
func (*DropIndex) StatementTag() string { return "DROP INDEX" }

// StatementType implements the Statement interface.
func (*DropTable) StatementType() StatementType { return DDL }

// StatementTag returns a short string identifying the type of statement.
func (*DropTable) StatementTag() string { return "DROP TABLE" }

// StatementType implements the Statement interface.
func (*DropView) StatementType() StatementType { return DDL }

// StatementTag returns a short string identifying the type of statement.
func (*DropView) StatementTag() string { return "DROP VIEW" }

// StatementType implements the Statement interface.
func (*DropSequence) StatementType() StatementType { return DDL }

// StatementTag returns a short string identifying the type of statement.
func (*DropSequence) StatementTag() string { return "DROP SEQUENCE" }

// StatementType implements the Statement interface.
func (*DropUser) StatementType() StatementType { return RowsAffected }

// StatementTag returns a short string identifying the type of statement.
func (*DropUser) StatementTag() string { return "DROP USER" }

// StatementType implements the Statement interface.
func (*DropRole) StatementType() StatementType { return RowsAffected }

// StatementTag returns a short string identifying the type of statement.
func (*DropRole) StatementTag() string { return "DROP ROLE" }

// StatementType implements the Statement interface.
func (*Execute) StatementType() StatementType { return Unknown }

// StatementTag returns a short string identifying the type of statement.
func (*Execute) StatementTag() string { return "EXECUTE" }

func (*Execute) independentFromParallelizedPriors() {}

// StatementType implements the Statement interface.
// func (*Explain) StatementType() StatementType { return Rows }
//
// // StatementTag returns a short string identifying the type of statement.
// func (*Explain) StatementTag() string { return "EXPLAIN" }
//
// func (*Explain) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*Export) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*Export) StatementTag() string { return "EXPORT" }

// StatementType implements the Statement interface.
func (*Grant) StatementType() StatementType { return DDL }

// StatementTag returns a short string identifying the type of statement.
func (*Grant) StatementTag() string { return "GRANT" }

func (*Grant) hiddenFromStats() {}

// StatementType implements the Statement interface.
// func (*GrantRole) StatementType() StatementType { return DDL }
//
// // StatementTag returns a short string identifying the type of statement.
// func (*GrantRole) StatementTag() string { return "GRANT" }
//
// func (*GrantRole) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (n *Insert) StatementType() StatementType { return n.Returning.statementType() }

// StatementTag returns a short string identifying the type of statement.
func (*Insert) StatementTag() string { return "INSERT" }

// StatementType implements the Statement interface.
func (n *Import) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*Import) StatementTag() string { return "IMPORT" }

// StatementType implements the Statement interface.
func (*ParenSelect) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ParenSelect) StatementTag() string { return "SELECT" }

// StatementType implements the Statement interface.
func (*Prepare) StatementType() StatementType { return Ack }

// StatementTag returns a short string identifying the type of statement.
func (*Prepare) StatementTag() string { return "PREPARE" }

func (*Prepare) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*ReleaseSavepoint) StatementType() StatementType { return Ack }

// StatementTag returns a short string identifying the type of statement.
func (*ReleaseSavepoint) StatementTag() string { return "RELEASE" }

func (*ReleaseSavepoint) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*RenameColumn) StatementType() StatementType { return DDL }

// StatementTag returns a short string identifying the type of statement.
func (*RenameColumn) StatementTag() string { return "RENAME COLUMN" }

// StatementType implements the Statement interface.
func (*RenameDatabase) StatementType() StatementType { return DDL }

// StatementTag returns a short string identifying the type of statement.
func (*RenameDatabase) StatementTag() string { return "RENAME DATABASE" }

// StatementType implements the Statement interface.
func (*RenameIndex) StatementType() StatementType { return DDL }

// StatementTag returns a short string identifying the type of statement.
func (*RenameIndex) StatementTag() string { return "RENAME INDEX" }

// StatementType implements the Statement interface.
func (*RenameTable) StatementType() StatementType { return DDL }

// StatementTag returns a short string identifying the type of statement.
func (n *RenameTable) StatementTag() string {
	if n.IsView {
		return "RENAME VIEW"
	} else if n.IsSequence {
		return "RENAME SEQUENCE"
	}
	return "RENAME TABLE"
}

// StatementType implements the Statement interface.
func (*Relocate) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (n *Relocate) StatementTag() string {
	if n.RelocateLease {
		return "EXPERIMENTAL_RELOCATE LEASE"
	}
	return "EXPERIMENTAL_RELOCATE"
}

// StatementType implements the Statement interface.
func (*Restore) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*Restore) StatementTag() string { return "RESTORE" }

func (*Restore) hiddenFromShowQueries() {}

// // StatementType implements the Statement interface.
// func (*Revoke) StatementType() StatementType { return DDL }
//
// // StatementTag returns a short string identifying the type of statement.
// func (*Revoke) StatementTag() string { return "REVOKE" }
//
// func (*Revoke) hiddenFromStats() {}
//
// // StatementType implements the Statement interface.
// func (*RevokeRole) StatementType() StatementType { return DDL }
//
// // StatementTag returns a short string identifying the type of statement.
// func (*RevokeRole) StatementTag() string { return "REVOKE" }
//
// func (*RevokeRole) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*RollbackToSavepoint) StatementType() StatementType { return Ack }

// StatementTag returns a short string identifying the type of statement.
func (*RollbackToSavepoint) StatementTag() string { return "ROLLBACK" }

func (*RollbackToSavepoint) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*RollbackTransaction) StatementType() StatementType { return Ack }

// StatementTag returns a short string identifying the type of statement.
func (*RollbackTransaction) StatementTag() string { return "ROLLBACK" }

func (*RollbackTransaction) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*Savepoint) StatementType() StatementType { return Ack }

// StatementTag returns a short string identifying the type of statement.
func (*Savepoint) StatementTag() string { return "SAVEPOINT" }

// StatementType implements the Statement interface.
func (*Scatter) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*Scatter) StatementTag() string { return "SCATTER" }

// StatementType implements the Statement interface.
func (*Scrub) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (n *Scrub) StatementTag() string { return "SCRUB" }

// StatementType implements the Statement interface.
func (*Select) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*Select) StatementTag() string { return "SELECT" }

// StatementType implements the Statement interface.
func (*SelectClause) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*SelectClause) StatementTag() string { return "SELECT" }

// StatementType implements the Statement interface.
func (*SetVar) StatementType() StatementType { return Ack }

// StatementTag returns a short string identifying the type of statement.
func (*SetVar) StatementTag() string { return "SET" }

func (*SetVar) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*SetClusterSetting) StatementType() StatementType { return Ack }

// StatementTag returns a short string identifying the type of statement.
func (*SetClusterSetting) StatementTag() string { return "SET CLUSTER SETTING" }

// StatementType implements the Statement interface.
func (*SetTransaction) StatementType() StatementType { return Ack }

// StatementTag returns a short string identifying the type of statement.
func (*SetTransaction) StatementTag() string { return "SET TRANSACTION" }

func (*SetTransaction) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*SetTracing) StatementType() StatementType { return Ack }

// StatementTag returns a short string identifying the type of statement.
func (*SetTracing) StatementTag() string { return "SET TRACING" }

func (*SetTracing) hiddenFromStats() {}

// observerStatement implements the ObserverStatement interface.
func (*SetTracing) observerStatement() {}

// StatementType implements the Statement interface.
func (*SetZoneConfig) StatementType() StatementType { return RowsAffected }

// StatementTag returns a short string identifying the type of statement.
func (*SetZoneConfig) StatementTag() string { return "CONFIGURE ZONE" }

// StatementType implements the Statement interface.
func (*SetSessionCharacteristics) StatementType() StatementType { return Ack }

// StatementTag returns a short string identifying the type of statement.
func (*SetSessionCharacteristics) StatementTag() string { return "SET" }

func (*SetSessionCharacteristics) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*ShowVar) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowVar) StatementTag() string { return "SHOW" }

func (*ShowVar) hiddenFromStats()                   {}
func (*ShowVar) independentFromParallelizedPriors() {}

// StatementType implements the Statement interface.
func (*ShowClusterSetting) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowClusterSetting) StatementTag() string { return "SHOW" }

func (*ShowClusterSetting) hiddenFromStats()                   {}
func (*ShowClusterSetting) independentFromParallelizedPriors() {}

// StatementType implements the Statement interface.
func (*ShowColumns) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowColumns) StatementTag() string { return "SHOW COLUMNS" }

func (*ShowColumns) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*ShowCreate) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowCreate) StatementTag() string { return "SHOW CREATE" }

func (*ShowCreate) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*ShowBackup) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowBackup) StatementTag() string { return "SHOW BACKUP" }

func (*ShowBackup) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*ShowDatabases) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowDatabases) StatementTag() string { return "SHOW DATABASES" }

func (*ShowDatabases) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*ShowTraceForSession) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowTraceForSession) StatementTag() string { return "SHOW TRACE FOR SESSION" }

func (*ShowTraceForSession) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*ShowGrants) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowGrants) StatementTag() string { return "SHOW GRANTS" }

func (*ShowGrants) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*ShowIndex) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowIndex) StatementTag() string { return "SHOW INDEX" }

func (*ShowIndex) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*ShowQueries) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowQueries) StatementTag() string { return "SHOW QUERIES" }

func (*ShowQueries) hiddenFromStats()                   {}
func (*ShowQueries) independentFromParallelizedPriors() {}

// StatementType implements the Statement interface.
func (*ShowJobs) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowJobs) StatementTag() string { return "SHOW JOBS" }

func (*ShowJobs) hiddenFromStats()                   {}
func (*ShowJobs) independentFromParallelizedPriors() {}

// StatementType implements the Statement interface.
func (*ShowRoleGrants) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowRoleGrants) StatementTag() string { return "SHOW GRANTS ON ROLE" }

func (*ShowRoleGrants) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*ShowSessions) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowSessions) StatementTag() string { return "SHOW SESSIONS" }

func (*ShowSessions) hiddenFromStats()                   {}
func (*ShowSessions) independentFromParallelizedPriors() {}

// StatementType implements the Statement interface.
func (*ShowTableStats) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowTableStats) StatementTag() string { return "SHOW STATISTICS" }

func (*ShowTableStats) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*ShowHistogram) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowHistogram) StatementTag() string { return "SHOW HISTOGRAM" }

func (*ShowHistogram) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*ShowSyntax) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowSyntax) StatementTag() string { return "SHOW SYNTAX" }

func (*ShowSyntax) observerStatement()                 {}
func (*ShowSyntax) hiddenFromStats()                   {}
func (*ShowSyntax) independentFromParallelizedPriors() {}

// StatementType implements the Statement interface.
func (*ShowTransactionStatus) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowTransactionStatus) StatementTag() string { return "SHOW TRANSACTION STATUS" }

func (*ShowTransactionStatus) observerStatement()                 {}
func (*ShowTransactionStatus) hiddenFromStats()                   {}
func (*ShowTransactionStatus) independentFromParallelizedPriors() {}

// StatementType implements the Statement interface.
func (*ShowUsers) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowUsers) StatementTag() string { return "SHOW USERS" }

func (*ShowUsers) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*ShowRoles) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowRoles) StatementTag() string { return "SHOW ROLES" }

func (*ShowRoles) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*ShowZoneConfig) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowZoneConfig) StatementTag() string { return "SHOW ZONE CONFIGURATION" }

// StatementType implements the Statement interface.
func (*ShowRanges) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowRanges) StatementTag() string { return "SHOW EXPERIMENTAL_RANGES" }

func (*ShowRanges) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*ShowFingerprints) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowFingerprints) StatementTag() string { return "SHOW EXPERIMENTAL_FINGERPRINTS" }

// StatementType implements the Statement interface.
func (*ShowConstraints) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowConstraints) StatementTag() string { return "SHOW CONSTRAINTS" }

func (*ShowConstraints) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*ShowTables) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowTables) StatementTag() string { return "SHOW TABLES" }

func (*ShowTables) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*ShowSchemas) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ShowSchemas) StatementTag() string { return "SHOW SCHEMAS" }

func (*ShowSchemas) hiddenFromStats() {}

// StatementType implements the Statement interface.
func (*Split) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*Split) StatementTag() string { return "SPLIT" }

// StatementType implements the Statement interface.
func (*Truncate) StatementType() StatementType { return Ack }

// StatementTag returns a short string identifying the type of statement.
func (*Truncate) StatementTag() string { return "TRUNCATE" }

// StatementType implements the Statement interface.
func (n *Update) StatementType() StatementType { return n.Returning.statementType() }

// StatementTag returns a short string identifying the type of statement.
func (*Update) StatementTag() string { return "UPDATE" }

// StatementType implements the Statement interface.
func (*UnionClause) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*UnionClause) StatementTag() string { return "UNION" }

// StatementType implements the Statement interface.
func (*ValuesClause) StatementType() StatementType { return Rows }

// StatementTag returns a short string identifying the type of statement.
func (*ValuesClause) StatementTag() string { return "VALUES" }

func (n *AlterIndex) String() string                { return AsString(n) }
func (n *AlterTable) String() string                { return AsString(n) }
func (n *AlterTableCmds) String() string            { return AsString(n) }
func (n *AlterTableAddColumn) String() string       { return AsString(n) }
func (n *AlterTableAddConstraint) String() string   { return AsString(n) }
func (n *AlterTableAlterColumnType) String() string { return AsString(n) }
func (n *AlterTableDropColumn) String() string      { return AsString(n) }
func (n *AlterTableDropConstraint) String() string  { return AsString(n) }
func (n *AlterTableDropNotNull) String() string     { return AsString(n) }
func (n *AlterTableDropStored) String() string      { return AsString(n) }
func (n *AlterTableSetDefault) String() string      { return AsString(n) }
func (n *AlterUserSetPassword) String() string      { return AsString(n) }
func (n *AlterSequence) String() string             { return AsString(n) }
func (n *Backup) String() string                    { return AsString(n) }
func (n *BeginTransaction) String() string          { return AsString(n) }
func (n *ControlJobs) String() string               { return AsString(n) }
func (n *CancelQueries) String() string             { return AsString(n) }
func (n *CancelSessions) String() string            { return AsString(n) }
func (n *CommitTransaction) String() string         { return AsString(n) }
func (n *CopyFrom) String() string                  { return AsString(n) }
func (n *CreateChangefeed) String() string          { return AsString(n) }
func (n *CreateDatabase) String() string            { return AsString(n) }
func (n *CreateIndex) String() string               { return AsString(n) }
func (n *CreateRole) String() string                { return AsString(n) }
func (n *CreateTable) String() string               { return AsString(n) }
func (n *CreateSequence) String() string            { return AsString(n) }
func (n *CreateStats) String() string               { return AsString(n) }
func (n *CreateUser) String() string                { return AsString(n) }
func (n *CreateView) String() string                { return AsString(n) }
func (n *Deallocate) String() string                { return AsString(n) }
func (n *Delete) String() string                    { return AsString(n) }
func (n *DropDatabase) String() string              { return AsString(n) }
func (n *DropIndex) String() string                 { return AsString(n) }
func (n *DropRole) String() string                  { return AsString(n) }
func (n *DropTable) String() string                 { return AsString(n) }
func (n *DropView) String() string                  { return AsString(n) }
func (n *DropSequence) String() string              { return AsString(n) }
func (n *DropUser) String() string                  { return AsString(n) }
func (n *Execute) String() string                   { return AsString(n) }
//func (n *Explain) String() string                   { return AsString(n) }
func (n *Export) String() string                    { return AsString(n) }
// func (n *Grant) String() string                     { return AsString(n) }
//func (n *GrantRole) String() string                 { return AsString(n) }
func (n *Insert) String() string                    { return AsString(n) }
func (n *Import) String() string                    { return AsString(n) }
func (n *ParenSelect) String() string               { return AsString(n) }
func (n *Prepare) String() string                   { return AsString(n) }
func (n *ReleaseSavepoint) String() string          { return AsString(n) }
func (n *Relocate) String() string                  { return AsString(n) }
func (n *RenameColumn) String() string              { return AsString(n) }
func (n *RenameDatabase) String() string            { return AsString(n) }
func (n *RenameIndex) String() string               { return AsString(n) }
func (n *RenameTable) String() string               { return AsString(n) }
func (n *Restore) String() string                   { return AsString(n) }
//func (n *Revoke) String() string                    { return AsString(n) }
//func (n *RevokeRole) String() string                { return AsString(n) }
func (n *RollbackToSavepoint) String() string       { return AsString(n) }
func (n *RollbackTransaction) String() string       { return AsString(n) }
func (n *Savepoint) String() string                 { return AsString(n) }
func (n *Scatter) String() string                   { return AsString(n) }
func (n *Scrub) String() string                     { return AsString(n) }
func (n *Select) String() string                    { return AsString(n) }
func (n *SelectClause) String() string              { return AsString(n) }
func (n *SetClusterSetting) String() string         { return AsString(n) }
func (n *SetZoneConfig) String() string             { return AsString(n) }
func (n *SetSessionCharacteristics) String() string { return AsString(n) }
func (n *SetTransaction) String() string            { return AsString(n) }
func (n *SetTracing) String() string                { return AsString(n) }
func (n *SetVar) String() string                    { return AsString(n) }
func (n *ShowBackup) String() string                { return AsString(n) }
func (n *ShowClusterSetting) String() string        { return AsString(n) }
func (n *ShowColumns) String() string               { return AsString(n) }
func (n *ShowConstraints) String() string           { return AsString(n) }
func (n *ShowCreate) String() string                { return AsString(n) }
func (n *ShowDatabases) String() string             { return AsString(n) }
func (n *ShowGrants) String() string                { return AsString(n) }
func (n *ShowHistogram) String() string             { return AsString(n) }
func (n *ShowIndex) String() string                 { return AsString(n) }
func (n *ShowJobs) String() string                  { return AsString(n) }
func (n *ShowQueries) String() string               { return AsString(n) }
func (n *ShowRanges) String() string                { return AsString(n) }
func (n *ShowRoleGrants) String() string            { return AsString(n) }
func (n *ShowRoles) String() string                 { return AsString(n) }
func (n *ShowSchemas) String() string               { return AsString(n) }
func (n *ShowSessions) String() string              { return AsString(n) }
func (n *ShowSyntax) String() string                { return AsString(n) }
func (n *ShowTableStats) String() string            { return AsString(n) }
func (n *ShowTables) String() string                { return AsString(n) }
func (n *ShowTraceForSession) String() string       { return AsString(n) }
func (n *ShowTransactionStatus) String() string     { return AsString(n) }
func (n *ShowUsers) String() string                 { return AsString(n) }
func (n *ShowVar) String() string                   { return AsString(n) }
func (n *ShowZoneConfig) String() string            { return AsString(n) }
func (n *ShowFingerprints) String() string          { return AsString(n) }
func (n *Split) String() string                     { return AsString(n) }
func (l *StatementList) String() string             { return AsString(l) }
func (n *Truncate) String() string                  { return AsString(n) }
func (n *UnionClause) String() string               { return AsString(n) }
func (n *Update) String() string                    { return AsString(n) }
func (n *ValuesClause) String() string              { return AsString(n) }
