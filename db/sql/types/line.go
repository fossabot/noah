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

package types

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Ready-Stock/Noah/db/sql/pgio"
	"github.com/pkg/errors"
)

type Line struct {
	A, B, C float64
	Status  Status
}

func (dst *Line) Set(src interface{}) error {
	return errors.Errorf("cannot convert %v to Line", src)
}

func (dst *Line) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Line) AssignTo(dst interface{}) error {
	return errors.Errorf("cannot assign %v to %T", src, dst)
}

func (dst *Line) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Line{Status: Null}
		return nil
	}

	if len(src) < 7 {
		return errors.Errorf("invalid length for Line: %v", len(src))
	}

	parts := strings.SplitN(string(src[1:len(src)-1]), ",", 3)
	if len(parts) < 3 {
		return errors.Errorf("invalid format for line")
	}

	a, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return err
	}

	b, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return err
	}

	c, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return err
	}

	*dst = Line{A: a, B: b, C: c, Status: Present}
	return nil
}

func (dst *Line) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Line{Status: Null}
		return nil
	}

	if len(src) != 24 {
		return errors.Errorf("invalid length for Line: %v", len(src))
	}

	a := binary.BigEndian.Uint64(src)
	b := binary.BigEndian.Uint64(src[8:])
	c := binary.BigEndian.Uint64(src[16:])

	*dst = Line{
		A:      math.Float64frombits(a),
		B:      math.Float64frombits(b),
		C:      math.Float64frombits(c),
		Status: Present,
	}
	return nil
}

func (src *Line) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = append(buf, fmt.Sprintf(`{%s,%s,%s}`,
		strconv.FormatFloat(src.A, 'f', -1, 64),
		strconv.FormatFloat(src.B, 'f', -1, 64),
		strconv.FormatFloat(src.C, 'f', -1, 64),
	)...)

	return buf, nil
}

func (src *Line) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = pgio.AppendUint64(buf, math.Float64bits(src.A))
	buf = pgio.AppendUint64(buf, math.Float64bits(src.B))
	buf = pgio.AppendUint64(buf, math.Float64bits(src.C))
	return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Line) Scan(src interface{}) error {
	if src == nil {
		*dst = Line{Status: Null}
		return nil
	}

	switch src := src.(type) {
	case string:
		return dst.DecodeText(nil, []byte(src))
	case []byte:
		srcCopy := make([]byte, len(src))
		copy(srcCopy, src)
		return dst.DecodeText(nil, srcCopy)
	}

	return errors.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src *Line) Value() (driver.Value, error) {
	return EncodeValueText(src)
}
