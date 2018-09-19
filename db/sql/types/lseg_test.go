package types_test

import (
	"testing"

	"github.com/Ready-Stock/Noah/db/sql/types"
	"github.com/Ready-Stock/Noah/db/sql/types/testutil"
)

func TestLsegTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "lseg", []interface{}{
		&types.Lseg{
			P:      [2]types.Vec2{{3.14, 1.678}, {7.1, 5.2345678901}},
			Status: types.Present,
		},
		&types.Lseg{
			P:      [2]types.Vec2{{7.1, 1.678}, {-13.14, -5.234}},
			Status: types.Present,
		},
		&types.Lseg{Status: types.Null},
	})
}
