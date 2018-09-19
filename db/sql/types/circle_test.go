package types_test

import (
	"testing"

	"github.com/Ready-Stock/Noah/db/sql/types"
	"github.com/Ready-Stock/Noah/db/sql/types/testutil"
)

func TestCircleTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "circle", []interface{}{
		&types.Circle{P: types.Vec2{1.234, 5.67890123}, R: 3.5, Status: types.Present},
		&types.Circle{P: types.Vec2{-1.234, -5.6789}, R: 12.9, Status: types.Present},
		&types.Circle{Status: types.Null},
	})
}
