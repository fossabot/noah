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
    "math"
    "strconv"
    "strings"

    "github.com/pkg/errors"
    "github.com/readystock/noah/db/sql/pgio"
)

type Path struct {
    P      []Vec2
    Closed bool
    Status Status
}

func (dst *Path) Set(src interface{}) error {
    return errors.Errorf("cannot convert %v to Path", src)
}

func (dst *Path) Get() interface{} {
    switch dst.Status {
    case Present:
        return dst
    case Null:
        return nil
    default:
        return dst.Status
    }
}

func (src *Path) AssignTo(dst interface{}) error {
    return errors.Errorf("cannot assign %v to %T", src, dst)
}

func (dst *Path) DecodeText(ci *ConnInfo, src []byte) error {
    if src == nil {
        *dst = Path{Status: Null}
        return nil
    }

    if len(src) < 7 {
        return errors.Errorf("invalid length for Path: %v", len(src))
    }

    closed := src[0] == '('
    points := make([]Vec2, 0)

    str := string(src[2:])

    for {
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

        points = append(points, Vec2{x, y})

        if end+3 < len(str) {
            str = str[end+3:]
        } else {
            break
        }
    }

    *dst = Path{P: points, Closed: closed, Status: Present}
    return nil
}

func (dst *Path) DecodeBinary(ci *ConnInfo, src []byte) error {
    if src == nil {
        *dst = Path{Status: Null}
        return nil
    }

    if len(src) < 5 {
        return errors.Errorf("invalid length for Path: %v", len(src))
    }

    closed := src[0] == 1
    pointCount := int(binary.BigEndian.Uint32(src[1:]))

    rp := 5

    if 5+pointCount*16 != len(src) {
        return errors.Errorf("invalid length for Path with %d points: %v", pointCount, len(src))
    }

    points := make([]Vec2, pointCount)
    for i := 0; i < len(points); i++ {
        x := binary.BigEndian.Uint64(src[rp:])
        rp += 8
        y := binary.BigEndian.Uint64(src[rp:])
        rp += 8
        points[i] = Vec2{math.Float64frombits(x), math.Float64frombits(y)}
    }

    *dst = Path{
        P:      points,
        Closed: closed,
        Status: Present,
    }
    return nil
}

func (src *Path) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
    switch src.Status {
    case Null:
        return nil, nil
    case Undefined:
        return nil, errUndefined
    }

    var startByte, endByte byte
    if src.Closed {
        startByte = '('
        endByte = ')'
    } else {
        startByte = '['
        endByte = ']'
    }
    buf = append(buf, startByte)

    for i, p := range src.P {
        if i > 0 {
            buf = append(buf, ',')
        }
        buf = append(buf, fmt.Sprintf(`(%s,%s)`,
            strconv.FormatFloat(p.X, 'f', -1, 64),
            strconv.FormatFloat(p.Y, 'f', -1, 64),
        )...)
    }

    return append(buf, endByte), nil
}

func (src *Path) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
    switch src.Status {
    case Null:
        return nil, nil
    case Undefined:
        return nil, errUndefined
    }

    var closeByte byte
    if src.Closed {
        closeByte = 1
    }
    buf = append(buf, closeByte)

    buf = pgio.AppendInt32(buf, int32(len(src.P)))

    for _, p := range src.P {
        buf = pgio.AppendUint64(buf, math.Float64bits(p.X))
        buf = pgio.AppendUint64(buf, math.Float64bits(p.Y))
    }

    return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Path) Scan(src interface{}) error {
    if src == nil {
        *dst = Path{Status: Null}
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
func (src *Path) Value() (driver.Value, error) {
    return EncodeValueText(src)
}
