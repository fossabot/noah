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
    "github.com/pkg/errors"
    "github.com/readystock/noah/db/sql/pgio"
    "time"
)

type Date struct {
    Time             time.Time
    Status           Status
    InfinityModifier InfinityModifier
}

const (
    negativeInfinityDayOffset = -2147483648
    infinityDayOffset         = 2147483647
)

func (dst *Date) Set(src interface{}) error {
    if src == nil {
        *dst = Date{Status: Null}
        return nil
    }

    switch value := src.(type) {
    case time.Time:
        *dst = Date{Time: value, Status: Present}
    default:
        if originalSrc, ok := underlyingTimeType(src); ok {
            return dst.Set(originalSrc)
        }
        return errors.Errorf("cannot convert %v to Date", value)
    }

    return nil
}

func (dst *Date) Get() interface{} {
    switch dst.Status {
    case Present:
        if dst.InfinityModifier != None {
            return dst.InfinityModifier
        }
        return dst.Time
    case Null:
        return nil
    default:
        return dst.Status
    }
}

func (src *Date) AssignTo(dst interface{}) error {
    switch src.Status {
    case Present:
        switch v := dst.(type) {
        case *time.Time:
            if src.InfinityModifier != None {
                return errors.Errorf("cannot assign %v to %T", src, dst)
            }
            *v = src.Time
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

func (dst *Date) DecodeText(ci *ConnInfo, src []byte) error {
    if src == nil {
        *dst = Date{Status: Null}
        return nil
    }

    sbuf := string(src)
    switch sbuf {
    case "infinity":
        *dst = Date{Status: Present, InfinityModifier: Infinity}
    case "-infinity":
        *dst = Date{Status: Present, InfinityModifier: -Infinity}
    default:
        t, err := time.ParseInLocation("2006-01-02", sbuf, time.UTC)
        if err != nil {
            return err
        }

        *dst = Date{Time: t, Status: Present}
    }

    return nil
}

func (dst *Date) DecodeBinary(ci *ConnInfo, src []byte) error {
    if src == nil {
        *dst = Date{Status: Null}
        return nil
    }

    if len(src) != 4 {
        return errors.Errorf("invalid length for date: %v", len(src))
    }

    dayOffset := int32(binary.BigEndian.Uint32(src))

    switch dayOffset {
    case infinityDayOffset:
        *dst = Date{Status: Present, InfinityModifier: Infinity}
    case negativeInfinityDayOffset:
        *dst = Date{Status: Present, InfinityModifier: -Infinity}
    default:
        t := time.Date(2000, 1, int(1+dayOffset), 0, 0, 0, 0, time.UTC)
        *dst = Date{Time: t, Status: Present}
    }

    return nil
}

func (src *Date) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
    switch src.Status {
    case Null:
        return nil, nil
    case Undefined:
        return nil, errUndefined
    }

    var s string

    switch src.InfinityModifier {
    case None:
        s = src.Time.Format("2006-01-02")
    case Infinity:
        s = "infinity"
    case NegativeInfinity:
        s = "-infinity"
    }

    return append(buf, s...), nil
}

func (src *Date) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
    switch src.Status {
    case Null:
        return nil, nil
    case Undefined:
        return nil, errUndefined
    }

    var daysSinceDateEpoch int32
    switch src.InfinityModifier {
    case None:
        tUnix := time.Date(src.Time.Year(), src.Time.Month(), src.Time.Day(), 0, 0, 0, 0, time.UTC).Unix()
        dateEpoch := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Unix()

        secSinceDateEpoch := tUnix - dateEpoch
        daysSinceDateEpoch = int32(secSinceDateEpoch / 86400)
    case Infinity:
        daysSinceDateEpoch = infinityDayOffset
    case NegativeInfinity:
        daysSinceDateEpoch = negativeInfinityDayOffset
    }

    return pgio.AppendInt32(buf, daysSinceDateEpoch), nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Date) Scan(src interface{}) error {
    if src == nil {
        *dst = Date{Status: Null}
        return nil
    }

    switch src := src.(type) {
    case string:
        return dst.DecodeText(nil, []byte(src))
    case []byte:
        srcCopy := make([]byte, len(src))
        copy(srcCopy, src)
        return dst.DecodeText(nil, srcCopy)
    case time.Time:
        *dst = Date{Time: src, Status: Present}
        return nil
    }

    return errors.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src *Date) Value() (driver.Value, error) {
    switch src.Status {
    case Present:
        if src.InfinityModifier != None {
            return src.InfinityModifier.String(), nil
        }
        return src.Time, nil
    case Null:
        return nil, nil
    default:
        return nil, errUndefined
    }
}
