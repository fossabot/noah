/*
 * Copyright (c) 2018 Ready Stock
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package context

import (
	"database/sql"
	"fmt"
	"github.com/kataras/go-errors"
	_ "github.com/lib/pq"
	"github.com/readystock/noah/Prototype/cluster"
	"sync"
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
	Rows    sql.Rows
}
type QueryResult struct {
	Error  error
	NodeID int
}

func (ctx *SessionContext) DistributeQuery(query string, nodes ...int) DistributedResponse {
	response := DistributedResponse{
		Success: true,
		Results: make([]QueryResult, 0),
		Errors:  make([]error, 0),
	}
	updated_nodes := make(chan NodeContext, len(nodes))
	results := make(chan QueryResult, len(nodes))
	errs := make(chan error, len(nodes))
	returned_rows := make(chan *sql.Rows, len(nodes))
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
					results <- QueryResult{
						Error:  err,
						NodeID: node.NodeID,
					}
					errs <- err
					updated_nodes <- node
					returned_rows <- nil
					return
				} else {
					if tx, err := db.Begin(); err != nil {
						results <- QueryResult{
							Error:  err,
							NodeID: node.NodeID,
						}
						errs <- err
						updated_nodes <- node
						returned_rows <- nil
						return
					} else {
						node.TX = tx
						node.DB = db
						node.TransactionState = StateInTxn
					}
				}
			}
			if rows, err := node.TX.Query(query); err != nil {
				results <- QueryResult{
					Error:  err,
					NodeID: node.NodeID,
				}
				errs <- err
				returned_rows <- nil
			} else {
				results <- QueryResult{
					NodeID: node.NodeID,
				}
				errs <- nil
				returned_rows <- rows
			}
			updated_nodes <- node
			return
		}(index, node)
	}
	wg.Wait()
	for i := 0; i < len(nodes); i++ {
		node := <-updated_nodes
		ctx.Nodes[node.NodeID] = node

		err := <-errs
		if err != nil {
			response.Errors = append(response.Errors, err)
		}

		result := <-results
		response.Results = append(response.Results, result)

		rows := <-returned_rows

		if rows != nil {

		}
	}
	response.Success = len(response.Errors) == 0
	return response
}

func (ctx *SessionContext) DoTxnOnAllNodes(action TxnAction, nodes ...int) DistributedResponse {
	fmt.Printf("Sending (%s) to nodes\n", action)
	response := DistributedResponse{
		Success: true,
		Results: make([]QueryResult, 0),
		Errors:  make([]error, 0),
	}
	errs := make(chan error, len(nodes))
	results := make(chan QueryResult, len(nodes))
	var wg sync.WaitGroup
	wg.Add(len(nodes))
	for index, node := range nodes {
		go func(index int, node_id int) {
			defer wg.Done()
			switch action {
			case TxnCommit:
				if err := ctx.Nodes[node_id].TX.Commit(); err != nil {
					results <- QueryResult{
						Error:  err,
						NodeID: node_id,
					}
					errs <- err
				} else {
					results <- QueryResult{
						Error:  nil,
						NodeID: node_id,
					}
					errs <- nil
				}
			case TxnRollback:
				if err := ctx.Nodes[node_id].TX.Rollback(); err != nil {
					results <- QueryResult{
						Error:  err,
						NodeID: node_id,
					}
					errs <- err
				} else {
					results <- QueryResult{
						Error:  nil,
						NodeID: node_id,
					}
					errs <- nil
				}
			default:
				err := errors.New("invalid action type (%s)").Format(action)
				results <- QueryResult{
					Error:  err,
					NodeID: node_id,
				}
				errs <- err
			}
		}(index, node)
	}
	wg.Wait()
	for i := 0; i < len(nodes); i++ {
		err := <-errs
		if err != nil {
			response.Errors = append(response.Errors, err)
		}

		result := <-results
		response.Results = append(response.Results, result)
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

func (ctx *SessionContext) GetAllNodes() []int {
	ids := make([]int, 0)
	for _, node := range cluster.Nodes {
		ids = append(ids, node.NodeID)
	}
	return ids
}
