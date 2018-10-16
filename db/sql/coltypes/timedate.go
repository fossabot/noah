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

package coltypes

import (
	"bytes"

	"github.com/Ready-Stock/Noah/db/sql/lex"
)

// TDate represents a DATE type.
type TDate struct{}

// TypeName implements the ColTypeFormatter interface.
func (node *TDate) TypeName() string { return "DATE" }

// Format implements the ColTypeFormatter interface.
func (node *TDate) Format(buf *bytes.Buffer, f lex.EncodeFlags) {
	buf.WriteString("DATE")
}

// TTime represents a TIME type.
type TTime struct{}

// TypeName implements the ColTypeFormatter interface.
func (node *TTime) TypeName() string { return "TIME" }

// Format implements the ColTypeFormatter interface.
func (node *TTime) Format(buf *bytes.Buffer, f lex.EncodeFlags) {
	buf.WriteString("TIME")
}

// TTimeTZ represents a TIMETZ type.
type TTimeTZ struct{}

// TypeName implements the ColTypeFormatter interface.
func (node *TTimeTZ) TypeName() string { return "TIMETZ" }

// Format implements the ColTypeFormatter interface.
func (node *TTimeTZ) Format(buf *bytes.Buffer, f lex.EncodeFlags) {
	buf.WriteString("TIME WITH TIME ZONE")
}

// TTimestamp represents a TIMESTAMP type.
type TTimestamp struct{}

// TypeName implements the ColTypeFormatter interface.
func (node *TTimestamp) TypeName() string { return "TIMESTAMP" }

// Format implements the ColTypeFormatter interface.
func (node *TTimestamp) Format(buf *bytes.Buffer, f lex.EncodeFlags) {
	buf.WriteString("TIMESTAMP")
}

// TTimestampTZ represents a TIMESTAMP type.
type TTimestampTZ struct{}

// TypeName implements the ColTypeFormatter interface.
func (node *TTimestampTZ) TypeName() string { return "TIMESTAMPTZ" }

// Format implements the ColTypeFormatter interface.
func (node *TTimestampTZ) Format(buf *bytes.Buffer, f lex.EncodeFlags) {
	buf.WriteString("TIMESTAMP WITH TIME ZONE")
}

// TInterval represents an INTERVAL type
type TInterval struct{}

// TypeName implements the ColTypeFormatter interface.
func (node *TInterval) TypeName() string { return "INTERVAL" }

// Format implements the ColTypeFormatter interface.
func (node *TInterval) Format(buf *bytes.Buffer, f lex.EncodeFlags) {
	buf.WriteString("INTERVAL")
}
