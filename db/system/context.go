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
 *
 * This application uses Open Source components. You can find the
 * source code of their open source projects along with license
 * information below. We acknowledge and are grateful to these
 * developers for their contributions to open source.
 *
 * Project: CockroachDB https://github.com/cockroachdb/cockroach
 * Copyright 2018 The Cockroach Authors.
 * License (Apache License 2.0) https://github.com/cockroachdb/cockroach/blob/master/LICENSE
 *
 * Project: Vitess https://github.com/vitessio/vitess
 * Copyright 2018 Google Inc.
 * License (Apache License 2.0) https://github.com/vitessio/vitess/blob/master/LICENSE
 *
 * Project: Citus https://github.com/citusdata/citus
 * Copyright 2018 Citus Data, Inc.
 * License (GNU Affero General Public License v3.0) https://github.com/citusdata/citus/blob/master/LICENSE
 *
 * Project: pg_query_go https://github.com/lfittl/pg_query_go
 * Copyright 2018 Lukas Fittl
 * License (3-Clause BSD) https://github.com/lfittl/pg_query_go/blob/master/LICENSE
 *
 * Project: pgx https://github.com/jackc/pgx
 * Copyright 2018 Jack Christensen
 * License (MIT) https://github.com/jackc/pgx/blob/master/LICENSE
 *
 * Project: BadgerDB https://github.com/dgraph-io/badger
 * Copyright 2018 Dgraph Labs, Inc. and Contributors
 * License (MIT) https://github.com/dgraph-io/badger/blob/master/LICENSE
 */

package system

import (
	"encoding/json"
	"fmt"
	"github.com/Ready-Stock/badger"
	"github.com/ahmetb/go-linq"
	"github.com/kataras/go-errors"
)

const (
	NodesPath                  = "/nodes/"
	SettingsPath               = "/settings/"
	PreloadPoolConnectionCount = 5
)

type SContext struct {
	Badger    *badger.DB
	WalBadger *badger.DB
	NodeIDs   *badger.Sequence
	Flags     SFlags
	node_info map[uint64]*NNode
	// node_pool     map[uint64]chan *pgx.Conn
}

type SFlags struct {
	HTTPPort      int
	PostgresPort  int
	DataDirectory string
	WalDirectory  string
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
			m[string(item.Key()[len(prefix)-1:])] = string(v)
		}
		return nil
	})
	return &m, e
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
