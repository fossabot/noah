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
    "bytes"
    "fmt"
    "sort"
    "strings"

    "github.com/readystock/noah/db/sql/pgwire/pgerror"
    "github.com/readystock/noah/db/sql/sem/types"
    "github.com/readystock/noah/db/util"
)

// PlaceholderTypes relates placeholder names to their resolved type.
type PlaceholderTypes map[string]types.T

// QueryArguments relates placeholder names to their provided query argument.
//
// A nil value represents a NULL argument.
type QueryArguments map[string]TypedExpr

var emptyQueryArgumentStr = "{}"

func (qa *QueryArguments) String() string {
    if len(*qa) == 0 {
        return emptyQueryArgumentStr
    }
    var buf bytes.Buffer
    buf.WriteByte('{')
    sep := ""
    for k, v := range *qa {
        fmt.Fprintf(&buf, "%s$%s:%q", sep, k, v)
        sep = ", "
    }
    buf.WriteByte('}')
    return buf.String()
}

// PlaceholderInfo defines the interface to SQL placeholders.
type PlaceholderInfo struct {
    Values QueryArguments
    // TypeHints contains the initially set type hints for each placeholder if
    // present, and will be filled in completely by the end of type checking
    // Hints that were present before type checking will not change, and hints
    // that were not present before type checking will be set to their
    // placeholder's inferred type.
    TypeHints PlaceholderTypes
    // Types contains the final types set for each placeholder after type
    // checking.
    Types PlaceholderTypes
    // permitUnassigned controls whether AssertAllAssigned returns an error when
    // there are unassigned placeholders. See PermitUnassigned().
    permitUnassigned bool
}

// MakePlaceholderInfo constructs an empty PlaceholderInfo.
func MakePlaceholderInfo() PlaceholderInfo {
    res := PlaceholderInfo{}
    res.Clear()
    return res
}

// Clear resets the placeholder info map.
func (p *PlaceholderInfo) Clear() {
    p.TypeHints = PlaceholderTypes{}
    p.Types = PlaceholderTypes{}
    p.Values = QueryArguments{}
    p.permitUnassigned = false
}

// Assign resets the PlaceholderInfo to the contents of src.
// If src is nil, the map is cleared.
func (p *PlaceholderInfo) Assign(src *PlaceholderInfo) {
    if src != nil {
        *p = *src
    } else {
        p.Clear()
    }
}

// PermitUnassigned permits unassigned placeholders during plan construction,
// so that EXPLAIN can work on statements with placeholders.
func (p *PlaceholderInfo) PermitUnassigned() {
    p.permitUnassigned = true
}

// AssertAllAssigned ensures that all placeholders that are used also have a
// value assigned, or that PermitUnassigned was called.
func (p *PlaceholderInfo) AssertAllAssigned() error {
    if p.permitUnassigned {
        return nil
    }
    var missing []string
    for pn := range p.Types {
        if _, ok := p.Values[pn]; !ok {
            missing = append(missing, "$"+pn)
        }
    }
    if len(missing) > 0 {
        sort.Strings(missing)
        return pgerror.NewErrorf(pgerror.CodeUndefinedParameterError,
            "no value provided for placeholder%s: %s",
            util.Pluralize(int64(len(missing))),
            strings.Join(missing, ", "),
        )
    }
    return nil
}

// Type returns the known type of a placeholder. If allowHints is true, will
// return a type hint if there's no known type yet but there is a type hint.
// Returns false in the 2nd value if the placeholder is not typed.
func (p *PlaceholderInfo) Type(name string, allowHints bool) (types.T, bool) {
    if t, ok := p.Types[name]; ok {
        return t, true
    } else if t, ok := p.TypeHints[name]; ok {
        return t, true
    }
    return nil, false
}

// Value returns the known value of a placeholder.  Returns false in
// the 2nd value if the placeholder does not have a value.
func (p *PlaceholderInfo) Value(name string) (TypedExpr, bool) {
    if v, ok := p.Values[name]; ok {
        return v, true
    }
    return nil, false
}

// SetType assigns a known type to a placeholder.
// Reports an error if another type was previously assigned.
func (p *PlaceholderInfo) SetType(name string, typ types.T) error {
    if t, ok := p.Types[name]; ok {
        if !typ.Equivalent(t) {
            return pgerror.NewErrorf(
                pgerror.CodeDatatypeMismatchError,
                "placeholder %s already has type %s, cannot assign %s", name, t, typ)
        }
        return nil
    }
    p.Types[name] = typ
    if _, ok := p.TypeHints[name]; !ok {
        // If the client didn't give us a type hint, we must communicate our
        // inferred type to pgwire so it can know how to parse incoming data.
        p.TypeHints[name] = typ
    }
    return nil
}

// SetTypeHints resets the type and values in the map and replaces the
// type hints map by an alias to src. If src is nil, the map is cleared.
// The type hints map is aliased because the invoking code from
// pgwire/v3.go for sql.Prepare needs to receive the updated type
// assignments after Prepare completes.
func (p *PlaceholderInfo) SetTypeHints(src PlaceholderTypes) {
    if src != nil {
        p.TypeHints = src
        p.Types = PlaceholderTypes{}
        p.Values = QueryArguments{}
    } else {
        p.Clear()
    }
}

// IsUnresolvedPlaceholder returns whether expr is an unresolved placeholder. In
// other words, it returns whether the provided expression is a placeholder
// expression or a placeholder expression within nested parentheses, and if so,
// whether the placeholder's type remains unset in the PlaceholderInfo.
func (p *PlaceholderInfo) IsUnresolvedPlaceholder(expr Expr) bool {
    if t, ok := StripParens(expr).(*Placeholder); ok {
        _, res := p.TypeHints[t.Name]
        return !res
    }
    return false
}
