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
    "encoding/hex"
    "fmt"

    "github.com/pkg/errors"
)

type UUID struct {
    Bytes  [16]byte
    Status Status
}

func (dst *UUID) Set(src interface{}) error {
    if src == nil {
        *dst = UUID{Status: Null}
        return nil
    }

    switch value := src.(type) {
    case [16]byte:
        *dst = UUID{Bytes: value, Status: Present}
    case []byte:
        if value != nil {
            if len(value) != 16 {
                return errors.Errorf("[]byte must be 16 bytes to convert to UUID: %d", len(value))
            }
            *dst = UUID{Status: Present}
            copy(dst.Bytes[:], value)
        } else {
            *dst = UUID{Status: Null}
        }
    case string:
        uuid, err := parseUUID(value)
        if err != nil {
            return err
        }
        *dst = UUID{Bytes: uuid, Status: Present}
    default:
        if originalSrc, ok := underlyingPtrType(src); ok {
            return dst.Set(originalSrc)
        }
        return errors.Errorf("cannot convert %v to UUID", value)
    }

    return nil
}

func (dst *UUID) Get() interface{} {
    switch dst.Status {
    case Present:
        return dst.Bytes
    case Null:
        return nil
    default:
        return dst.Status
    }
}

func (src *UUID) AssignTo(dst interface{}) error {
    switch src.Status {
    case Present:
        switch v := dst.(type) {
        case *[16]byte:
            *v = src.Bytes
            return nil
        case *[]byte:
            *v = make([]byte, 16)
            copy(*v, src.Bytes[:])
            return nil
        case *string:
            *v = EncodeUUID(src.Bytes)
            return nil
        default:
            if nextDst, retry := GetAssignToDstType(v); retry {
                return src.AssignTo(nextDst)
            }
        }
    case Null:
        return NullAssignTo(dst)
    }

    return errors.Errorf("cannot assign %v into %T", src, dst)
}

// parseUUID converts a string UUID in standard form to a byte array.
func parseUUID(src string) (dst [16]byte, err error) {
    src = src[0:8] + src[9:13] + src[14:18] + src[19:23] + src[24:]
    buf, err := hex.DecodeString(src)
    if err != nil {
        return dst, err
    }

    copy(dst[:], buf)
    return dst, err
}

// encodeUUID converts a uuid byte array to UUID standard string form.
func EncodeUUID(src [16]byte) string {
    return fmt.Sprintf("%x-%x-%x-%x-%x", src[0:4], src[4:6], src[6:8], src[8:10], src[10:16])
}

func (dst *UUID) DecodeText(ci *ConnInfo, src []byte) error {
    if src == nil {
        *dst = UUID{Status: Null}
        return nil
    }

    if len(src) != 36 {
        return errors.Errorf("invalid length for UUID: %v", len(src))
    }

    buf, err := parseUUID(string(src))
    if err != nil {
        return err
    }

    *dst = UUID{Bytes: buf, Status: Present}
    return nil
}

func (dst *UUID) DecodeBinary(ci *ConnInfo, src []byte) error {
    if src == nil {
        *dst = UUID{Status: Null}
        return nil
    }

    if len(src) != 16 {
        return errors.Errorf("invalid length for UUID: %v", len(src))
    }

    *dst = UUID{Status: Present}
    copy(dst.Bytes[:], src)
    return nil
}

func (src *UUID) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
    switch src.Status {
    case Null:
        return nil, nil
    case Undefined:
        return nil, errUndefined
    }

    return append(buf, EncodeUUID(src.Bytes)...), nil
}

func (src *UUID) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
    switch src.Status {
    case Null:
        return nil, nil
    case Undefined:
        return nil, errUndefined
    }

    return append(buf, src.Bytes[:]...), nil
}

// Scan implements the database/sql Scanner interface.
func (dst *UUID) Scan(src interface{}) error {
    if src == nil {
        *dst = UUID{Status: Null}
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
func (src *UUID) Value() (driver.Value, error) {
    return EncodeValueText(src)
}
