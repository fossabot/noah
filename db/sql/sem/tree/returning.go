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

// ReturningClause represents the returning clause on a statement.
type ReturningClause interface {
	NodeFormatter
	// statementType returns the StatementType of statements that include
	// the implementors variant of a RETURNING clause.
	statementType() StatementType
	returningClause()
}

var _ ReturningClause = &ReturningExprs{}
var _ ReturningClause = &ReturningNothing{}
var _ ReturningClause = &NoReturningClause{}

// ReturningExprs represents RETURNING expressions.
type ReturningExprs SelectExprs

// Format implements the NodeFormatter interface.
func (r *ReturningExprs) Format(ctx *FmtCtx) {
	ctx.WriteString(" RETURNING ")
	ctx.FormatNode((*SelectExprs)(r))
}

// ReturningNothingClause is a shared instance to avoid unnecessary allocations.
var ReturningNothingClause = &ReturningNothing{}

// ReturningNothing represents RETURNING NOTHING.
type ReturningNothing struct{}

// Format implements the NodeFormatter interface.
func (*ReturningNothing) Format(ctx *FmtCtx) {
	ctx.WriteString(" RETURNING NOTHING")
}

// AbsentReturningClause is a ReturningClause variant representing the absence of
// a RETURNING clause.
var AbsentReturningClause = &NoReturningClause{}

// NoReturningClause represents the absence of a RETURNING clause.
type NoReturningClause struct{}

// Format implements the NodeFormatter interface.
func (*NoReturningClause) Format(_ *FmtCtx) {}

// used by parent statements to determine their own StatementType.
func (*ReturningExprs) statementType() StatementType    { return Rows }
func (*ReturningNothing) statementType() StatementType  { return RowsAffected }
func (*NoReturningClause) statementType() StatementType { return RowsAffected }

func (*ReturningExprs) returningClause()    {}
func (*ReturningNothing) returningClause()  {}
func (*NoReturningClause) returningClause() {}

// HasReturningClause determines if a ReturningClause is present, given a
// variant of the ReturningClause interface.
func HasReturningClause(clause ReturningClause) bool {
	_, ok := clause.(*NoReturningClause)
	return !ok
}
