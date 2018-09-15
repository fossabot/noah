package sql

import (
	"github.com/Ready-Stock/Noah/db/sql/plan"
	"github.com/Ready-Stock/pg_query_go"
	. "github.com/ahmetb/go-linq"
	"github.com/kataras/go-errors"
)

func (stmt *SelectStatement) getTargetNodes(ex *connExecutor) ([]uint64, error) {
	accounts, err := stmt.getAccountIDs()
	if err != nil {
		return nil, err
	}

	if len(accounts) == 1 {
		return ex.GetNodesForAccountID(&accounts[0])
	} else if len(accounts) > 1 {
		node_ids := make([]uint64, 0)
		From(accounts).SelectManyT(func(id uint64) Query {
			if ids, err := ex.GetNodesForAccountID(&id); err == nil {
				return From(ids)
			}
			return From(make([]uint64, 0))
		}).Distinct().ToSlice(&node_ids)
		if len(node_ids) == 0 {
			return nil, errors.New("could not find nodes for account IDs")
		} else {
			return node_ids, nil
		}
	} else {
		return ex.GetNodesForAccountID(nil)
	}
}

func (stmt *SelectStatement) getAccountIDs() ([]uint64, error) {

	return nil, nil
}

func (stmt *SelectStatement) compilePlan(ex *connExecutor, nodes []uint64) ([]plan.NodeExecutionPlan, error) {
	plans := make([]plan.NodeExecutionPlan, len(nodes))
	deparsed, err := pg_query.Deparse(stmt.Statement)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(plans); i++ {
		plans[i] = plan.NodeExecutionPlan{
			CompiledQuery: *deparsed,
			NodeID:        nodes[i],
		}
	}
	return plans, nil
}
