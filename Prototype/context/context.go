package context

import (
	"database/sql"
	"sync"
	_ "github.com/lib/pq"
	"github.com/Ready-Stock/Noah/Prototype/cluster"
	"github.com/kataras/go-errors"
)

type TxnState string

const (
	StateNoTxn        TxnState = "StateNoTxn"
	StateInTxn        TxnState = "StateInTxn"
	StateCommittedTxn TxnState = "StateCommittedTxn"
)

type TxnAction string

const (
	TxnCommit   TxnAction = "TxnCommit"
	TxnRollback TxnAction = "TxnRollback"
)

type SessionContext struct {
	TransactionState TxnState
	Nodes            map[int]NodeContext
}

type NodeContext struct {
	NodeID           int
	TransactionState TxnState
	DB               *sql.DB
	TX               *sql.Tx
}

type DistributedResponse struct {
	Success bool
	Results []QueryResult
	Errors  []error
}
type QueryResult struct {
	Rows   *sql.Rows
	Error  error
	NodeID int
}

func (ctx *SessionContext) DistributeQuery(query string, nodes ...int) DistributedResponse {
	response := DistributedResponse{
		Success: true,
		Results: make([]QueryResult, len(nodes)),
		Errors:  make([]error, 0),
	}
	updated_nodes := make(chan *NodeContext, len(nodes))
	errors := make([]error, len(nodes))
	var wg sync.WaitGroup
	wg.Add(len(nodes))
	for index, node := range nodes {
		go func(index int, node_id int) {
			defer wg.Done()
			var node NodeContext
			if n, ok := ctx.Nodes[node_id]; !ok {
				node = NodeContext{
					NodeID:           node_id,
					DB:               nil,
					TransactionState: StateNoTxn,
				}
			} else {
				node = n
			}
			if node.DB == nil {
				if db, err := sql.Open("postgres", cluster.Nodes[node.NodeID].ConnectionString); err != nil {
					response.Results[index] = QueryResult{
						Error:  err,
						NodeID: node.NodeID,
					}
					errors[index] = err
				} else {
					if tx, err := db.Begin(); err != nil {
						response.Results[index] = QueryResult{
							Error:  err,
							NodeID: node.NodeID,
						}
						errors[index] = err
					} else {
						node.TX = tx
						node.DB = db
						node.TransactionState = StateInTxn
					}
				}
			}

			if rows, err := node.TX.Query(query); err != nil {
				response.Results[index] = QueryResult{
					Error:  err,
					NodeID: node.NodeID,
				}
				errors[index] = err
			} else {
				response.Results[index] = QueryResult{
					Rows:   rows,
					NodeID: node.NodeID,
				}
				errors[index] = nil
			}
			ctx.Nodes[node_id] = node
		}(index, node)
	}
	wg.Wait()
	for _, err := range errors {
		if err != nil {
			response.Errors = append(response.Errors, err)
		}
	}
	response.Success = len(response.Errors) == 0
	return response
}

func (ctx *SessionContext) DoTxnOnAllNodes(action TxnAction, nodes ...int) DistributedResponse {
	response := DistributedResponse{
		Success: true,
		Results: make([]QueryResult, len(nodes)),
		Errors:  make([]error, 0),
	}
	errors := make([]error, len(nodes))
	var wg sync.WaitGroup
	wg.Add(len(nodes))
	for index, node := range nodes {
		go func(index int, node_id int) {
			defer wg.Done()
			switch action {
			case TxnCommit:
				if err := ctx.Nodes[node_id].TX.Commit(); err != nil {
					response.Results[index] = QueryResult{
						Error:  err,
						NodeID: node_id,
					}
					errors[index] = err
				}
			case TxnRollback:
				if err := ctx.Nodes[node_id].TX.Rollback(); err != nil {
					response.Results[index] = QueryResult{
						Error:  err,
						NodeID: node_id,
					}
					errors[index] = err
				}
			}
		}(index, node)
	}
	wg.Wait()
	for _, err := range errors {
		if err != nil {
			response.Errors = append(response.Errors, err)
		}
	}
	response.Success = len(response.Errors) == 0
	return response
}

func (ctx *SessionContext) HandleResponse(response DistributedResponse) error {
	if !response.Success {
		if len(response.Errors) > 0 {
			return response.Errors[0]
		} else {
			return errors.New("could not execute query successfully on all nodes")
		}
	} else {
		return nil
	}
}

func (ctx *SessionContext) GetAllNodes() ([]int) {
	ids := make([]int, 0)
	for _, node := range cluster.Nodes {
		ids = append(ids, node.NodeID)
	}
	return ids
}
