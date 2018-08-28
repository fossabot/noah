package system

import (
	"github.com/Ready-Stock/badger"
	"github.com/Ready-Stock/pgx"
		"encoding/json"
	"github.com/Ready-Stock/Noah/Configuration"
	"strconv"
	"fmt"
	"github.com/kataras/go-errors"
)

const (
	NodesPath = "/nodes/"
	NodesIDsPath = "/sequences/nodes/ids"
)

type SContext struct {
	Badger    *badger.DB
	node_info map[int]*NNode
	node_pool map[int]struct{
		Pool chan *pgx.Conn
	}
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
	go func(){
		ctx.node_pool = map[int]struct{ Pool chan *pgx.Conn }{}
		for _, node := range ctx.node_info {
			conn, err := pgx.Connect(pgx.ConnConfig{
				Host: node.IPAddress,
				Port: node.Port,
				Database: node.Database,
				User: node.User,
				Password: node.Password,
				TLSConfig: nil,
			})
			if err != nil {
				return
			}
			p := make(chan *pgx.Conn, 1)
			p <- conn
			ctx.node_pool[node.NodeID] = struct{ Pool chan *pgx.Conn }{Pool: p}
		}
	}()

	return nil
}

func (ctx *SContext) GetConnection(NodeID int) (*pgx.Conn, error) {
	if pool, ok := ctx.node_pool[NodeID]; !ok {
		return nil, errors.New("node does not exist in pool")
	} else {
		return <- pool.Pool, nil
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

func (ctx *SContext) UpsertNodes(nodes []Conf.NodeConfig) error {
	tx := ctx.Badger.NewTransaction(true)
	for i := 0; i < len(nodes); i++ {
		node := nodes[i]
		n := NNode{
			NodeID:node.NodeID,
			Region:"",
			IPAddress:node.Address,
			Port:node.Port,
			Database:node.Database,
			User:node.User,
			Password:node.Password,
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
