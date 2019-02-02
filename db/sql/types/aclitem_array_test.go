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

func TestACLItemArrayTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "aclitem[]", []interface{}{
		&types.ACLItemArray{
			Elements:   nil,
			Dimensions: nil,
			Status:     types.Present,
		},
		&types.ACLItemArray{
			Elements: []types.ACLItem{
				{String: "=r/postgres", Status: types.Present},
				{Status: types.Null},
			},
			Dimensions: []types.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.ACLItemArray{Status: types.Null},
		&types.ACLItemArray{
			Elements: []types.ACLItem{
				{String: "=r/postgres", Status: types.Present},
				{String: "postgres=arwdDxt/postgres", Status: types.Present},
				{String: `postgres=arwdDxt/" tricky, ' } "" \ test user "`, Status: types.Present},
				{String: "=r/postgres", Status: types.Present},
				{Status: types.Null},
				{String: "=r/postgres", Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{{Length: 3, LowerBound: 1}, {Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.ACLItemArray{
			Elements: []types.ACLItem{
				{String: "=r/postgres", Status: types.Present},
				{String: "postgres=arwdDxt/postgres", Status: types.Present},
				{String: "=r/postgres", Status: types.Present},
				{String: "postgres=arwdDxt/postgres", Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{
				{Length: 2, LowerBound: 4},
				{Length: 2, LowerBound: 2},
			},
			Status: types.Present,
		},
	})
}

func TestACLItemArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.ACLItemArray
	}{
		{
			source: []string{"=r/postgres"},
			result: types.ACLItemArray{
				Elements:   []types.ACLItem{{String: "=r/postgres", Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present},
		},
		{
			source: ([]string)(nil),
			result: types.ACLItemArray{Status: types.Null},
		},
	}

	for i, tt := range successfulTests {
		var r types.ACLItemArray
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestACLItemArrayAssignTo(t *testing.T) {
	var stringSlice []string
	type _stringSlice []string
	var namedStringSlice _stringSlice

	simpleTests := []struct {
		src      types.ACLItemArray
		dst      interface{}
		expected interface{}
	}{
		{
			src: types.ACLItemArray{
				Elements:   []types.ACLItem{{String: "=r/postgres", Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &stringSlice,
			expected: []string{"=r/postgres"},
		},
		{
			src: types.ACLItemArray{
				Elements:   []types.ACLItem{{String: "=r/postgres", Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &namedStringSlice,
			expected: _stringSlice{"=r/postgres"},
		},
		{
			src:      types.ACLItemArray{Status: types.Null},
			dst:      &stringSlice,
			expected: ([]string)(nil),
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
		src types.ACLItemArray
		dst interface{}
	}{
		{
			src: types.ACLItemArray{
				Elements:   []types.ACLItem{{Status: types.Null}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst: &stringSlice,
		},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}
}
