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
 */

package types_test

import (
	"reflect"
	"testing"

    "github.com/Ready-Stock/noah/db/sql/types"
    "github.com/Ready-Stock/noah/db/sql/types/testutil"
)

func TestFloat8Transcode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "float8", []interface{}{
		&types.Float8{Float: -1, Status: types.Present},
		&types.Float8{Float: 0, Status: types.Present},
		&types.Float8{Float: 0.00001, Status: types.Present},
		&types.Float8{Float: 1, Status: types.Present},
		&types.Float8{Float: 9999.99, Status: types.Present},
		&types.Float8{Float: 0, Status: types.Null},
	})
}

func TestFloat8Set(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.Float8
	}{
		{source: float32(1), result: types.Float8{Float: 1, Status: types.Present}},
		{source: float64(1), result: types.Float8{Float: 1, Status: types.Present}},
		{source: int8(1), result: types.Float8{Float: 1, Status: types.Present}},
		{source: int16(1), result: types.Float8{Float: 1, Status: types.Present}},
		{source: int32(1), result: types.Float8{Float: 1, Status: types.Present}},
		{source: int64(1), result: types.Float8{Float: 1, Status: types.Present}},
		{source: int8(-1), result: types.Float8{Float: -1, Status: types.Present}},
		{source: int16(-1), result: types.Float8{Float: -1, Status: types.Present}},
		{source: int32(-1), result: types.Float8{Float: -1, Status: types.Present}},
		{source: int64(-1), result: types.Float8{Float: -1, Status: types.Present}},
		{source: uint8(1), result: types.Float8{Float: 1, Status: types.Present}},
		{source: uint16(1), result: types.Float8{Float: 1, Status: types.Present}},
		{source: uint32(1), result: types.Float8{Float: 1, Status: types.Present}},
		{source: uint64(1), result: types.Float8{Float: 1, Status: types.Present}},
		{source: "1", result: types.Float8{Float: 1, Status: types.Present}},
		{source: _int8(1), result: types.Float8{Float: 1, Status: types.Present}},
	}

	for i, tt := range successfulTests {
		var r types.Float8
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if r != tt.result {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestFloat8AssignTo(t *testing.T) {
	var i8 int8
	var i16 int16
	var i32 int32
	var i64 int64
	var i int
	var ui8 uint8
	var ui16 uint16
	var ui32 uint32
	var ui64 uint64
	var ui uint
	var pi8 *int8
	var _i8 _int8
	var _pi8 *_int8
	var f32 float32
	var f64 float64
	var pf32 *float32
	var pf64 *float64

	simpleTests := []struct {
		src      types.Float8
		dst      interface{}
		expected interface{}
	}{
		{src: types.Float8{Float: 42, Status: types.Present}, dst: &f32, expected: float32(42)},
		{src: types.Float8{Float: 42, Status: types.Present}, dst: &f64, expected: float64(42)},
		{src: types.Float8{Float: 42, Status: types.Present}, dst: &i16, expected: int16(42)},
		{src: types.Float8{Float: 42, Status: types.Present}, dst: &i32, expected: int32(42)},
		{src: types.Float8{Float: 42, Status: types.Present}, dst: &i64, expected: int64(42)},
		{src: types.Float8{Float: 42, Status: types.Present}, dst: &i, expected: int(42)},
		{src: types.Float8{Float: 42, Status: types.Present}, dst: &ui8, expected: uint8(42)},
		{src: types.Float8{Float: 42, Status: types.Present}, dst: &ui16, expected: uint16(42)},
		{src: types.Float8{Float: 42, Status: types.Present}, dst: &ui32, expected: uint32(42)},
		{src: types.Float8{Float: 42, Status: types.Present}, dst: &ui64, expected: uint64(42)},
		{src: types.Float8{Float: 42, Status: types.Present}, dst: &ui, expected: uint(42)},
		{src: types.Float8{Float: 42, Status: types.Present}, dst: &_i8, expected: _int8(42)},
        {src: types.Float8{Float: 0, Status: types.Null}, dst: &pi8, expected: (*int8)(nil)},
        {src: types.Float8{Float: 0, Status: types.Null}, dst: &_pi8, expected: (*_int8)(nil)},
	}

	for i, tt := range simpleTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Interface(); dst != tt.expected {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}

	pointerAllocTests := []struct {
		src      types.Float8
		dst      interface{}
		expected interface{}
	}{
		{src: types.Float8{Float: 42, Status: types.Present}, dst: &pf32, expected: float32(42)},
		{src: types.Float8{Float: 42, Status: types.Present}, dst: &pf64, expected: float64(42)},
	}

	for i, tt := range pointerAllocTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Elem().Interface(); dst != tt.expected {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}

	errorTests := []struct {
		src types.Float8
		dst interface{}
	}{
		{src: types.Float8{Float: 150, Status: types.Present}, dst: &i8},
		{src: types.Float8{Float: 40000, Status: types.Present}, dst: &i16},
		{src: types.Float8{Float: -1, Status: types.Present}, dst: &ui8},
		{src: types.Float8{Float: -1, Status: types.Present}, dst: &ui16},
		{src: types.Float8{Float: -1, Status: types.Present}, dst: &ui32},
		{src: types.Float8{Float: -1, Status: types.Present}, dst: &ui64},
		{src: types.Float8{Float: -1, Status: types.Present}, dst: &ui},
		{src: types.Float8{Float: 0, Status: types.Null}, dst: &i32},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}
}
