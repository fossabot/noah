package sql

import (
"github.com/Ready-Stock/Noah/db/sql/plan"
)

type IQueryStatement interface {
	Execute(ex *connExecutor, res RestrictedCommandResult)
	compilePlan(ex *connExecutor, nodes []uint64) ([]plan.NodeExecutionPlan, error)
	getTargetNodes(ex *connExecutor) ([]int, error)
}

