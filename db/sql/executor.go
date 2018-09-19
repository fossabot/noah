package sql

import (
	"github.com/Ready-Stock/Noah/db/sql/plan"
)

func (ex *connExecutor) ExecutePlans(plans []plan.NodeExecutionPlan) error {
	for i, p := range plans {
		go func(ex *connExecutor, index int, pln plan.NodeExecutionPlan) {
			ex.Info("Executing query: `%s` on node %d", pln.CompiledQuery, pln.NodeID)
			// _, ok := ex.Nodes[pln.NodeID]
			// if !ok {
			// 	ex.Warn("Allocating connection to node %d for session.", pln.NodeID)
			// } else {
			//
			// }
		}(ex, i, p)
	}
	return nil
}
