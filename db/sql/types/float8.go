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
 * Copyright 2018  'pq' Contributors Portions Copyright (C) 2018 Blake Mizerany
 * License https://github.com/lib/pq/blob/master/LICENSE.md
 *
 * Project: go-linq github.com/ahmetb/go-linq
 * Copyright 2018 Ahmet Alp Balkan
 * License https://github.com/ahmetb/go-linq/blob/master/LICENSE
 */

package types

import (
    "database/sql/driver"
    "encoding/binary"
    "math"
    "strconv"

    "github.com/pkg/errors"
    "github.com/readystock/noah/db/sql/pgio"
)

type Float8 struct {
    Float  float64
    Status Status
}

func (dst *Float8) GetFloat() float64 {
    return float64(dst.Float)
}

func (dst *Float8) Set(src interface{}) error {
    if src == nil {
        *dst = Float8{Status: Null}
        return nil
    }

    switch value := src.(type) {
    case float32:
        *dst = Float8{Float: float64(value), Status: Present}
    case float64:
        *dst = Float8{Float: value, Status: Present}
    case int8:
        *dst = Float8{Float: float64(value), Status: Present}
    case uint8:
        *dst = Float8{Float: float64(value), Status: Present}
    case int16:
        *dst = Float8{Float: float64(value), Status: Present}
    case uint16:
        *dst = Float8{Float: float64(value), Status: Present}
    case int32:
        *dst = Float8{Float: float64(value), Status: Present}
    case uint32:
        *dst = Float8{Float: float64(value), Status: Present}
    case int64:
        f64 := float64(value)
        if int64(f64) == value {
            *dst = Float8{Float: f64, Status: Present}
        } else {
            return errors.Errorf("%v cannot be exactly represented as float64", value)
        }
    case uint64:
        f64 := float64(value)
        if uint64(f64) == value {
            *dst = Float8{Float: f64, Status: Present}
        } else {
            return errors.Errorf("%v cannot be exactly represented as float64", value)
        }
    case int:
        f64 := float64(value)
        if int(f64) == value {
            *dst = Float8{Float: f64, Status: Present}
        } else {
            return errors.Errorf("%v cannot be exactly represented as float64", value)
        }
    case uint:
        f64 := float64(value)
        if uint(f64) == value {
            *dst = Float8{Float: f64, Status: Present}
        } else {
            return errors.Errorf("%v cannot be exactly represented as float64", value)
        }
    case string:
        num, err := strconv.ParseFloat(value, 64)
        if err != nil {
            return err
        }
        *dst = Float8{Float: float64(num), Status: Present}
    default:
        if originalSrc, ok := underlyingNumberType(src); ok {
            return dst.Set(originalSrc)
        }
        return errors.Errorf("cannot convert %v to Float8", value)
    }

    return nil
}

func (dst *Float8) Get() interface{} {
    switch dst.Status {
    case Present:
        return dst.Float
    case Null:
        return nil
    default:
        return dst.Status
    }
}

func (src *Float8) AssignTo(dst interface{}) error {
    return float64AssignTo(src.Float, src.Status, dst)
}

func (dst *Float8) DecodeText(ci *ConnInfo, src []byte) error {
    if src == nil {
        *dst = Float8{Status: Null}
        return nil
    }

    n, err := strconv.ParseFloat(string(src), 64)
    if err != nil {
        return err
    }

    *dst = Float8{Float: n, Status: Present}
    return nil
}

func (dst *Float8) DecodeBinary(ci *ConnInfo, src []byte) error {
    if src == nil {
        *dst = Float8{Status: Null}
        return nil
    }

    if len(src) != 8 {
        return errors.Errorf("invalid length for float4: %v", len(src))
    }

    n := int64(binary.BigEndian.Uint64(src))

    *dst = Float8{Float: math.Float64frombits(uint64(n)), Status: Present}
    return nil
}

func (src *Float8) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
    switch src.Status {
    case Null:
        return nil, nil
    case Undefined:
        return nil, errUndefined
    }

    buf = append(buf, strconv.FormatFloat(float64(src.Float), 'f', -1, 64)...)
    return buf, nil
}

func (src *Float8) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
    switch src.Status {
    case Null:
        return nil, nil
    case Undefined:
        return nil, errUndefined
    }

    buf = pgio.AppendUint64(buf, math.Float64bits(src.Float))
    return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Float8) Scan(src interface{}) error {
    if src == nil {
        *dst = Float8{Status: Null}
        return nil
    }

    switch src := src.(type) {
    case float64:
        *dst = Float8{Float: src, Status: Present}
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
func (src *Float8) Value() (driver.Value, error) {
    switch src.Status {
    case Present:
        return src.Float, nil
    case Null:
        return nil, nil
    default:
        return nil, errUndefined
    }
}
