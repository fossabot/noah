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

package tree

// Split represents an `ALTER TABLE/INDEX .. SPLIT AT ..` statement.
type Split struct {
	// Only one of Table and Index can be set.
	Table *NormalizableTableName
	Index *TableNameWithIndex
	// Each row contains values for the columns in the PK or index (or a prefix
	// of the columns).
	Rows *Select
}

// Format implements the NodeFormatter interface.
func (node *Split) Format(ctx *FmtCtx) {
	ctx.WriteString("ALTER ")
	if node.Index != nil {
		ctx.WriteString("INDEX ")
		ctx.FormatNode(node.Index)
	} else {
		ctx.WriteString("TABLE ")
		ctx.FormatNode(node.Table)
	}
	ctx.WriteString(" SPLIT AT ")
	ctx.FormatNode(node.Rows)
}

// Relocate represents an `ALTER TABLE/INDEX .. EXPERIMENTAL_RELOCATE ..`
// statement.
type Relocate struct {
	// Only one of Table and Index can be set.
	// TODO(a-robinson): It's not great that this can only work on ranges that
	// are part of a currently valid table or index.
	Table *NormalizableTableName
	Index *TableNameWithIndex
	// Each row contains an array with store ids and values for the columns in the
	// PK or index (or a prefix of the columns).
	// See docs/RFCS/sql_split_syntax.md.
	Rows          *Select
	RelocateLease bool
}

// Format implements the NodeFormatter interface.
func (node *Relocate) Format(ctx *FmtCtx) {
	ctx.WriteString("ALTER ")
	if node.Index != nil {
		ctx.WriteString("INDEX ")
		ctx.FormatNode(node.Index)
	} else {
		ctx.WriteString("TABLE ")
		ctx.FormatNode(node.Table)
	}
	ctx.WriteString(" EXPERIMENTAL_RELOCATE ")
	if node.RelocateLease {
		ctx.WriteString("LEASE ")
	}
	ctx.FormatNode(node.Rows)
}

// Scatter represents an `ALTER TABLE/INDEX .. SCATTER ..`
// statement.
type Scatter struct {
	// Only one of Table and Index can be set.
	Table *NormalizableTableName
	Index *TableNameWithIndex
	// Optional from and to values for the columns in the PK or index (or a prefix
	// of the columns).
	// See docs/RFCS/sql_split_syntax.md.
	From, To Exprs
}

// Format implements the NodeFormatter interface.
func (node *Scatter) Format(ctx *FmtCtx) {
	ctx.WriteString("ALTER ")
	if node.Index != nil {
		ctx.WriteString("INDEX ")
		ctx.FormatNode(node.Index)
	} else {
		ctx.WriteString("TABLE ")
		ctx.FormatNode(node.Table)
	}
	ctx.WriteString(" SCATTER")
	if node.From != nil {
		ctx.WriteString(" FROM (")
		ctx.FormatNode(&node.From)
		ctx.WriteString(") TO (")
		ctx.FormatNode(&node.To)
		ctx.WriteString(")")
	}
}
