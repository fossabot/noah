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

func TestByteaArrayTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "bytea[]", []interface{}{
		&types.ByteaArray{
			Elements:   nil,
			Dimensions: nil,
			Status:     types.Present,
		},
		&types.ByteaArray{
			Elements: []types.Bytea{
				{Bytes: []byte{1, 2, 3}, Status: types.Present},
				{Status: types.Null},
			},
			Dimensions: []types.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.ByteaArray{Status: types.Null},
		&types.ByteaArray{
			Elements: []types.Bytea{
				{Bytes: []byte{1, 2, 3}, Status: types.Present},
				{Bytes: []byte{1, 2, 3}, Status: types.Present},
				{Bytes: []byte{}, Status: types.Present},
				{Bytes: []byte{1, 2, 3}, Status: types.Present},
				{Status: types.Null},
				{Bytes: []byte{1}, Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{{Length: 3, LowerBound: 1}, {Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.ByteaArray{
			Elements: []types.Bytea{
				{Bytes: []byte{1, 2, 3}, Status: types.Present},
				{Bytes: []byte{}, Status: types.Present},
				{Bytes: []byte{1, 2, 3}, Status: types.Present},
				{Bytes: []byte{1}, Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{
				{Length: 2, LowerBound: 4},
				{Length: 2, LowerBound: 2},
			},
			Status: types.Present,
		},
	})
}

func TestByteaArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.ByteaArray
	}{
		{
			source: [][]byte{{1, 2, 3}},
			result: types.ByteaArray{
				Elements:   []types.Bytea{{Bytes: []byte{1, 2, 3}, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present},
		},
		{
			source: ([][]byte)(nil),
			result: types.ByteaArray{Status: types.Null},
		},
	}

	for i, tt := range successfulTests {
		var r types.ByteaArray
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestByteaArrayAssignTo(t *testing.T) {
	var byteByteSlice [][]byte

	simpleTests := []struct {
		src      types.ByteaArray
		dst      interface{}
		expected interface{}
	}{
		{
			src: types.ByteaArray{
				Elements:   []types.Bytea{{Bytes: []byte{1, 2, 3}, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &byteByteSlice,
			expected: [][]byte{{1, 2, 3}},
		},
		{
			src:      types.ByteaArray{Status: types.Null},
			dst:      &byteByteSlice,
			expected: ([][]byte)(nil),
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
}
