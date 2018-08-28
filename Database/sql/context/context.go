package context

import (
	"github.com/Ready-Stock/Noah/Database/system"
	"github.com/Ready-Stock/badger"
	"fmt"
	"strconv"
	"github.com/Ready-Stock/pgx"
)

const (
	AccountsNodesPath = "/accounts/%d/nodes/"
	NodesPath         = "/nodes/"
)

type NContext struct {
	SystemContext system.SContext
	Nodes map[int]*pgx.Tx
}

type NodeExecutionPlan struct {
	CompiledQuery string
	NodeID        int
}

type NodeExecutionResult struct {
	NodeExecutionPlan
	Error error
}

func (ctx *NContext) GetNodesForAccountID(AccountID *int) ([]int, error) {
	node_ids := make([]int, 0)
	search_prefix := make([]byte, 0)
	if AccountID == nil {
		search_prefix = []byte(NodesPath)
	} else {
		search_prefix = []byte(fmt.Sprintf(AccountsNodesPath, AccountID))
	}
	if err := ctx.SystemContext.Badger.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(search_prefix); it.ValidForPrefix(search_prefix); it.Next() {
			k := it.Item().Key()
			if id, err := strconv.Atoi(string(k)); err != nil {
				return err
			} else {
				node_ids = append(node_ids, id)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return node_ids, nil
}

func (ctx *NContext) ExecutePlans(plans []NodeExecutionPlan) error {


	return nil
}
