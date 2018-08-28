package _select

import (
	"github.com/Ready-Stock/Noah/Database/sql/context"
	. "github.com/ahmetb/go-linq"
	"github.com/kataras/go-errors"
	"github.com/Ready-Stock/pg_query_go"
)

func (stmt *SelectStatement) getTargetNodes(ctx *context.NContext) ([]int, error) {
	accounts, err := stmt.getAccountIDs()
	if err != nil {
		return nil, err
	}

	if len(accounts) == 1 {
		return ctx.GetNodesForAccountID(&accounts[0])
	} else if len(accounts) > 1 {
		node_ids := make([]int, 0)
		From(accounts).SelectManyT(func(id int) Query {
			if ids, err := ctx.GetNodesForAccountID(&id); err == nil {
				return From(ids)
			}
			return From(make([]int, 0))
		}).Distinct().ToSlice(&node_ids)
		if len(node_ids) == 0 {
			return nil, errors.New("could not find nodes for account IDs")
		} else {
			return node_ids, nil
		}
	} else {
		return ctx.GetNodesForAccountID(nil)
	}
}

func (stmt *SelectStatement) getAccountIDs() ([]int, error) {

	return nil, nil
}

func (stmt *SelectStatement) compilePlan(ctx *context.NContext, nodes []int) ([]context.NodeExecutionPlan, error) {
	plans := make([]context.NodeExecutionPlan, len(nodes))
	deparsed, err := pg_query.Deparse(stmt.Statement)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(plans); i++ {
		plans[i] = context.NodeExecutionPlan{
			CompiledQuery: *deparsed,
			NodeID: nodes[i],
		}
	}
	return plans, nil
}

