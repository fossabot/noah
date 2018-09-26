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

package uint128

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/pkg/errors"
)

// Uint128 is a big-endian 128 bit unsigned integer which wraps two uint64s.
type Uint128 struct {
	Hi, Lo uint64
}

// GetBytes returns a big-endian byte representation.
func (u Uint128) GetBytes() []byte {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf[:8], u.Hi)
	binary.BigEndian.PutUint64(buf[8:], u.Lo)
	return buf
}

// String returns a hexadecimal string representation.
func (u Uint128) String() string {
	return hex.EncodeToString(u.GetBytes())
}

// Equal returns whether or not the Uint128 are equivalent.
func (u Uint128) Equal(o Uint128) bool {
	return u.Hi == o.Hi && u.Lo == o.Lo
}

// Compare compares the two Uint128.
func (u Uint128) Compare(o Uint128) int {
	if u.Hi > o.Hi {
		return 1
	} else if u.Hi < o.Hi {
		return -1
	} else if u.Lo > o.Lo {
		return 1
	} else if u.Lo < o.Lo {
		return -1
	}
	return 0
}

// Add returns a new Uint128 incremented by n.
func (u Uint128) Add(n uint64) Uint128 {
	lo := u.Lo + n
	hi := u.Hi
	if u.Lo > lo {
		hi++
	}
	return Uint128{hi, lo}
}

// Sub returns a new Uint128 decremented by n.
func (u Uint128) Sub(n uint64) Uint128 {
	lo := u.Lo - n
	hi := u.Hi
	if u.Lo < lo {
		hi--
	}
	return Uint128{hi, lo}
}

// And returns a new Uint128 that is the bitwise AND of two Uint128 values.
func (u Uint128) And(o Uint128) Uint128 {
	return Uint128{u.Hi & o.Hi, u.Lo & o.Lo}
}

// Or returns a new Uint128 that is the bitwise OR of two Uint128 values.
func (u Uint128) Or(o Uint128) Uint128 {
	return Uint128{u.Hi | o.Hi, u.Lo | o.Lo}
}

// Xor returns a new Uint128 that is the bitwise XOR of two Uint128 values.
func (u Uint128) Xor(o Uint128) Uint128 {
	return Uint128{u.Hi ^ o.Hi, u.Lo ^ o.Lo}
}

// FromBytes parses the byte slice as a 128 bit big-endian unsigned integer.
func FromBytes(b []byte) Uint128 {
	hi := binary.BigEndian.Uint64(b[:8])
	lo := binary.BigEndian.Uint64(b[8:])
	return Uint128{hi, lo}
}

// FromString parses a hexadecimal string as a 128-bit big-endian unsigned integer.
func FromString(s string) (Uint128, error) {
	if len(s) > 32 {
		return Uint128{}, errors.Errorf("input string %s too large for uint128", s)
	}
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return Uint128{}, errors.Wrapf(err, "could not decode %s as hex", s)
	}

	// Grow the byte slice if it's smaller than 16 bytes, by prepending 0s
	if len(bytes) < 16 {
		bytesCopy := make([]byte, 16)
		copy(bytesCopy[(16-len(bytes)):], bytes)
		bytes = bytesCopy
	}

	return FromBytes(bytes), nil
}

// FromInts takes in two unsigned 64-bit integers and constructs a Uint128.
func FromInts(hi uint64, lo uint64) Uint128 {
	return Uint128{hi, lo}
}
