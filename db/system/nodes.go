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
 *
 * Project: Sonyflake https://github.com/sony/sonyflake
 * Copyright 2018 Sony Corporation
 * License (MIT) https://github.com/sony/sonyflake/blob/master/LICENSE
 *
 * Project: Raft https://github.com/hashicorp/raft
 * Copyright 2018 HashiCorp
 * License (MPL-2.0) https://github.com/hashicorp/raft/blob/master/LICENSE
 */

package system

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/golang/protobuf/proto"
	"github.com/kataras/go-errors"
)

type SNode baseContext

type NodeScope int

const (
	StandardNodes NodeScope = 1
	ReplicaNodes  NodeScope = 2
	AllNodes      NodeScope = 4
)

func (ctx *SNode) GetNodes() (nodes []NNode, e error) {
	nodes = make([]NNode, 0)
	values, err := ctx.db.GetPrefix([]byte(nodesPath))
	if err != nil {
		return nil, err
	}
	for _, val := range values {
		node := NNode{}
		if err := proto.Unmarshal(val.Value, &node); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (ctx *SNode) GetLiveNodes(scope NodeScope) (n []NNode, e error) {
	if nodes, err := ctx.GetNodes(); err != nil {
		return nil, err
	} else {
		linq.From(nodes).WhereT(func(node NNode) bool {
			if scope == AllNodes {
				return node.IsAlive
			} else if scope == ReplicaNodes {
				return node.IsAlive && node.ReplicaOf > 0
			} else if scope == StandardNodes {
				return node.IsAlive && node.ReplicaOf == 0
			}
			return false
		}).ToSlice(&n)
	}
	return n, e
}

func (ctx *SNode) GetNode(nodeId uint64) (node NNode, e error) {
	value, err := ctx.db.Get([]byte(fmt.Sprintf("%s%d", nodesPath, nodeId)))
	if err != nil {
		return node, err
	}
	if err := proto.Unmarshal(value, &node); err != nil {
		return node, err
	}
	return node, nil
}

func (ctx *SNode) AddNode(node NNode) error {
	existingNodes, err := ctx.GetNodes()
	if err != nil {
		return err
	}
	if linq.From(existingNodes).AnyWithT(func(existing NNode) bool {
		return node.Database == existing.Database && node.Address == existing.Address && node.Port == existing.Port
	}) {
		return errors.New("a node already exists with the same connection string.")
	}
	id, err := ctx.db.NextSequenceValueById("_noah.nodes_")
	if err != nil {
		return err
	}
	node.NodeId = *id
	b, err := proto.Marshal(&node)
	if err != nil {
		return err
	}
	return ctx.db.Set([]byte(fmt.Sprintf("%s%d", nodesPath, node.NodeId)), b)
}

func (ctx *SNode) SetNodeLive(nodeId uint64, isAlive bool) (err error) {
	path := []byte(fmt.Sprintf("%s%d", nodesPath, nodeId))
	nodeBytes, err := ctx.db.Get(path)
	if err != nil {
		return err
	}
	node := NNode{}
	if err := proto.Unmarshal(nodeBytes, &node); err != nil {
		return err
	}
	node.IsAlive = isAlive
	b, err := proto.Marshal(&node)
	if err != nil {
		return err
	}
	return ctx.db.Set(path, b)
}
