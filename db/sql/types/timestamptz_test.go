package types_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/Ready-Stock/Noah/db/sql/types"
	"github.com/Ready-Stock/Noah/db/sql/types/testutil"
)

func TestTimestamptzTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscodeEqFunc(t, "timestamptz", []interface{}{
		&types.Timestamptz{Time: time.Date(1800, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Present},
		&types.Timestamptz{Time: time.Date(1900, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Present},
		&types.Timestamptz{Time: time.Date(1905, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Present},
		&types.Timestamptz{Time: time.Date(1940, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Present},
		&types.Timestamptz{Time: time.Date(1960, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Present},
		&types.Timestamptz{Time: time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Present},
		&types.Timestamptz{Time: time.Date(1999, 12, 31, 0, 0, 0, 0, time.Local), Status: types.Present},
		&types.Timestamptz{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Present},
		&types.Timestamptz{Time: time.Date(2000, 1, 2, 0, 0, 0, 0, time.Local), Status: types.Present},
		&types.Timestamptz{Time: time.Date(2200, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Present},
		&types.Timestamptz{Status: types.Null},
		&types.Timestamptz{Status: types.Present, InfinityModifier: types.Infinity},
		&types.Timestamptz{Status: types.Present, InfinityModifier: -types.Infinity},
	}, func(a, b interface{}) bool {
		at := a.(types.Timestamptz)
		bt := b.(types.Timestamptz)

		return at.Time.Equal(bt.Time) && at.Status == bt.Status && at.InfinityModifier == bt.InfinityModifier
	})
}

func TestTimestamptzSet(t *testing.T) {
	type _time time.Time

	successfulTests := []struct {
		source interface{}
		result types.Timestamptz
	}{
		{source: time.Date(1900, 1, 1, 0, 0, 0, 0, time.Local), result: types.Timestamptz{Time: time.Date(1900, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Present}},
		{source: time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local), result: types.Timestamptz{Time: time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Present}},
		{source: time.Date(1999, 12, 31, 12, 59, 59, 0, time.Local), result: types.Timestamptz{Time: time.Date(1999, 12, 31, 12, 59, 59, 0, time.Local), Status: types.Present}},
		{source: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), result: types.Timestamptz{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Present}},
		{source: time.Date(2000, 1, 1, 0, 0, 1, 0, time.Local), result: types.Timestamptz{Time: time.Date(2000, 1, 1, 0, 0, 1, 0, time.Local), Status: types.Present}},
		{source: time.Date(2200, 1, 1, 0, 0, 0, 0, time.Local), result: types.Timestamptz{Time: time.Date(2200, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Present}},
		{source: _time(time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)), result: types.Timestamptz{Time: time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Present}},
	}

	for i, tt := range successfulTests {
		var r types.Timestamptz
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if r != tt.result {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestTimestamptzAssignTo(t *testing.T) {
	var tim time.Time
	var ptim *time.Time

	simpleTests := []struct {
		src      types.Timestamptz
		dst      interface{}
		expected interface{}
	}{
		{src: types.Timestamptz{Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Present}, dst: &tim, expected: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local)},
		{src: types.Timestamptz{Time: time.Time{}, Status: types.Null}, dst: &ptim, expected: ((*time.Time)(nil))},
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
		src      types.Timestamptz
		dst      interface{}
		expected interface{}
	}{
		{src: types.Timestamptz{Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Present}, dst: &ptim, expected: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local)},
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
		src types.Timestamptz
		dst interface{}
	}{
		{src: types.Timestamptz{Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local), InfinityModifier: types.Infinity, Status: types.Present}, dst: &tim},
		{src: types.Timestamptz{Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local), InfinityModifier: types.NegativeInfinity, Status: types.Present}, dst: &tim},
		{src: types.Timestamptz{Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Null}, dst: &tim},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}
}
