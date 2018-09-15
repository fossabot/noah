package system

import (
	"github.com/Ready-Stock/badger"
	"github.com/Ready-Stock/pgx"
	"encoding/json"
	"strconv"
	"fmt"
		"github.com/Ready-Stock/Noah/conf"
	"github.com/kataras/go-errors"
)

const (
	NodesPath    = "/nodes/"
	NodesIDsPath = "/sequences/nodes/ids"
	PreloadPoolConnectionCount = 5
)

type SContext struct {
	Badger        *badger.DB
	NodeIDs		  *badger.Sequence
	Flags 		  SFlags
	node_info     map[int]*NNode
	node_pool     map[int]chan *pgx.Conn
}

type SFlags struct {
	HTTPPort int
	PostgresPort int
	DataDirectory string
}

func (ctx *SContext) GetNodes() (n []NNode, e error) {
	n = make([]NNode, 0)
	e = ctx.Badger.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(NodesPath)
		ctx.node_info = map[int]*NNode{}
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			v, err := item.Value()
			if err != nil {
				return err
			}
			node := NNode{}
			if err := json.Unmarshal(v, &node); err != nil {
				return err
			}
			n = append(n, node)
		}
		return nil
	})
	return n, e
}

func (ctx *SContext) AddNode(node NNode) (error) {
	return ctx.Badger.Update(func(txn *badger.Txn) error {
		node_id, err := ctx.NodeIDs.Next()
		if err != nil {
			return err
		}
		json, err := json.Marshal(node)
		if err != nil {
			return err
		}
		if err := txn.Set([]byte(fmt.Sprintf("%s%d", NodesPath, node_id)), json); err != nil {
			txn.Discard()
			return err
		} else {
			return txn.Commit(nil)
		}
	})
}




func (ctx *SContext) loadStoredNodes() error {
	return ctx.Badger.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(NodesPath)
		ctx.node_info = map[int]*NNode{}
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			v, err := item.Value()
			if err != nil {
				return err
			}
			node := NNode{}
			if err := json.Unmarshal(v, &node); err != nil {
				return err
			}
			ctx.node_info[node.NodeID] = &node
		}
		return nil
	})
}

func (ctx *SContext) StartConnectionPool() error {
	if err := ctx.loadStoredNodes(); err != nil {
		return err
	}
	ctx.node_pool = map[int]chan *pgx.Conn{}
	for _, node := range ctx.node_info {
		conn, err := pgx.Connect(pgx.ConnConfig{
			Host:      node.IPAddress,
			Port:      node.Port,
			Database:  node.Database,
			User:      node.User,
			Password:  node.Password,
			TLSConfig: nil,
		})
		if err != nil {
			return err
		}
		ctx.node_pool[node.NodeID] = make(chan *pgx.Conn, PreloadPoolConnectionCount)
		ctx.node_pool[node.NodeID] <- conn
	}
	return nil
}

func (ctx *SContext) GetConnection(NodeID int) (*pgx.Conn, error) {
	if pool, ok := ctx.node_pool[NodeID]; !ok || len(pool) == 0 {
		// It's possible (somehow) for the node to have info but not be in the pool. If the info doesn't exist return an error.
		if node, ok := ctx.node_info[NodeID]; !ok {
			return nil, errors.New("could not retrieve connection details for node ID %d").Format(NodeID)
		} else {
			conn, err := pgx.Connect(pgx.ConnConfig{
				Host:      node.IPAddress,
				Port:      node.Port,
				Database:  node.Database,
				User:      node.User,
				Password:  node.Password,
				TLSConfig: nil,
			})
			if err != nil {
				return nil, err
			}
			ctx.node_pool[NodeID] = make(chan *pgx.Conn, PreloadPoolConnectionCount)
			return conn, nil
		}
	} else {
		return <-pool, nil
	}
}

func (ctx *SContext) ReturnConnection(NodeID int, conn *pgx.Conn) {
	if pool, ok := ctx.node_pool[NodeID]; !ok || len(pool) >= PreloadPoolConnectionCount {
		// If we cannot return the connection to the pool, close the connection.
		// Or if we already have 5 connections pooled, then close this extra connection.
		conn.Close()
	} else {
		pool <- conn
	}
}

func (ctx *SContext) GetNextNodeID() (*int, error) {
	s, err := ctx.Badger.GetSequence([]byte(NodesIDsPath), 1)
	if err != nil {
		return nil, err
	}
	val, err := s.Next()
	if err != nil {
		return nil, err
	}
	s_string := strconv.FormatUint(val, 10)
	s_int, _ := strconv.Atoi(s_string)
	return &s_int, nil
}

func (ctx *SContext) UpsertNodes(nodes []conf.NodeConfig) error {
	tx := ctx.Badger.NewTransaction(true)
	for i := 0; i < len(nodes); i++ {
		node := nodes[i]
		n := NNode{
			NodeID:    node.NodeID,
			Region:    "",
			IPAddress: node.Address,
			Port:      node.Port,
			Database:  node.Database,
			User:      node.User,
			Password:  node.Password,
		}
		njson, err := json.Marshal(n)
		if err != nil {
			return err
		}
		if err := tx.Set([]byte(fmt.Sprintf("%s%d", NodesPath, n.NodeID)), njson); err != nil {
			tx.Discard()
			return err
		}
	}
	return tx.Commit(nil)
}
