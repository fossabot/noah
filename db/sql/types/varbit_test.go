package types_test

import (
	"testing"

	"github.com/Ready-Stock/Noah/db/sql/types"
	"github.com/Ready-Stock/Noah/db/sql/types/testutil"
)

func TestVarbitTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "varbit", []interface{}{
		&types.Varbit{Bytes: []byte{}, Len: 0, Status: types.Present},
		&types.Varbit{Bytes: []byte{0, 1, 128, 254, 255}, Len: 40, Status: types.Present},
		&types.Varbit{Bytes: []byte{0, 1, 128, 254, 128}, Len: 33, Status: types.Present},
		&types.Varbit{Status: types.Null},
	})
}

func TestVarbitNormalize(t *testing.T) {
	testutil.TestSuccessfulNormalize(t, []testutil.NormalizeTest{
		{
			SQL:   "select B'111111111'",
			Value: &types.Varbit{Bytes: []byte{255, 128}, Len: 9, Status: types.Present},
		},
	})
}
