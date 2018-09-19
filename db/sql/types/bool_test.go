package types_test

import (
	"reflect"
	"testing"

	"github.com/Ready-Stock/Noah/db/sql/types"
	"github.com/Ready-Stock/Noah/db/sql/types/testutil"
)

func TestBoolTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "bool", []interface{}{
		&types.Bool{Bool: false, Status: types.Present},
		&types.Bool{Bool: true, Status: types.Present},
		&types.Bool{Bool: false, Status: types.Null},
	})
}

func TestBoolSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.Bool
	}{
		{source: true, result: types.Bool{Bool: true, Status: types.Present}},
		{source: false, result: types.Bool{Bool: false, Status: types.Present}},
		{source: "true", result: types.Bool{Bool: true, Status: types.Present}},
		{source: "false", result: types.Bool{Bool: false, Status: types.Present}},
		{source: "t", result: types.Bool{Bool: true, Status: types.Present}},
		{source: "f", result: types.Bool{Bool: false, Status: types.Present}},
		{source: _bool(true), result: types.Bool{Bool: true, Status: types.Present}},
		{source: _bool(false), result: types.Bool{Bool: false, Status: types.Present}},
		{source: nil, result: types.Bool{Status: types.Null}},
	}

	for i, tt := range successfulTests {
		var r types.Bool
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if r != tt.result {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestBoolAssignTo(t *testing.T) {
	var b bool
	var _b _bool
	var pb *bool
	var _pb *_bool

	simpleTests := []struct {
		src      types.Bool
		dst      interface{}
		expected interface{}
	}{
		{src: types.Bool{Bool: false, Status: types.Present}, dst: &b, expected: false},
		{src: types.Bool{Bool: true, Status: types.Present}, dst: &b, expected: true},
		{src: types.Bool{Bool: false, Status: types.Present}, dst: &_b, expected: _bool(false)},
		{src: types.Bool{Bool: true, Status: types.Present}, dst: &_b, expected: _bool(true)},
		{src: types.Bool{Bool: false, Status: types.Null}, dst: &pb, expected: ((*bool)(nil))},
		{src: types.Bool{Bool: false, Status: types.Null}, dst: &_pb, expected: ((*_bool)(nil))},
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
		src      types.Bool
		dst      interface{}
		expected interface{}
	}{
		{src: types.Bool{Bool: true, Status: types.Present}, dst: &pb, expected: true},
		{src: types.Bool{Bool: true, Status: types.Present}, dst: &_pb, expected: _bool(true)},
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
}
