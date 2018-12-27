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
    "fmt"
    "github.com/pkg/errors"
    "github.com/readystock/noah/db/sql/pgio"
    "math"
    "strconv"
    "strings"
)

type Circle struct {
    P      Vec2
    R      float64
    Status Status
}

func (dst *Circle) Set(src interface{}) error {
    return errors.Errorf("cannot convert %v to Circle", src)
}

func (dst *Circle) Get() interface{} {
    switch dst.Status {
    case Present:
        return dst
    case Null:
        return nil
    default:
        return dst.Status
    }
}

func (src *Circle) AssignTo(dst interface{}) error {
    return errors.Errorf("cannot assign %v to %T", src, dst)
}

func (dst *Circle) DecodeText(ci *ConnInfo, src []byte) error {
    if src == nil {
        *dst = Circle{Status: Null}
        return nil
    }

    if len(src) < 9 {
        return errors.Errorf("invalid length for Circle: %v", len(src))
    }

    str := string(src[2:])
    end := strings.IndexByte(str, ',')
    x, err := strconv.ParseFloat(str[:end], 64)
    if err != nil {
        return err
    }

    str = str[end+1:]
    end = strings.IndexByte(str, ')')

    y, err := strconv.ParseFloat(str[:end], 64)
    if err != nil {
        return err
    }

    str = str[end+2 : len(str)-1]

    r, err := strconv.ParseFloat(str, 64)
    if err != nil {
        return err
    }

    *dst = Circle{P: Vec2{x, y}, R: r, Status: Present}
    return nil
}

func (dst *Circle) DecodeBinary(ci *ConnInfo, src []byte) error {
    if src == nil {
        *dst = Circle{Status: Null}
        return nil
    }

    if len(src) != 24 {
        return errors.Errorf("invalid length for Circle: %v", len(src))
    }

    x := binary.BigEndian.Uint64(src)
    y := binary.BigEndian.Uint64(src[8:])
    r := binary.BigEndian.Uint64(src[16:])

    *dst = Circle{
        P:      Vec2{math.Float64frombits(x), math.Float64frombits(y)},
        R:      math.Float64frombits(r),
        Status: Present,
    }
    return nil
}

func (src *Circle) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
    switch src.Status {
    case Null:
        return nil, nil
    case Undefined:
        return nil, errUndefined
    }

    buf = append(buf, fmt.Sprintf(`<(%s,%s),%s>`,
        strconv.FormatFloat(src.P.X, 'f', -1, 64),
        strconv.FormatFloat(src.P.Y, 'f', -1, 64),
        strconv.FormatFloat(src.R, 'f', -1, 64),
    )...)

    return buf, nil
}

func (src *Circle) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
    switch src.Status {
    case Null:
        return nil, nil
    case Undefined:
        return nil, errUndefined
    }

    buf = pgio.AppendUint64(buf, math.Float64bits(src.P.X))
    buf = pgio.AppendUint64(buf, math.Float64bits(src.P.Y))
    buf = pgio.AppendUint64(buf, math.Float64bits(src.R))
    return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Circle) Scan(src interface{}) error {
    if src == nil {
        *dst = Circle{Status: Null}
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
func (src *Circle) Value() (driver.Value, error) {
    return EncodeValueText(src)
}
