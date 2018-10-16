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
	"fmt"
)

// Table patterns are used by e.g. GRANT statements, to designate
// zero, one or more table names.  For example:
//   GRANT ... ON foo ...
//   GRANT ... ON * ...
//   GRANT ... ON db.*  ...
//
// The other syntax nodes hold a TablePattern reference.  This is
// initially populated during parsing with an UnresolvedName, which
// can be transformed to either a TableName (single name) or
// AllTablesSelector instance (all tables of a given database) using
// NormalizeTablePattern().

// TablePattern is the common interface to UnresolvedName, TableName
// and AllTablesSelector.
type TablePattern interface {
	fmt.Stringer
	NodeFormatter

	// NormalizeTablePattern() guarantees to return a pattern that is
	// not an UnresolvedName. This converts the UnresolvedName to an
	// AllTablesSelector or TableName as necessary.
	NormalizeTablePattern() (TablePattern, error)
}

var _ TablePattern = &UnresolvedName{}
var _ TablePattern = &TableName{}
var _ TablePattern = &AllTablesSelector{}

// NormalizeTablePattern resolves an UnresolvedName to either a
// TableName or AllTablesSelector.
func (n *UnresolvedName) NormalizeTablePattern() (TablePattern, error) {
	return classifyTablePattern(n)
}

// NormalizeTablePattern implements the TablePattern interface.
func (t *TableName) NormalizeTablePattern() (TablePattern, error) { return t, nil }

// AllTablesSelector corresponds to a selection of all
// tables in a database, e.g. when used with GRANT.
type AllTablesSelector struct {
	TableNamePrefix
}

// Format implements the NodeFormatter interface.
func (at *AllTablesSelector) Format(ctx *FmtCtx) {
	at.TableNamePrefix.Format(ctx)
	if at.ExplicitSchema || ctx.alwaysFormatTablePrefix() {
		ctx.WriteByte('.')
	}
	ctx.WriteByte('*')
}
func (at *AllTablesSelector) String() string { return AsString(at) }

// NormalizeTablePattern implements the TablePattern interface.
func (at *AllTablesSelector) NormalizeTablePattern() (TablePattern, error) { return at, nil }

// TablePatterns implement a comma-separated list of table patterns.
// Used by e.g. the GRANT statement.
type TablePatterns []TablePattern

// Format implements the NodeFormatter interface.
func (tt *TablePatterns) Format(ctx *FmtCtx) {
	for i, t := range *tt {
		if i > 0 {
			ctx.WriteString(", ")
		}
		ctx.FormatNode(t)
	}
}
