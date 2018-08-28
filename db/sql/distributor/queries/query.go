package queries

import (
	"github.com/Ready-Stock/Noah/db/sql/context"
)

type IQueryStatement interface {
	Execute(ctx *context.NContext)
	compilePlan(ctx *context.NContext, nodes []int) ([]context.NodeExecutionPlan, error)
	getTargetNodes(ctx *context.NContext) ([]int, error)
	getAccountIDs() ([]int, error)
}

