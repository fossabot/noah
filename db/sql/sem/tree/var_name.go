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
	"github.com/Ready-Stock/Noah/db/sql/sem/types"
)

// VarName occurs inside scalar expressions.
//
// Immediately after parsing, the following types can occur:
//
// - UnqualifiedStar: a naked star as argument to a function, e.g. count(*),
//   or at the top level of a SELECT clause.
//   See also uses of StarExpr() and StarSelectExpr() in the grammar.
//
// - UnresolvedName: other names of the form `a.b....e` or `a.b...e.*`.
//
// Consumers of variable names do not like UnresolvedNames and instead
// expect either AllColumnsSelector or ColumnItem. Use
// NormalizeVarName() for this.
//
// After a ColumnItem is available, it should be further resolved, for this
// the Resolve() method should be used; see name_resolution.go.
type VarName interface {
	TypedExpr

	// NormalizeVarName() guarantees to return a variable name
	// that is not an UnresolvedName. This converts the UnresolvedName
	// to an AllColumnsSelector or ColumnItem as necessary.
	NormalizeVarName() (VarName, error)
}

var _ VarName = &UnresolvedName{}
var _ VarName = UnqualifiedStar{}
var _ VarName = &AllColumnsSelector{}
var _ VarName = &TupleStar{}
var _ VarName = &ColumnItem{}

// UnqualifiedStar corresponds to a standalone '*' in a scalar
// expression.
type UnqualifiedStar struct{}

// Format implements the NodeFormatter interface.
func (UnqualifiedStar) Format(ctx *FmtCtx) { ctx.WriteByte('*') }
func (u UnqualifiedStar) String() string   { return AsString(u) }

// NormalizeVarName implements the VarName interface.
func (u UnqualifiedStar) NormalizeVarName() (VarName, error) { return u, nil }

var singletonStarName VarName = UnqualifiedStar{}

// StarExpr is a convenience function that represents an unqualified "*".
func StarExpr() VarName { return singletonStarName }

// ResolvedType implements the TypedExpr interface.
func (UnqualifiedStar) ResolvedType() types.T {
	panic("unqualified stars ought to be replaced before this point")
}

// Variable implements the VariableExpr interface.
func (UnqualifiedStar) Variable() {}

// UnresolvedName is defined in name_part.go. It also implements the
// VarName interface, and thus TypedExpr too.

// ResolvedType implements the TypedExpr interface.
func (*UnresolvedName) ResolvedType() types.T {
	panic("unresolved names ought to be replaced before this point")
}

// Variable implements the VariableExpr interface.  Although, the
// UnresolvedName ought to be replaced to an IndexedVar before the points the
// VariableExpr interface is used.
func (*UnresolvedName) Variable() {}

// NormalizeVarName implements the VarName interface.
func (n *UnresolvedName) NormalizeVarName() (VarName, error) {
	return classifyColumnItem(n)
}

// AllColumnsSelector corresponds to a selection of all
// columns in a table when used in a SELECT clause.
// (e.g. `table.*`).
type AllColumnsSelector struct {
	// TableName corresponds to the table prefix, before the star. The
	// UnresolvedName within is guaranteed to not contain a star itself.
	TableName UnresolvedName
}

// Format implements the NodeFormatter interface.
func (a *AllColumnsSelector) Format(ctx *FmtCtx) {
	ctx.FormatNode(&a.TableName)
	ctx.WriteString(".*")
}
func (a *AllColumnsSelector) String() string { return AsString(a) }

// NormalizeVarName implements the VarName interface.
func (a *AllColumnsSelector) NormalizeVarName() (VarName, error) { return a, nil }

// Variable implements the VariableExpr interface.  Although, the
// AllColumnsSelector ought to be replaced to an IndexedVar before the points the
// VariableExpr interface is used.
func (a *AllColumnsSelector) Variable() {}

// ResolvedType implements the TypedExpr interface.
func (*AllColumnsSelector) ResolvedType() types.T {
	panic("all-columns selectors ought to be replaced before this point")
}

// ColumnItem corresponds to the name of a column in an expression.
type ColumnItem struct {
	// TableName holds the table prefix, if the name refers to a column.
	//
	// This uses UnresolvedName because we need to preserve the
	// information about which parts were initially specified in the SQL
	// text. ColumnItems are intermediate data structures anyway, that
	// still need to undergo name resolution.
	TableName UnresolvedName
	// ColumnName names the designated column.
	ColumnName Name

	// This column is a selector column expression used in a SELECT
	// for an UPDATE/DELETE.
	// TODO(vivek): Do not artificially create such expressions
	// when scanning columns for an UPDATE/DELETE.
	ForUpdateOrDelete bool
}

// Format implements the NodeFormatter interface.
func (c *ColumnItem) Format(ctx *FmtCtx) {
	if c.TableName.NumParts > 0 {
		c.TableName.Format(ctx)
		ctx.WriteByte('.')
	}
	ctx.FormatNode(&c.ColumnName)
}
func (c *ColumnItem) String() string { return AsString(c) }

// NormalizeVarName implements the VarName interface.
func (c *ColumnItem) NormalizeVarName() (VarName, error) { return c, nil }

// Column retrieves the unqualified column name.
func (c *ColumnItem) Column() string {
	return string(c.ColumnName)
}

// Variable implements the VariableExpr interface.
//
// Note that in common uses, ColumnItem ought to be replaced to an
// IndexedVar prior to evaluation.
func (c *ColumnItem) Variable() {}

// ResolvedType implements the TypedExpr interface.
func (c *ColumnItem) ResolvedType() types.T {
	if presetTypesForTesting == nil {
		return nil
	}
	return presetTypesForTesting[c.String()]
}

// NewColumnItem constructs a column item from an already valid
// TableName. This can be used for e.g. pretty-printing.
func NewColumnItem(tn *TableName, colName Name) *ColumnItem {
	c := MakeColumnItem(tn, colName)
	return &c
}

// MakeColumnItem constructs a column item from an already valid
// TableName. This can be used for e.g. pretty-printing.
func MakeColumnItem(tn *TableName, colName Name) ColumnItem {
	c := ColumnItem{
		TableName: UnresolvedName{
			Parts: NameParts{tn.Table(), tn.Schema(), tn.Catalog()},
		},
		ColumnName: colName,
	}
	if tn.ExplicitCatalog {
		c.TableName.NumParts = 3
	} else if tn.ExplicitSchema {
		c.TableName.NumParts = 2
	} else {
		c.TableName.NumParts = 1
	}
	return c
}
