package system

import (
	"encoding/json"
	"fmt"
	"github.com/Ready-Stock/badger"
	"github.com/Ready-Stock/pgx"
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
	node_info     map[uint64]*NNode
	node_pool     map[uint64]chan *pgx.Conn
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
		node.NodeID = node_id
		json, err := json.Marshal(node)
		if err != nil {
			return err
		}
		return txn.Set([]byte(fmt.Sprintf("%s%d", NodesPath, node_id)), json);
	})
}




func (ctx *SContext) loadStoredNodes() error {
	return ctx.Badger.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(NodesPath)
		ctx.node_info = map[uint64]*NNode{}
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
