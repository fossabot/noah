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

package types

import (
	"database/sql/driver"
	"encoding/binary"
	"math"
	"strconv"

	"github.com/Ready-Stock/Noah/db/sql/pgio"
	"github.com/pkg/errors"
)

type Float4 struct {
	Float  float32
	Status Status
}

func (dst *Float4) GetFloat() float64 {
	return float64(dst.Float)
}

func (dst *Float4) Set(src interface{}) error {
	if src == nil {
		*dst = Float4{Status: Null}
		return nil
	}

	switch value := src.(type) {
	case float32:
		*dst = Float4{Float: value, Status: Present}
	case float64:
		*dst = Float4{Float: float32(value), Status: Present}
	case int8:
		*dst = Float4{Float: float32(value), Status: Present}
	case uint8:
		*dst = Float4{Float: float32(value), Status: Present}
	case int16:
		*dst = Float4{Float: float32(value), Status: Present}
	case uint16:
		*dst = Float4{Float: float32(value), Status: Present}
	case int32:
		f32 := float32(value)
		if int32(f32) == value {
			*dst = Float4{Float: f32, Status: Present}
		} else {
			return errors.Errorf("%v cannot be exactly represented as float32", value)
		}
	case uint32:
		f32 := float32(value)
		if uint32(f32) == value {
			*dst = Float4{Float: f32, Status: Present}
		} else {
			return errors.Errorf("%v cannot be exactly represented as float32", value)
		}
	case int64:
		f32 := float32(value)
		if int64(f32) == value {
			*dst = Float4{Float: f32, Status: Present}
		} else {
			return errors.Errorf("%v cannot be exactly represented as float32", value)
		}
	case uint64:
		f32 := float32(value)
		if uint64(f32) == value {
			*dst = Float4{Float: f32, Status: Present}
		} else {
			return errors.Errorf("%v cannot be exactly represented as float32", value)
		}
	case int:
		f32 := float32(value)
		if int(f32) == value {
			*dst = Float4{Float: f32, Status: Present}
		} else {
			return errors.Errorf("%v cannot be exactly represented as float32", value)
		}
	case uint:
		f32 := float32(value)
		if uint(f32) == value {
			*dst = Float4{Float: f32, Status: Present}
		} else {
			return errors.Errorf("%v cannot be exactly represented as float32", value)
		}
	case string:
		num, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}
		*dst = Float4{Float: float32(num), Status: Present}
	default:
		if originalSrc, ok := underlyingNumberType(src); ok {
			return dst.Set(originalSrc)
		}
		return errors.Errorf("cannot convert %v to Float8", value)
	}

	return nil
}

func (dst *Float4) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.Float
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Float4) AssignTo(dst interface{}) error {
	return float64AssignTo(float64(src.Float), src.Status, dst)
}

func (dst *Float4) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Float4{Status: Null}
		return nil
	}

	n, err := strconv.ParseFloat(string(src), 32)
	if err != nil {
		return err
	}

	*dst = Float4{Float: float32(n), Status: Present}
	return nil
}

func (dst *Float4) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Float4{Status: Null}
		return nil
	}

	if len(src) != 4 {
		return errors.Errorf("invalid length for float4: %v", len(src))
	}

	n := int32(binary.BigEndian.Uint32(src))

	*dst = Float4{Float: math.Float32frombits(uint32(n)), Status: Present}
	return nil
}

func (src *Float4) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = append(buf, strconv.FormatFloat(float64(src.Float), 'f', -1, 32)...)
	return buf, nil
}

func (src *Float4) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = pgio.AppendUint32(buf, math.Float32bits(src.Float))
	return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Float4) Scan(src interface{}) error {
	if src == nil {
		*dst = Float4{Status: Null}
		return nil
	}

	switch src := src.(type) {
	case float64:
		*dst = Float4{Float: float32(src), Status: Present}
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
func (src *Float4) Value() (driver.Value, error) {
	switch src.Status {
	case Present:
		return float64(src.Float), nil
	case Null:
		return nil, nil
	default:
		return nil, errUndefined
	}
}
