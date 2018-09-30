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
 */

package types_test

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/Ready-Stock/Noah/db/sql/types"
	"github.com/Ready-Stock/Noah/db/sql/types/testutil"
)

func TestNumericArrayTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "numeric[]", []interface{}{
		&types.NumericArray{
			Elements:   nil,
			Dimensions: nil,
			Status:     types.Present,
		},
		&types.NumericArray{
			Elements: []types.Numeric{
				{Int: big.NewInt(1), Status: types.Present},
				{Status: types.Null},
			},
			Dimensions: []types.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.NumericArray{Status: types.Null},
		&types.NumericArray{
			Elements: []types.Numeric{
				{Int: big.NewInt(1), Status: types.Present},
				{Int: big.NewInt(2), Status: types.Present},
				{Int: big.NewInt(3), Status: types.Present},
				{Int: big.NewInt(4), Status: types.Present},
				{Status: types.Null},
				{Int: big.NewInt(6), Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{{Length: 3, LowerBound: 1}, {Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.NumericArray{
			Elements: []types.Numeric{
				{Int: big.NewInt(1), Status: types.Present},
				{Int: big.NewInt(2), Status: types.Present},
				{Int: big.NewInt(3), Status: types.Present},
				{Int: big.NewInt(4), Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{
				{Length: 2, LowerBound: 4},
				{Length: 2, LowerBound: 2},
			},
			Status: types.Present,
		},
	})
}

func TestNumericArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.NumericArray
	}{
		{
			source: []float32{1},
			result: types.NumericArray{
				Elements:   []types.Numeric{{Int: big.NewInt(1), Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present},
		},
		{
			source: []float64{1},
			result: types.NumericArray{
				Elements:   []types.Numeric{{Int: big.NewInt(1), Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present},
		},
		{
			source: (([]float32)(nil)),
			result: types.NumericArray{Status: types.Null},
		},
	}

	for i, tt := range successfulTests {
		var r types.NumericArray
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestNumericArrayAssignTo(t *testing.T) {
	var float32Slice []float32
	var float64Slice []float64

	simpleTests := []struct {
		src      types.NumericArray
		dst      interface{}
		expected interface{}
	}{
		{
			src: types.NumericArray{
				Elements:   []types.Numeric{{Int: big.NewInt(1), Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &float32Slice,
			expected: []float32{1},
		},
		{
			src: types.NumericArray{
				Elements:   []types.Numeric{{Int: big.NewInt(1), Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &float64Slice,
			expected: []float64{1},
		},
		{
			src:      types.NumericArray{Status: types.Null},
			dst:      &float32Slice,
			expected: (([]float32)(nil)),
		},
	}

	for i, tt := range simpleTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Interface(); !reflect.DeepEqual(dst, tt.expected) {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}

	errorTests := []struct {
		src types.NumericArray
		dst interface{}
	}{
		{
			src: types.NumericArray{
				Elements:   []types.Numeric{{Status: types.Null}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst: &float32Slice,
		},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}

}
