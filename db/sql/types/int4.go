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

package types

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"math"
	"strconv"

	"github.com/Ready-Stock/Noah/db/sql/pgio"
	"github.com/pkg/errors"
)

type Int4 struct {
	Int    int32
	Status Status
}

func (dst *Int4) GetInt() int64 {
	return int64(dst.Int)
}

func (dst *Int4) Set(src interface{}) error {
	if src == nil {
		*dst = Int4{Status: Null}
		return nil
	}

	switch value := src.(type) {
	case int8:
		*dst = Int4{Int: int32(value), Status: Present}
	case uint8:
		*dst = Int4{Int: int32(value), Status: Present}
	case int16:
		*dst = Int4{Int: int32(value), Status: Present}
	case uint16:
		*dst = Int4{Int: int32(value), Status: Present}
	case int32:
		*dst = Int4{Int: int32(value), Status: Present}
	case uint32:
		if value > math.MaxInt32 {
			return errors.Errorf("%d is greater than maximum value for Int4", value)
		}
		*dst = Int4{Int: int32(value), Status: Present}
	case int64:
		if value < math.MinInt32 {
			return errors.Errorf("%d is greater than maximum value for Int4", value)
		}
		if value > math.MaxInt32 {
			return errors.Errorf("%d is greater than maximum value for Int4", value)
		}
		*dst = Int4{Int: int32(value), Status: Present}
	case uint64:
		if value > math.MaxInt32 {
			return errors.Errorf("%d is greater than maximum value for Int4", value)
		}
		*dst = Int4{Int: int32(value), Status: Present}
	case int:
		if value < math.MinInt32 {
			return errors.Errorf("%d is greater than maximum value for Int4", value)
		}
		if value > math.MaxInt32 {
			return errors.Errorf("%d is greater than maximum value for Int4", value)
		}
		*dst = Int4{Int: int32(value), Status: Present}
	case uint:
		if value > math.MaxInt32 {
			return errors.Errorf("%d is greater than maximum value for Int4", value)
		}
		*dst = Int4{Int: int32(value), Status: Present}
	case string:
		num, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}
		*dst = Int4{Int: int32(num), Status: Present}
	default:
		if originalSrc, ok := underlyingNumberType(src); ok {
			return dst.Set(originalSrc)
		}
		return errors.Errorf("cannot convert %v to Int4", value)
	}

	return nil
}

func (dst *Int4) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.Int
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Int4) AssignTo(dst interface{}) error {
	return int64AssignTo(int64(src.Int), src.Status, dst)
}

func (dst *Int4) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Int4{Status: Null}
		return nil
	}

	n, err := strconv.ParseInt(string(src), 10, 32)
	if err != nil {
		return err
	}

	*dst = Int4{Int: int32(n), Status: Present}
	return nil
}

func (dst *Int4) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Int4{Status: Null}
		return nil
	}

	if len(src) != 4 {
		return errors.Errorf("invalid length for int4: %v", len(src))
	}

	n := int32(binary.BigEndian.Uint32(src))
	*dst = Int4{Int: n, Status: Present}
	return nil
}

func (src *Int4) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	return append(buf, strconv.FormatInt(int64(src.Int), 10)...), nil
}

func (src *Int4) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	return pgio.AppendInt32(buf, src.Int), nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Int4) Scan(src interface{}) error {
	if src == nil {
		*dst = Int4{Status: Null}
		return nil
	}

	switch src := src.(type) {
	case int64:
		if src < math.MinInt32 {
			return errors.Errorf("%d is greater than maximum value for Int4", src)
		}
		if src > math.MaxInt32 {
			return errors.Errorf("%d is greater than maximum value for Int4", src)
		}
		*dst = Int4{Int: int32(src), Status: Present}
		return nil
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
func (src *Int4) Value() (driver.Value, error) {
	switch src.Status {
	case Present:
		return int64(src.Int), nil
	case Null:
		return nil, nil
	default:
		return nil, errUndefined
	}
}

func (src *Int4) MarshalJSON() ([]byte, error) {
	switch src.Status {
	case Present:
		return []byte(strconv.FormatInt(int64(src.Int), 10)), nil
	case Null:
		return []byte("null"), nil
	case Undefined:
		return nil, errUndefined
	}

	return nil, errBadStatus
}

func (dst *Int4) UnmarshalJSON(b []byte) error {
	var n int32
	err := json.Unmarshal(b, &n)
	if err != nil {
		return err
	}

	*dst = Int4{Int: n, Status: Present}

	return nil
}
