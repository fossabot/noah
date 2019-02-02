/*
 * Copyright (c) 2019 Ready Stock
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
 */

package types_test

import (
	"reflect"
	"testing"

	"github.com/readystock/noah/db/sql/types"
	"github.com/readystock/noah/db/sql/types/testutil"
)

func TestBoolArrayTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "bool[]", []interface{}{
		&types.BoolArray{
			Elements:   nil,
			Dimensions: nil,
			Status:     types.Present,
		},
		&types.BoolArray{
			Elements: []types.Bool{
				{Bool: true, Status: types.Present},
				{Status: types.Null},
			},
			Dimensions: []types.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.BoolArray{Status: types.Null},
		&types.BoolArray{
			Elements: []types.Bool{
				{Bool: true, Status: types.Present},
				{Bool: true, Status: types.Present},
				{Bool: false, Status: types.Present},
				{Bool: true, Status: types.Present},
				{Status: types.Null},
				{Bool: false, Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{{Length: 3, LowerBound: 1}, {Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.BoolArray{
			Elements: []types.Bool{
				{Bool: true, Status: types.Present},
				{Bool: false, Status: types.Present},
				{Bool: true, Status: types.Present},
				{Bool: false, Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{
				{Length: 2, LowerBound: 4},
				{Length: 2, LowerBound: 2},
			},
			Status: types.Present,
		},
	})
}

func TestBoolArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.BoolArray
	}{
		{
			source: []bool{true},
			result: types.BoolArray{
				Elements:   []types.Bool{{Bool: true, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present},
		},
		{
			source: ([]bool)(nil),
			result: types.BoolArray{Status: types.Null},
		},
	}

	for i, tt := range successfulTests {
		var r types.BoolArray
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestBoolArrayAssignTo(t *testing.T) {
	var boolSlice []bool
	type _boolSlice []bool
	var namedBoolSlice _boolSlice

	simpleTests := []struct {
		src      types.BoolArray
		dst      interface{}
		expected interface{}
	}{
		{
			src: types.BoolArray{
				Elements:   []types.Bool{{Bool: true, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &boolSlice,
			expected: []bool{true},
		},
		{
			src: types.BoolArray{
				Elements:   []types.Bool{{Bool: true, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &namedBoolSlice,
			expected: _boolSlice{true},
		},
		{
			src:      types.BoolArray{Status: types.Null},
			dst:      &boolSlice,
			expected: ([]bool)(nil),
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
		src types.BoolArray
		dst interface{}
	}{
		{
			src: types.BoolArray{
				Elements:   []types.Bool{{Status: types.Null}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst: &boolSlice,
		},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}

}
