package distributor

import (
	"database/sql"
	"sync"
	_ "github.com/lib/pq"
	"github.com/Ready-Stock/Noah/Prototype/datums"
)

type QueryResult struct {
	Rows   *sql.Rows
	Error  error
	NodeID int
}

func DistributeQuery(query string, nodes ...datums.Node) []QueryResult {
	responses := make([]QueryResult, len(nodes))
	var wg sync.WaitGroup
	wg.Add(len(nodes))
	for index, node := range nodes {
		go func(index int, node datums.Node) {
			defer wg.Done()
			if db, err := sql.Open("postgres", node.ConnectionString); err != nil {
				responses[index] = QueryResult{
					Error:  err,
					NodeID: node.NodeID,
				}
			} else {
				if rows, err := db.Query(query); err != nil {
					responses[index] = QueryResult{
						Error:  err,
						NodeID: node.NodeID,
					}
				} else {
					responses[index] = QueryResult{
						Rows:   rows,
						NodeID: node.NodeID,
					}
				}
			}
		}(index, node)
	}
	wg.Wait()
	return responses
}
