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

// This code was derived from https://github.com/youtube/vitess.

package tree

// Insert represents an INSERT statement.
type Insert struct {
	With       *With
	Table      TableExpr
	Columns    NameList
	Rows       *Select
	OnConflict *OnConflict
	Returning  ReturningClause
}

// Format implements the NodeFormatter interface.
func (node *Insert) Format(ctx *FmtCtx) {
	ctx.FormatNode(node.With)
	if node.OnConflict.IsUpsertAlias() {
		ctx.WriteString("UPSERT")
	} else {
		ctx.WriteString("INSERT")
	}
	ctx.WriteString(" INTO ")
	ctx.FormatNode(node.Table)
	if node.Columns != nil {
		ctx.WriteByte('(')
		ctx.FormatNode(&node.Columns)
		ctx.WriteByte(')')
	}
	if node.DefaultValues() {
		ctx.WriteString(" DEFAULT VALUES")
	} else {
		ctx.WriteByte(' ')
		ctx.FormatNode(node.Rows)
	}
	if node.OnConflict != nil && !node.OnConflict.IsUpsertAlias() {
		ctx.WriteString(" ON CONFLICT")
		if len(node.OnConflict.Columns) > 0 {
			ctx.WriteString(" (")
			ctx.FormatNode(&node.OnConflict.Columns)
			ctx.WriteString(")")
		}
		if node.OnConflict.DoNothing {
			ctx.WriteString(" DO NOTHING")
		} else {
			ctx.WriteString(" DO UPDATE SET ")
			ctx.FormatNode(&node.OnConflict.Exprs)
			if node.OnConflict.Where != nil {
				ctx.FormatNode(node.OnConflict.Where)
			}
		}
	}
	ctx.FormatNode(node.Returning)
}

// DefaultValues returns true iff only default values are being inserted.
func (node *Insert) DefaultValues() bool {
	return node.Rows.Select == nil
}

// OnConflict represents an `ON CONFLICT (columns) DO UPDATE SET exprs WHERE
// where` clause.
//
// The zero value for OnConflict is used to signal the UPSERT short form, which
// uses the primary key for as the conflict index and the values being inserted
// for Exprs.
type OnConflict struct {
	Columns   NameList
	Exprs     UpdateExprs
	Where     *Where
	DoNothing bool
}

// IsUpsertAlias returns true if the UPSERT syntactic sugar was used.
func (oc *OnConflict) IsUpsertAlias() bool {
	return oc != nil && oc.Columns == nil && oc.Exprs == nil && oc.Where == nil && !oc.DoNothing
}
