package sql

import (
	"github.com/readystock/noah/db/system"
	"testing"
)

func Test_Create_CompilePlan(t *testing.T) {
	stmt := CreateStatement{}
	plans, err := stmt.compilePlan(nil, make([]system.NNode, 1))
	if err != nil {
		panic(err)
	}

	if len(plans) == 0 {
		panic("no plans returned")
	}
}
