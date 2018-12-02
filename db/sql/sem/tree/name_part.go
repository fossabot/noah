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
    "github.com/Ready-Stock/noah/db/sql/lex"
)

// A Name is an SQL identifier.
//
// In general, a Name is the result of parsing a name nonterminal, which is used
// in the grammar where reserved keywords cannot be distinguished from
// identifiers. A Name that matches a reserved keyword must thus be quoted when
// formatted. (Names also need quoting for a variety of other reasons; see
// isBareIdentifier.)
//
// For historical reasons, some Names are instead the result of parsing
// `unrestricted_name` nonterminals. See UnrestrictedName for details.
type Name string

// Format implements the NodeFormatter interface.
func (n *Name) Format(ctx *FmtCtx) {
	f := ctx.flags
	if f.HasFlags(FmtAnonymize) {
		ctx.WriteByte('_')
	} else {
		lex.EncodeRestrictedSQLIdent(ctx.Buffer, string(*n), f.EncodeFlags())
	}
}

// NameStringP escapes an identifier stored in a heap string to a SQL
// identifier, avoiding a heap allocation.
func NameStringP(s *string) string {
	return ((*Name)(s)).String()
}

// NameString escapes an identifier stored in a string to a SQL
// identifier.
func NameString(s string) string {
	return ((*Name)(&s)).String()
}

// ErrNameString escapes an identifier stored a string to a SQL
// identifier suitable for printing in error messages, avoiding a heap
// allocation.
func ErrNameString(s *string) string {
    return ErrString((*Name)(s))
}

// Normalize normalizes to lowercase and Unicode Normalization Form C
// (NFC).
func (n Name) Normalize() string {
	return lex.NormalizeName(string(n))
}

// An UnrestrictedName is a Name that does not need to be escaped when it
// matches a reserved keyword.
//
// In general, an UnrestrictedName is the result of parsing an unrestricted_name
// nonterminal, which is used in the grammar where reserved keywords can be
// unambiguously interpreted as identifiers. When formatted, an UnrestrictedName
// that matches a reserved keyword thus does not need to be quoted.
//
// For historical reasons, some unrestricted_name nonterminals are instead
// parsed as Names. The only user-visible impact of this is that we are too
// aggressive about quoting names in certain positions. New grammar rules should
// prefer to parse unrestricted_name nonterminals into UnrestrictedNames.
type UnrestrictedName string

// Format implements the NodeFormatter interface.
func (u *UnrestrictedName) Format(ctx *FmtCtx) {
	f := ctx.flags
	if f.HasFlags(FmtAnonymize) {
		ctx.WriteByte('_')
	} else {
		lex.EncodeUnrestrictedSQLIdent(ctx.Buffer, string(*u), f.EncodeFlags())
	}
}

// ToStrings converts the name list to an array of regular strings.
func (l NameList) ToStrings() []string {
	if l == nil {
		return nil
	}
	names := make([]string, len(l))
	for i, n := range l {
		names[i] = string(n)
	}
	return names
}

// A NameList is a list of identifiers.
type NameList []Name

// Format implements the NodeFormatter interface.
func (l *NameList) Format(ctx *FmtCtx) {
	for i := range *l {
		if i > 0 {
			ctx.WriteString(", ")
		}
		ctx.FormatNode(&(*l)[i])
	}
}

// ArraySubscript corresponds to the syntax `<name>[ ... ]`.
type ArraySubscript struct {
	Begin Expr
	End   Expr
	Slice bool
}

// Format implements the NodeFormatter interface.
func (a *ArraySubscript) Format(ctx *FmtCtx) {
	ctx.WriteByte('[')
	if a.Begin != nil {
		ctx.FormatNode(a.Begin)
	}
	if a.Slice {
		ctx.WriteByte(':')
		if a.End != nil {
			ctx.FormatNode(a.End)
		}
	}
	ctx.WriteByte(']')
}

// UnresolvedName corresponds to an unresolved qualified name.
type UnresolvedName struct {
	// NumParts indicates the number of name parts specified, including
	// the star. Always 1 or greater.
	NumParts int

	// Star indicates the name ends with a star.
	// In that case, Parts below is empty in the first position.
	Star bool

	// Parts are the name components, in reverse order.
	// There are at most 4: column, table, schema, catalog/db.
	//
	// Note: NameParts has a fixed size so that we avoid a heap
	// allocation for the slice every time we construct an
	// UnresolvedName. It does imply however that Parts does not have
	// a meaningful "length"; its actual length (the number of parts
	// specified) is populated in NumParts above.
	Parts NameParts
}

// NameParts is the array of strings that composes the path in an
// UnresolvedName.
type NameParts = [4]string

// Format implements the NodeFormatter interface.
func (u *UnresolvedName) Format(ctx *FmtCtx) {
	stopAt := 1
	if u.Star {
		stopAt = 2
	}
	// Every part after that is necessarily an unrestricted name.
	for i := u.NumParts; i >= stopAt; i-- {
		// The first part to print is the last item in u.Parts.  It is also
		// a potentially restricted name to disambiguate from keywords in
		// the grammar, so print it out as a "Name".
		if i == u.NumParts {
			ctx.FormatNode((*Name)(&u.Parts[i-1]))
		} else {
			ctx.FormatNode((*UnrestrictedName)(&u.Parts[i-1]))
		}
		if i > 1 {
			ctx.WriteByte('.')
		}
	}
	if u.Star {
		ctx.WriteByte('*')
	}
}
func (u *UnresolvedName) String() string { return AsString(u) }

// NewUnresolvedName constructs an UnresolvedName from some strings.
func NewUnresolvedName(args ...string) *UnresolvedName {
	n := MakeUnresolvedName(args...)
	return &n
}

// MakeUnresolvedName constructs an UnresolvedName from some strings.
func MakeUnresolvedName(args ...string) UnresolvedName {
	n := UnresolvedName{NumParts: len(args)}
	for i := 0; i < len(args); i++ {
		n.Parts[i] = args[len(args)-1-i]
	}
	return n
}
