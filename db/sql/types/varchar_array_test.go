package types_test

import (
	"reflect"
	"testing"

	"github.com/Ready-Stock/Noah/db/sql/types"
	"github.com/Ready-Stock/Noah/db/sql/types/testutil"
)

func TestVarcharArrayTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "varchar[]", []interface{}{
		&types.VarcharArray{
			Elements:   nil,
			Dimensions: nil,
			Status:     types.Present,
		},
		&types.VarcharArray{
			Elements: []types.Varchar{
				{String: "foo", Status: types.Present},
				{Status: types.Null},
			},
			Dimensions: []types.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.VarcharArray{Status: types.Null},
		&types.VarcharArray{
			Elements: []types.Varchar{
				{String: "bar ", Status: types.Present},
				{String: "NuLL", Status: types.Present},
				{String: `wow"quz\`, Status: types.Present},
				{String: "", Status: types.Present},
				{Status: types.Null},
				{String: "null", Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{{Length: 3, LowerBound: 1}, {Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.VarcharArray{
			Elements: []types.Varchar{
				{String: "bar", Status: types.Present},
				{String: "baz", Status: types.Present},
				{String: "quz", Status: types.Present},
				{String: "foo", Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{
				{Length: 2, LowerBound: 4},
				{Length: 2, LowerBound: 2},
			},
			Status: types.Present,
		},
	})
}

func TestVarcharArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.VarcharArray
	}{
		{
			source: []string{"foo"},
			result: types.VarcharArray{
				Elements:   []types.Varchar{{String: "foo", Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present},
		},
		{
			source: (([]string)(nil)),
			result: types.VarcharArray{Status: types.Null},
		},
	}

	for i, tt := range successfulTests {
		var r types.VarcharArray
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestVarcharArrayAssignTo(t *testing.T) {
	var stringSlice []string
	type _stringSlice []string
	var namedStringSlice _stringSlice

	simpleTests := []struct {
		src      types.VarcharArray
		dst      interface{}
		expected interface{}
	}{
		{
			src: types.VarcharArray{
				Elements:   []types.Varchar{{String: "foo", Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &stringSlice,
			expected: []string{"foo"},
		},
		{
			src: types.VarcharArray{
				Elements:   []types.Varchar{{String: "bar", Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &namedStringSlice,
			expected: _stringSlice{"bar"},
		},
		{
			src:      types.VarcharArray{Status: types.Null},
			dst:      &stringSlice,
			expected: (([]string)(nil)),
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
		src types.VarcharArray
		dst interface{}
	}{
		{
			src: types.VarcharArray{
				Elements:   []types.Varchar{{Status: types.Null}},
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
