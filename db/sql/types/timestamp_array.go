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

package types

import (
    "database/sql/driver"
    "encoding/binary"
    "time"

    "github.com/pkg/errors"
    "github.com/readystock/noah/db/sql/pgio"
)

type TimestampArray struct {
    Elements   []Timestamp
    Dimensions []ArrayDimension
    Status     Status
}

func (dst *TimestampArray) Set(src interface{}) error {
    // untyped nil and typed nil interfaces are different
    if src == nil {
        *dst = TimestampArray{Status: Null}
        return nil
    }

    switch value := src.(type) {

    case []time.Time:
        if value == nil {
            *dst = TimestampArray{Status: Null}
        } else if len(value) == 0 {
            *dst = TimestampArray{Status: Present}
        } else {
            elements := make([]Timestamp, len(value))
            for i := range value {
                if err := elements[i].Set(value[i]); err != nil {
                    return err
                }
            }
            *dst = TimestampArray{
                Elements:   elements,
                Dimensions: []ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
                Status:     Present,
            }
        }

    default:
        if originalSrc, ok := underlyingSliceType(src); ok {
            return dst.Set(originalSrc)
        }
        return errors.Errorf("cannot convert %v to TimestampArray", value)
    }

    return nil
}

func (dst *TimestampArray) Get() interface{} {
    switch dst.Status {
    case Present:
        return dst
    case Null:
        return nil
    default:
        return dst.Status
    }
}

func (src *TimestampArray) AssignTo(dst interface{}) error {
    switch src.Status {
    case Present:
        switch v := dst.(type) {

        case *[]time.Time:
            *v = make([]time.Time, len(src.Elements))
            for i := range src.Elements {
                if err := src.Elements[i].AssignTo(&((*v)[i])); err != nil {
                    return err
                }
            }
            return nil

        default:
            if nextDst, retry := GetAssignToDstType(dst); retry {
                return src.AssignTo(nextDst)
            }
        }
    case Null:
        return NullAssignTo(dst)
    }

    return errors.Errorf("cannot decode %#v into %T", src, dst)
}

func (dst *TimestampArray) DecodeText(ci *ConnInfo, src []byte) error {
    if src == nil {
        *dst = TimestampArray{Status: Null}
        return nil
    }

    uta, err := ParseUntypedTextArray(string(src))
    if err != nil {
        return err
    }

    var elements []Timestamp

    if len(uta.Elements) > 0 {
        elements = make([]Timestamp, len(uta.Elements))

        for i, s := range uta.Elements {
            var elem Timestamp
            var elemSrc []byte
            if s != "NULL" {
                elemSrc = []byte(s)
            }
            err = elem.DecodeText(ci, elemSrc)
            if err != nil {
                return err
            }

            elements[i] = elem
        }
    }

    *dst = TimestampArray{Elements: elements, Dimensions: uta.Dimensions, Status: Present}

    return nil
}

func (dst *TimestampArray) DecodeBinary(ci *ConnInfo, src []byte) error {
    if src == nil {
        *dst = TimestampArray{Status: Null}
        return nil
    }

    var arrayHeader ArrayHeader
    rp, err := arrayHeader.DecodeBinary(ci, src)
    if err != nil {
        return err
    }

    if len(arrayHeader.Dimensions) == 0 {
        *dst = TimestampArray{Dimensions: arrayHeader.Dimensions, Status: Present}
        return nil
    }

    elementCount := arrayHeader.Dimensions[0].Length
    for _, d := range arrayHeader.Dimensions[1:] {
        elementCount *= d.Length
    }

    elements := make([]Timestamp, elementCount)

    for i := range elements {
        elemLen := int(int32(binary.BigEndian.Uint32(src[rp:])))
        rp += 4
        var elemSrc []byte
        if elemLen >= 0 {
            elemSrc = src[rp : rp+elemLen]
            rp += elemLen
        }
        err = elements[i].DecodeBinary(ci, elemSrc)
        if err != nil {
            return err
        }
    }

    *dst = TimestampArray{Elements: elements, Dimensions: arrayHeader.Dimensions, Status: Present}
    return nil
}

func (src *TimestampArray) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
    switch src.Status {
    case Null:
        return nil, nil
    case Undefined:
        return nil, errUndefined
    }

    if len(src.Dimensions) == 0 {
        return append(buf, '{', '}'), nil
    }

    buf = EncodeTextArrayDimensions(buf, src.Dimensions)

    // dimElemCounts is the multiples of elements that each array lies on. For
    // example, a single dimension array of length 4 would have a dimElemCounts of
    // [4]. A multi-dimensional array of lengths [3,5,2] would have a
    // dimElemCounts of [30,10,2]. This is used to simplify when to render a '{'
    // or '}'.
    dimElemCounts := make([]int, len(src.Dimensions))
    dimElemCounts[len(src.Dimensions)-1] = int(src.Dimensions[len(src.Dimensions)-1].Length)
    for i := len(src.Dimensions) - 2; i > -1; i-- {
        dimElemCounts[i] = int(src.Dimensions[i].Length) * dimElemCounts[i+1]
    }

    inElemBuf := make([]byte, 0, 32)
    for i, elem := range src.Elements {
        if i > 0 {
            buf = append(buf, ',')
        }

        for _, dec := range dimElemCounts {
            if i%dec == 0 {
                buf = append(buf, '{')
            }
        }

        elemBuf, err := elem.EncodeText(ci, inElemBuf)
        if err != nil {
            return nil, err
        }
        if elemBuf == nil {
            buf = append(buf, `NULL`...)
        } else {
            buf = append(buf, QuoteArrayElementIfNeeded(string(elemBuf))...)
        }

        for _, dec := range dimElemCounts {
            if (i+1)%dec == 0 {
                buf = append(buf, '}')
            }
        }
    }

    return buf, nil
}

func (src *TimestampArray) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
    switch src.Status {
    case Null:
        return nil, nil
    case Undefined:
        return nil, errUndefined
    }

    arrayHeader := ArrayHeader{
        Dimensions: src.Dimensions,
    }

    if dt, ok := ci.DataTypeForName("timestamp"); ok {
        arrayHeader.ElementOID = int32(dt.OID)
    } else {
        return nil, errors.Errorf("unable to find oid for type name %v", "timestamp")
    }

    for i := range src.Elements {
        if src.Elements[i].Status == Null {
            arrayHeader.ContainsNull = true
            break
        }
    }

    buf = arrayHeader.EncodeBinary(ci, buf)

    for i := range src.Elements {
        sp := len(buf)
        buf = pgio.AppendInt32(buf, -1)

        elemBuf, err := src.Elements[i].EncodeBinary(ci, buf)
        if err != nil {
            return nil, err
        }
        if elemBuf != nil {
            buf = elemBuf
            pgio.SetInt32(buf[sp:], int32(len(buf[sp:])-4))
        }
    }

    return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *TimestampArray) Scan(src interface{}) error {
    if src == nil {
        return dst.DecodeText(nil, nil)
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
func (src *TimestampArray) Value() (driver.Value, error) {
    buf, err := src.EncodeText(nil, nil)
    if err != nil {
        return nil, err
    }
    if buf == nil {
        return nil, nil
    }

    return string(buf), nil
}
