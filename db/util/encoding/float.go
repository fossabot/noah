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

package encoding

import (
    "math"

    "github.com/pkg/errors"
)

// EncodeFloatAscending returns the resulting byte slice with the encoded float64
// appended to b. The encoded format for a float64 value f is, for positive f, the
// encoding of the 64 bits (in IEEE 754 format) re-interpreted as an int64 and
// encoded using EncodeUint64Ascending. For negative f, we keep the sign bit and
// invert all other bits, encoding this value using EncodeUint64Descending. This
// approach was inspired by in github.com/google/orderedcode/orderedcode.go.
//
// One of five single-byte prefix tags are appended to the front of the encoding.
// These tags enforce logical ordering of keys for both ascending and descending
// encoding directions. The tags split the encoded floats into five categories:
// - NaN for an ascending encoding direction
// - Negative valued floats
// - Zero (positive and negative)
// - Positive valued floats
// - NaN for a descending encoding direction
// This ordering ensures that NaNs are always sorted first in either encoding
// direction, and that after them a logical ordering is followed.
func EncodeFloatAscending(b []byte, f float64) []byte {
    // Handle the simplistic cases first.
    switch {
    case math.IsNaN(f):
        return append(b, floatNaN)
    case f == 0:
        // This encodes both positive and negative zero the same. Negative zero uses
        // composite indexes to decode itself correctly.
        return append(b, floatZero)
    }
    u := math.Float64bits(f)
    if u&(1<<63) != 0 {
        u = ^u
        b = append(b, floatNeg)
    } else {
        b = append(b, floatPos)
    }
    return EncodeUint64Ascending(b, u)
}

// EncodeFloatDescending is the descending version of EncodeFloatAscending.
func EncodeFloatDescending(b []byte, f float64) []byte {
    if math.IsNaN(f) {
        return append(b, floatNaNDesc)
    }
    return EncodeFloatAscending(b, -f)
}

// DecodeFloatAscending returns the remaining byte slice after decoding and the decoded
// float64 from buf.
func DecodeFloatAscending(buf []byte) ([]byte, float64, error) {
    if PeekType(buf) != Float {
        return buf, 0, errors.Errorf("did not find marker")
    }
    switch buf[0] {
    case floatNaN, floatNaNDesc:
        return buf[1:], math.NaN(), nil
    case floatNeg:
        b, u, err := DecodeUint64Ascending(buf[1:])
        if err != nil {
            return b, 0, err
        }
        u = ^u
        return b, math.Float64frombits(u), nil
    case floatZero:
        return buf[1:], 0, nil
    case floatPos:
        b, u, err := DecodeUint64Ascending(buf[1:])
        if err != nil {
            return b, 0, err
        }
        return b, math.Float64frombits(u), nil
    default:
        return nil, 0, errors.Errorf("unknown prefix of the encoded byte slice: %q", buf)
    }
}

// DecodeFloatDescending decodes floats encoded with EncodeFloatDescending.
func DecodeFloatDescending(buf []byte) ([]byte, float64, error) {
    b, r, err := DecodeFloatAscending(buf)
    return b, -r, err
}
