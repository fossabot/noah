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

package coltypes

import (
	"bytes"
	"fmt"

    "github.com/Ready-Stock/noah/db/sql/lex"
)

// TBool represents a BOOLEAN type.
type TBool struct {
	Name string
}

// TypeName implements the ColTypeFormatter interface.
func (node *TBool) TypeName() string { return node.Name }

// Format implements the ColTypeFormatter interface.
func (node *TBool) Format(buf *bytes.Buffer, f lex.EncodeFlags) {
	buf.WriteString(node.Name)
}

// TInt represents an INT, INTEGER, SMALLINT or BIGINT type.
type TInt struct {
	Name          string
	Width         int
	ImplicitWidth bool
}

// TypeName implements the ColTypeFormatter interface.
func (node *TInt) TypeName() string { return node.Name }

// Format implements the ColTypeFormatter interface.
func (node *TInt) Format(buf *bytes.Buffer, f lex.EncodeFlags) {
	buf.WriteString(node.Name)
	if node.Width > 0 && !node.ImplicitWidth {
		fmt.Fprintf(buf, "(%d)", node.Width)
	}
}

var serialIntTypes = map[string]struct{}{
	SmallSerial.Name: {},
	Serial.Name:      {},
	BigSerial.Name:   {},
	Serial2.Name:     {},
	Serial4.Name:     {},
	Serial8.Name:     {},
}

// IsSerial returns true when this column should be given a DEFAULT of a unique,
// incrementing function.
func (node *TInt) IsSerial() bool {
	_, ok := serialIntTypes[node.Name]
	return ok
}

// TFloat represents a REAL, DOUBLE or FLOAT type.
type TFloat struct {
	Name          string
	Prec          int
	Width         int
	PrecSpecified bool // true if the value of Prec is not the default
}

// TypeName implements the ColTypeFormatter interface.
func (node *TFloat) TypeName() string { return node.Name }

// Format implements the ColTypeFormatter interface.
func (node *TFloat) Format(buf *bytes.Buffer, f lex.EncodeFlags) {
	buf.WriteString(node.Name)
	if node.Prec > 0 {
		fmt.Fprintf(buf, "(%d)", node.Prec)
	}
}

// TDecimal represents a DECIMAL or NUMERIC type.
type TDecimal struct {
	Name  string
	Prec  int
	Scale int
}

// TypeName implements the ColTypeFormatter interface.
func (node *TDecimal) TypeName() string { return node.Name }

// Format implements the ColTypeFormatter interface.
func (node *TDecimal) Format(buf *bytes.Buffer, f lex.EncodeFlags) {
	buf.WriteString(node.Name)
	if node.Prec > 0 {
		fmt.Fprintf(buf, "(%d", node.Prec)
		if node.Scale > 0 {
			fmt.Fprintf(buf, ",%d", node.Scale)
		}
		buf.WriteByte(')')
	}
}
