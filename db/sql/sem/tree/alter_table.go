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

package tree

import (
    "github.com/Ready-Stock/noah/db/sql/coltypes"
)

// AlterTable represents an ALTER TABLE statement.
type AlterTable struct {
	IfExists bool
	Table    NormalizableTableName
	Cmds     AlterTableCmds
}

// Format implements the NodeFormatter interface.
func (node *AlterTable) Format(ctx *FmtCtx) {
	ctx.WriteString("ALTER TABLE ")
	if node.IfExists {
		ctx.WriteString("IF EXISTS ")
	}
	ctx.FormatNode(&node.Table)
	ctx.FormatNode(&node.Cmds)
}

// AlterTableCmds represents a list of table alterations.
type AlterTableCmds []AlterTableCmd

// Format implements the NodeFormatter interface.
func (node *AlterTableCmds) Format(ctx *FmtCtx) {
	for i, n := range *node {
		if i > 0 {
			ctx.WriteString(",")
		}
		ctx.FormatNode(n)
	}
}

// AlterTableCmd represents a table modification operation.
type AlterTableCmd interface {
	NodeFormatter
	// Placeholder function to ensure that only desired types
	// (AlterTable*) conform to the AlterTableCmd interface.
	alterTableCmd()
}

func (*AlterTableAddColumn) alterTableCmd()          {}
func (*AlterTableAddConstraint) alterTableCmd()      {}
func (*AlterTableAlterColumnType) alterTableCmd()    {}
func (*AlterTableDropColumn) alterTableCmd()         {}
func (*AlterTableDropConstraint) alterTableCmd()     {}
func (*AlterTableDropNotNull) alterTableCmd()        {}
func (*AlterTableDropStored) alterTableCmd()         {}
func (*AlterTableSetAudit) alterTableCmd()           {}
func (*AlterTableSetDefault) alterTableCmd()         {}
func (*AlterTableValidateConstraint) alterTableCmd() {}
func (*AlterTablePartitionBy) alterTableCmd()        {}
func (*AlterTableInjectStats) alterTableCmd()        {}

var _ AlterTableCmd = &AlterTableAddColumn{}
var _ AlterTableCmd = &AlterTableAddConstraint{}
var _ AlterTableCmd = &AlterTableAlterColumnType{}
var _ AlterTableCmd = &AlterTableDropColumn{}
var _ AlterTableCmd = &AlterTableDropConstraint{}
var _ AlterTableCmd = &AlterTableDropNotNull{}
var _ AlterTableCmd = &AlterTableDropStored{}
var _ AlterTableCmd = &AlterTableSetAudit{}
var _ AlterTableCmd = &AlterTableSetDefault{}
var _ AlterTableCmd = &AlterTableValidateConstraint{}
var _ AlterTableCmd = &AlterTablePartitionBy{}
var _ AlterTableCmd = &AlterTableInjectStats{}

// ColumnMutationCmd is the subset of AlterTableCmds that modify an
// existing column.
type ColumnMutationCmd interface {
	AlterTableCmd
	GetColumn() Name
}

// AlterTableAddColumn represents an ADD COLUMN command.
type AlterTableAddColumn struct {
	IfNotExists bool
	ColumnDef   *ColumnTableDef
}

// Format implements the NodeFormatter interface.
func (node *AlterTableAddColumn) Format(ctx *FmtCtx) {
	ctx.WriteString(" ADD COLUMN ")
	if node.IfNotExists {
		ctx.WriteString("IF NOT EXISTS ")
	}
	ctx.FormatNode(node.ColumnDef)
}

// ValidationBehavior specifies whether or not a constraint is validated.
type ValidationBehavior int

const (
	// ValidationDefault is the default validation behavior (immediate).
	ValidationDefault ValidationBehavior = iota
	// ValidationSkip skips validation of any existing data.
	ValidationSkip
)

// AlterTableAddConstraint represents an ADD CONSTRAINT command.
type AlterTableAddConstraint struct {
	ConstraintDef      ConstraintTableDef
	ValidationBehavior ValidationBehavior
}

// Format implements the NodeFormatter interface.
func (node *AlterTableAddConstraint) Format(ctx *FmtCtx) {
	ctx.WriteString(" ADD ")
	ctx.FormatNode(node.ConstraintDef)
	if node.ValidationBehavior == ValidationSkip {
		ctx.WriteString(" NOT VALID")
	}
}

// AlterTableAlterColumnType represents an ALTER TABLE ALTER COLUMN TYPE command.
type AlterTableAlterColumnType struct {
	Collation string
	Column    Name
	ToType    coltypes.T
	Using     Expr
}

// Format implements the NodeFormatter interface.
func (node *AlterTableAlterColumnType) Format(ctx *FmtCtx) {
	ctx.WriteString(" ALTER COLUMN ")
	ctx.FormatNode(&node.Column)
	ctx.WriteString(" SET DATA TYPE ")
	node.ToType.Format(ctx.Buffer, ctx.flags.EncodeFlags())
	if len(node.Collation) > 0 {
		ctx.WriteString(" COLLATE ")
		ctx.WriteString(node.Collation)
	}
	if node.Using != nil {
		ctx.WriteString(" USING ")
		ctx.FormatNode(node.Using)
	}
}

// GetColumn implements the ColumnMutationCmd interface.
func (node *AlterTableAlterColumnType) GetColumn() Name {
	return node.Column
}

// AlterTableDropColumn represents a DROP COLUMN command.
type AlterTableDropColumn struct {
	IfExists     bool
	Column       Name
	DropBehavior DropBehavior
}

// Format implements the NodeFormatter interface.
func (node *AlterTableDropColumn) Format(ctx *FmtCtx) {
	ctx.WriteString(" DROP COLUMN ")
	if node.IfExists {
		ctx.WriteString("IF EXISTS ")
	}
	ctx.FormatNode(&node.Column)
	if node.DropBehavior != DropDefault {
		ctx.Printf(" %s", node.DropBehavior)
	}
}

// AlterTableDropConstraint represents a DROP CONSTRAINT command.
type AlterTableDropConstraint struct {
	IfExists     bool
	Constraint   Name
	DropBehavior DropBehavior
}

// Format implements the NodeFormatter interface.
func (node *AlterTableDropConstraint) Format(ctx *FmtCtx) {
	ctx.WriteString(" DROP CONSTRAINT ")
	if node.IfExists {
		ctx.WriteString("IF EXISTS ")
	}
	ctx.FormatNode(&node.Constraint)
	if node.DropBehavior != DropDefault {
		ctx.Printf(" %s", node.DropBehavior)
	}
}

// AlterTableValidateConstraint represents a VALIDATE CONSTRAINT command.
type AlterTableValidateConstraint struct {
	Constraint Name
}

// Format implements the NodeFormatter interface.
func (node *AlterTableValidateConstraint) Format(ctx *FmtCtx) {
	ctx.WriteString(" VALIDATE CONSTRAINT ")
	ctx.FormatNode(&node.Constraint)
}

// AlterTableSetDefault represents an ALTER COLUMN SET DEFAULT
// or DROP DEFAULT command.
type AlterTableSetDefault struct {
	Column  Name
	Default Expr
}

// GetColumn implements the ColumnMutationCmd interface.
func (node *AlterTableSetDefault) GetColumn() Name {
	return node.Column
}

// Format implements the NodeFormatter interface.
func (node *AlterTableSetDefault) Format(ctx *FmtCtx) {
	ctx.WriteString(" ALTER COLUMN ")
	ctx.FormatNode(&node.Column)
	if node.Default == nil {
		ctx.WriteString(" DROP DEFAULT")
	} else {
		ctx.WriteString(" SET DEFAULT ")
		ctx.FormatNode(node.Default)
	}
}

// AlterTableDropNotNull represents an ALTER COLUMN DROP NOT NULL
// command.
type AlterTableDropNotNull struct {
	Column Name
}

// GetColumn implements the ColumnMutationCmd interface.
func (node *AlterTableDropNotNull) GetColumn() Name {
	return node.Column
}

// Format implements the NodeFormatter interface.
func (node *AlterTableDropNotNull) Format(ctx *FmtCtx) {
	ctx.WriteString(" ALTER COLUMN ")
	ctx.FormatNode(&node.Column)
	ctx.WriteString(" DROP NOT NULL")
}

// AlterTableDropStored represents an ALTER COLUMN DROP STORED command
// to remove the computed-ness from a column.
type AlterTableDropStored struct {
	Column Name
}

// GetColumn implemnets the ColumnMutationCmd interface.
func (node *AlterTableDropStored) GetColumn() Name {
	return node.Column
}

// Format implements the NodeFormatter interface.
func (node *AlterTableDropStored) Format(ctx *FmtCtx) {
	ctx.WriteString(" ALTER COLUMN ")
	ctx.FormatNode(&node.Column)
	ctx.WriteString(" DROP STORED")
}

// AlterTablePartitionBy represents an ALTER TABLE PARTITION BY
// command.
type AlterTablePartitionBy struct {
	*PartitionBy
}

// Format implements the NodeFormatter interface.
func (node *AlterTablePartitionBy) Format(ctx *FmtCtx) {
	ctx.FormatNode(node.PartitionBy)
}

// AuditMode represents a table audit mode
type AuditMode int

const (
	// AuditModeDisable is the default mode - no audit.
	AuditModeDisable AuditMode = iota
	// AuditModeReadWrite enables audit on read or write statements.
	AuditModeReadWrite
)

var auditModeName = [...]string{
	AuditModeDisable:   "OFF",
	AuditModeReadWrite: "READ WRITE",
}

func (m AuditMode) String() string {
	return auditModeName[m]
}

// AlterTableSetAudit represents an ALTER TABLE AUDIT SET statement.
type AlterTableSetAudit struct {
	Mode AuditMode
}

// Format implements the NodeFormatter interface.
func (node *AlterTableSetAudit) Format(ctx *FmtCtx) {
	ctx.WriteString(" EXPERIMENTAL_AUDIT SET ")
	ctx.WriteString(node.Mode.String())
}

// AlterTableInjectStats represents an ALTER TABLE INJECT STATISTICS statement.
type AlterTableInjectStats struct {
	Stats Expr
}

// Format implements the NodeFormatter interface.
func (node *AlterTableInjectStats) Format(ctx *FmtCtx) {
	ctx.WriteString(" INJECT STATISTICS ")
	ctx.FormatNode(node.Stats)
}
