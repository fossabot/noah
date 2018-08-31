package sql

import (
	"fmt"
	"github.com/Ready-Stock/Noah/db/sql/plan"
	"github.com/Ready-Stock/pgx"
)

func (ex *connExecutor) ExecutePlans(plans []plan.NodeExecutionPlan) error {
	nodes := []pgx.ConnConfig{
		{
			Database: "ready_one",
		},
		{
			Database: "ready_two",
		},
	}
	for i, node_conf := range nodes {
		fmt.Printf("\t[Execute] Executing query: %s on node %d database %s\n", plans[i].CompiledQuery, plans[i].NodeID, node_conf.Database)
	}
	return nil
}
