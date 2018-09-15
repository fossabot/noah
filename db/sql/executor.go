package sql

import (
	"github.com/Ready-Stock/Noah/db/sql/plan"
	"github.com/kataras/golog"
)

func (ex *connExecutor) ExecutePlans(plans []plan.NodeExecutionPlan) error {
	for i, p := range plans {
		go func(ex *connExecutor, index int, pln plan.NodeExecutionPlan) {
			golog.Infof("[%s] Executing query: `%s` on node %d", ex.ClientAddress, pln.CompiledQuery, pln.NodeID)
			if node, ok := ex.Nodes[pln.NodeID]; !ok {
				golog.Debug()
			}
		}(ex, i, p)
	}
	return nil
}
