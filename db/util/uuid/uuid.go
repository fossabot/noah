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

package uuid

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/readystock/noah/db/util/uint128"
	"github.com/satori/go.uuid"
)

// UUID is a thin wrapper around "github.com/satori/go.uuid".UUID that can be
// used as a gogo/protobuf customtype.
type UUID struct {
	uuid.UUID
}

// Nil is the empty UUID with all 128 bits set to zero.
var Nil = UUID{uuid.Nil}

// Short returns the first eight characters of the output of String().
func (u UUID) Short() string {
	return u.String()[:8]
}

// ShortStringer implements fmt.Stringer to output Short() on String().
type ShortStringer UUID

// String is part of fmt.Stringer.
func (s ShortStringer) String() string {
	return UUID(s).Short()
}

var _ fmt.Stringer = ShortStringer{}

// Bytes shadows (*github.com/satori/go.uuid.UUID).Bytes() to prevent UUID
// from implementing github.com/golang/protobuf/proto.raw, the semantics of
// which do not match the semantics of the shadowed method. See
// https://github.com/golang/protobuf/blob/5386fff/proto/text.go#L173:L176.
//
// TODO(tamird): remove when fixed upstream. See
// https://github.com/gogo/protobuf/pull/227 and
// https://github.com/golang/protobuf/issues/311.
func (UUID) Bytes() {
	panic("intentionally shadowed; use GetBytes()")
}

// Silence unused warning for UUID.Bytes.
var _ = UUID.Bytes

// Equal returns true iff the receiver equals the argument.
//
// This method exists only to conform to the API expected by gogoproto's
// generated Equal implementations.
func (u UUID) Equal(t UUID) bool {
	return u == t
}

// GetBytes returns the UUID as a byte slice.
func (u UUID) GetBytes() []byte {
	return u.UUID.Bytes()
}

// ToUint128 returns the UUID as a Uint128.
func (u UUID) ToUint128() uint128.Uint128 {
	return uint128.FromBytes(u.GetBytes())
}

// Size returns the marshaled size of u, in bytes.
func (u UUID) Size() int {
	return len(u.UUID)
}

// MarshalTo marshals u to data.
func (u UUID) MarshalTo(data []byte) (int, error) {
	return copy(data, u.UUID.Bytes()), nil
}

// Unmarshal unmarshals data to u.
func (u *UUID) Unmarshal(data []byte) error {
	return u.UUID.UnmarshalBinary(data)
}

// MarshalJSON returns the JSON encoding of u.
func (u UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

// UnmarshalJSON unmarshals the JSON encoded data into u.
func (u *UUID) UnmarshalJSON(data []byte) error {
	var uuidString string
	if err := json.Unmarshal(data, &uuidString); err != nil {
		return err
	}
	uuid, err := FromString(uuidString)
	*u = uuid
	return err
}

// MakeV4 delegates to "github.com/satori/go.uuid".NewV4 and wraps the result in
// a UUID.
func MakeV4() UUID {
	uuid, _ := uuid.NewV4()
	return UUID{uuid}
}

// NewPopulatedUUID returns a populated UUID.
func NewPopulatedUUID(r interface {
	Int63() int64
}) *UUID {
	var u uuid.UUID
	binary.LittleEndian.PutUint64(u[:8], uint64(r.Int63()))
	binary.LittleEndian.PutUint64(u[8:], uint64(r.Int63()))
	return &UUID{u}
}

// FromBytes delegates to "github.com/satori/go.uuid".FromBytes and wraps the
// result in a UUID.
func FromBytes(input []byte) (UUID, error) {
	u, err := uuid.FromBytes(input)
	return UUID{u}, err
}

// FromString delegates to "github.com/satori/go.uuid".FromString and wraps the
// result in a UUID.
func FromString(input string) (UUID, error) {
	u, err := uuid.FromString(input)
	return UUID{u}, err
}

// FromUint128 delegates to "github.com/satori/go.uuid".FromBytes and wraps the
// result in a UUID.
func FromUint128(input uint128.Uint128) UUID {
	u, err := uuid.FromBytes(input.GetBytes())
	if err != nil {
		panic(errors.Wrap(err, "should never happen with 16 byte slice"))
	}
	return UUID{u}
}
