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
 */

package tree

import (
	"fmt"
)

// ScrubType describes the SCRUB statement operation.
type ScrubType int

const (
	// ScrubTable describes the SCRUB operation SCRUB TABLE.
	ScrubTable = iota
	// ScrubDatabase describes the SCRUB operation SCRUB DATABASE.
	ScrubDatabase = iota
)

// Scrub represents a SCRUB statement.
type Scrub struct {
	Typ     ScrubType
	Options ScrubOptions
	// Table is only set during SCRUB TABLE statements.
	Table NormalizableTableName
	// Database is only set during SCRUB DATABASE statements.
	Database Name
	AsOf     AsOfClause
}

// Format implements the NodeFormatter interface.
func (n *Scrub) Format(ctx *FmtCtx) {
	ctx.WriteString("EXPERIMENTAL SCRUB ")
	switch n.Typ {
	case ScrubTable:
		ctx.WriteString("TABLE ")
		n.Table.Format(ctx)
	case ScrubDatabase:
		ctx.WriteString("DATABASE ")
		ctx.FormatNode(&n.Database)
	default:
		panic("Unhandled ScrubType")
	}

	if n.AsOf.Expr != nil {
		ctx.WriteByte(' ')
		ctx.FormatNode(&n.AsOf)
	}

	if len(n.Options) > 0 {
		ctx.WriteString(" WITH OPTIONS ")
		ctx.FormatNode(&n.Options)
	}
}

// ScrubOptions corresponds to a comma-delimited list of scrub options.
type ScrubOptions []ScrubOption

// Format implements the NodeFormatter interface.
func (n *ScrubOptions) Format(ctx *FmtCtx) {
	for i, option := range *n {
		if i > 0 {
			ctx.WriteString(", ")
		}
		ctx.FormatNode(option)
	}
}

func (n *ScrubOptions) String() string { return AsString(n) }

// ScrubOption represents a scrub option.
type ScrubOption interface {
	fmt.Stringer
	NodeFormatter

	scrubOptionType()
}

// scrubOptionType implements the ScrubOption interface
func (*ScrubOptionIndex) scrubOptionType()      {}
func (*ScrubOptionPhysical) scrubOptionType()   {}
func (*ScrubOptionConstraint) scrubOptionType() {}

func (n *ScrubOptionIndex) String() string      { return AsString(n) }
func (n *ScrubOptionPhysical) String() string   { return AsString(n) }
func (n *ScrubOptionConstraint) String() string { return AsString(n) }

// ScrubOptionIndex represents an INDEX scrub check.
type ScrubOptionIndex struct {
	IndexNames NameList
}

// Format implements the NodeFormatter interface.
func (n *ScrubOptionIndex) Format(ctx *FmtCtx) {
	ctx.WriteString("INDEX ")
	if n.IndexNames != nil {
		ctx.WriteByte('(')
		ctx.FormatNode(&n.IndexNames)
		ctx.WriteByte(')')
	} else {
		ctx.WriteString("ALL")
	}
}

// ScrubOptionPhysical represents a PHYSICAL scrub check.
type ScrubOptionPhysical struct{}

// Format implements the NodeFormatter interface.
func (n *ScrubOptionPhysical) Format(ctx *FmtCtx) {
	ctx.WriteString("PHYSICAL")
}

// ScrubOptionConstraint represents a CONSTRAINT scrub check.
type ScrubOptionConstraint struct {
	ConstraintNames NameList
}

// Format implements the NodeFormatter interface.
func (n *ScrubOptionConstraint) Format(ctx *FmtCtx) {
	ctx.WriteString("CONSTRAINT ")
	if n.ConstraintNames != nil {
		ctx.WriteByte('(')
		ctx.FormatNode(&n.ConstraintNames)
		ctx.WriteByte(')')
	} else {
		ctx.WriteString("ALL")
	}
}
