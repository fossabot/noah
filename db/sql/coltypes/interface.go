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
	"fmt"

	"github.com/Ready-Stock/Noah/db/sql/lex"
)

// ColTypeFormatter knows how to format a ColType to a bytes.Buffer.
type ColTypeFormatter interface {
	fmt.Stringer

	// TypeName returns the base name of the type, suitable to generate
	// column names for cast expressions.
	TypeName() string

	// Format returns a non-lossy string representation of the coltype.
	// NOTE: It is important that two coltypes that should be different print out
	//       different string representations. The optimizer relies on unique
	//       string representations in order to intern the coltypes during
	//       memoization.
	Format(buf *bytes.Buffer, flags lex.EncodeFlags)
}

// ColTypeAsString print a T to a string.
func ColTypeAsString(n ColTypeFormatter) string {
	var buf bytes.Buffer
	n.Format(&buf, lex.EncNoFlags)
	return buf.String()
}

// CastTargetType represents a type that is a valid cast target.
type CastTargetType interface {
	ColTypeFormatter
	castTargetType()
}

// T represents a type in a column definition.
type T interface {
	CastTargetType

	columnType()
}

func (*TBool) columnType()           {}
func (*TInt) columnType()            {}
func (*TFloat) columnType()          {}
func (*TDecimal) columnType()        {}
func (*TDate) columnType()           {}
func (*TTime) columnType()           {}
func (*TTimeTZ) columnType()         {}
func (*TTimestamp) columnType()      {}
func (*TTimestampTZ) columnType()    {}
func (*TInterval) columnType()       {}
func (*TJSON) columnType()           {}
func (*TUUID) columnType()           {}
func (*TIPAddr) columnType()         {}
func (*TString) columnType()         {}
func (*TName) columnType()           {}
func (*TBytes) columnType()          {}
func (*TCollatedString) columnType() {}
func (*TArray) columnType()          {}
func (*TVector) columnType()         {}
func (TTuple) columnType()           {}
func (*TOid) columnType()            {}

// All Ts also implement CastTargetType.
func (*TBool) castTargetType()           {}
func (*TInt) castTargetType()            {}
func (*TFloat) castTargetType()          {}
func (*TDecimal) castTargetType()        {}
func (*TDate) castTargetType()           {}
func (*TTime) castTargetType()           {}
func (*TTimeTZ) castTargetType()         {}
func (*TTimestamp) castTargetType()      {}
func (*TTimestampTZ) castTargetType()    {}
func (*TInterval) castTargetType()       {}
func (*TJSON) castTargetType()           {}
func (*TUUID) castTargetType()           {}
func (*TIPAddr) castTargetType()         {}
func (*TString) castTargetType()         {}
func (*TName) castTargetType()           {}
func (*TBytes) castTargetType()          {}
func (*TCollatedString) castTargetType() {}
func (*TArray) castTargetType()          {}
func (*TVector) castTargetType()         {}
func (TTuple) castTargetType()           {}
func (*TOid) castTargetType()            {}

func (node *TBool) String() string           { return ColTypeAsString(node) }
func (node *TInt) String() string            { return ColTypeAsString(node) }
func (node *TFloat) String() string          { return ColTypeAsString(node) }
func (node *TDecimal) String() string        { return ColTypeAsString(node) }
func (node *TDate) String() string           { return ColTypeAsString(node) }
func (node *TTime) String() string           { return ColTypeAsString(node) }
func (node *TTimeTZ) String() string         { return ColTypeAsString(node) }
func (node *TTimestamp) String() string      { return ColTypeAsString(node) }
func (node *TTimestampTZ) String() string    { return ColTypeAsString(node) }
func (node *TInterval) String() string       { return ColTypeAsString(node) }
func (node *TJSON) String() string           { return ColTypeAsString(node) }
func (node *TUUID) String() string           { return ColTypeAsString(node) }
func (node *TIPAddr) String() string         { return ColTypeAsString(node) }
func (node *TString) String() string         { return ColTypeAsString(node) }
func (node *TName) String() string           { return ColTypeAsString(node) }
func (node *TBytes) String() string          { return ColTypeAsString(node) }
func (node *TCollatedString) String() string { return ColTypeAsString(node) }
func (node *TArray) String() string          { return ColTypeAsString(node) }
func (node *TVector) String() string         { return ColTypeAsString(node) }
func (node TTuple) String() string           { return ColTypeAsString(node) }
func (node *TOid) String() string            { return ColTypeAsString(node) }
