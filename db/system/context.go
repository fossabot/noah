package system

import (
	"encoding/json"
	"fmt"
	"github.com/Ready-Stock/badger"
	"github.com/Ready-Stock/pgx"
	"github.com/ahmetb/go-linq"
	"github.com/kataras/go-errors"
)

const (
	NodesPath    = "/nodes/"
	SettingsPath = "/settings/"
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

func (ctx *SContext) GetNode(NodeID uint64) (n *NNode, e error) {
	node := NNode{}
	e = ctx.Badger.View(func(txn *badger.Txn) error {
		j, err := txn.Get([]byte(fmt.Sprintf("%s%d", NodesPath, NodeID)))
		if err != nil {
			return err
		}
		v, err := j.Value()
		if err != nil {
			return err
		}
		err = json.Unmarshal(v, &node)
		return err
	})
	return &node, e
}

func (ctx *SContext) AddNode(node NNode) (error) {
	return ctx.Badger.Update(func(txn *badger.Txn) error {
		existing_nodes, err := ctx.GetNodes()
		if linq.From(existing_nodes).AnyWithT(func(existing NNode) bool {
			return node.Database == existing.Database && node.IPAddress == existing.IPAddress && node.Port == existing.Port
		}) {
			return errors.New("a node already exists with the same connection string.")
		}
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

func (ctx *SContext) GetSettings() (*map[string]string, error) {
	m := map[string]string{}
	e := ctx.Badger.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(SettingsPath)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			v, err := item.Value()
			if err != nil {
				return err
			}
			m[string(item.Key()[len(prefix) - 1:])] = string(v)
		}
		return nil
	})
	return &m, e
}

func (ctx *SContext) SetSetting(SettingName string) (error) {
	return nil
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
