package types_test

import (
	"testing"

	"github.com/Ready-Stock/Noah/db/sql/types"
	"github.com/Ready-Stock/Noah/db/sql/types/testutil"
)

func TestBPCharArrayTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "char(8)[]", []interface{}{
		&types.BPCharArray{
			Elements:   nil,
			Dimensions: nil,
			Status:     types.Present,
		},
		&types.BPCharArray{
			Elements: []types.BPChar{
				types.BPChar{String: "foo     ", Status: types.Present},
				types.BPChar{Status: types.Null},
			},
			Dimensions: []types.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.BPCharArray{Status: types.Null},
		&types.BPCharArray{
			Elements: []types.BPChar{
				types.BPChar{String: "bar     ", Status: types.Present},
				types.BPChar{String: "NuLL    ", Status: types.Present},
				types.BPChar{String: `wow"quz\`, Status: types.Present},
				types.BPChar{String: "1       ", Status: types.Present},
				types.BPChar{String: "1       ", Status: types.Present},
				types.BPChar{String: "null    ", Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{
				{Length: 3, LowerBound: 1},
				{Length: 2, LowerBound: 1},
			},
			Status: types.Present,
		},
		&types.BPCharArray{
			Elements: []types.BPChar{
				types.BPChar{String: " bar    ", Status: types.Present},
				types.BPChar{String: "    baz ", Status: types.Present},
				types.BPChar{String: "    quz ", Status: types.Present},
				types.BPChar{String: "foo     ", Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{
				{Length: 2, LowerBound: 4},
				{Length: 2, LowerBound: 2},
			},
			Status: types.Present,
		},
	})
}
