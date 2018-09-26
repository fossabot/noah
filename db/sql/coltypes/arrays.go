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

package coltypes

import (
	"bytes"

	"github.com/Ready-Stock/Noah/db/sql/lex"
)

// TArray represents an ARRAY column type.
type TArray struct {
	Name string
	// ParamTyp is the type of the elements in this array.
	ParamType T
	Bounds    []int32
}

// TypeName implements the ColTypeFormatter interface.
func (node *TArray) TypeName() string { return node.Name }

// Format implements the ColTypeFormatter interface.
func (node *TArray) Format(buf *bytes.Buffer, f lex.EncodeFlags) {
	buf.WriteString(node.Name)
	if collation, ok := node.ParamType.(*TCollatedString); ok {
		buf.WriteString(" COLLATE ")
		lex.EncodeUnrestrictedSQLIdent(buf, collation.Locale, f)
	}
}

// canBeInArrayColType returns true if the given T is a valid
// element type for an array column type.
func canBeInArrayColType(t T) bool {
	switch t.(type) {
	case *TJSON:
		return false
	default:
		return true
	}
}

// TVector is the base for VECTOR column types, which are Postgres's
// older, limited version of ARRAYs. These are not meant to be persisted,
// because ARRAYs are a strict superset.
type TVector struct {
	Name      string
	ParamType T
}

// TypeName implements the ColTypeFormatter interface.
func (node *TVector) TypeName() string { return node.Name }

// Format implements the ColTypeFormatter interface.
func (node *TVector) Format(buf *bytes.Buffer, _ lex.EncodeFlags) {
	buf.WriteString(node.Name)
}
