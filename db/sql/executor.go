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
			Host: "127.0.0.1",
			Port: 5432,
			User: "postgres",
			Password: "Spring!2016",
		},
		{
			Database: "ready_two",
			Host: "127.0.0.1",
			Port: 5432,
			User: "postgres",
			Password: "Spring!2016",
		},
	}
	for i, node_conf := range nodes {
		go func(index int, node pgx.ConnConfig) {
			fmt.Printf("\t[Execute] Executing query: `%s` on node %d database %s\n", plans[index].CompiledQuery, plans[index].NodeID, node.Database)

		}(i, node_conf)
	}
	return nil
}
