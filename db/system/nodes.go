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

package system

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/golang/protobuf/proto"
	"github.com/kataras/go-errors"
	"strconv"
	"strings"
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

func (ctx *SNode) GetFirstNode(scope NodeScope) (*NNode, error) {
	predicate := func(node NNode) bool {
		if scope == AllNodes {
			return node.IsAlive
		} else if scope == ReplicaNodes {
			return node.IsAlive && node.ReplicaOf > 0
		} else if scope == StandardNodes {
			return node.IsAlive && node.ReplicaOf == 0
		}
		return false
	}

	values, err := ctx.db.GetPrefix([]byte(nodesPath))
	if err != nil {
		return nil, err
	}

	for _, val := range values {
		node := NNode{}
		if err := proto.Unmarshal(val.Value, &node); err != nil {
			return nil, err
		}
		if predicate(node) {
			return &node, nil
		}
	}
	return nil, errors.New("no nodes in the cluster match the scope")
}

func (ctx *SNode) GetNode(nodeId uint64) (*NNode, error) {
	value, err := ctx.db.Get([]byte(fmt.Sprintf("%s%d", nodesPath, nodeId)))
	if err != nil {
		return nil, err
	}
	node := NNode{}
	if len(value) == 0 {
		return nil, nil
	}
	err = proto.Unmarshal(value, &node)
	if err != nil {
		return nil, err
	}
	return &node, err
}

func (ctx *SNode) AddNode(node NNode) (*NNode, error) {
	existingNodes, err := ctx.GetNodes()
	if err != nil {
		return nil, err
	}
	if linq.From(existingNodes).AnyWithT(func(existing NNode) bool {
		return node.Database == existing.Database && node.Address == existing.Address && node.Port == existing.Port
	}) {
		return nil, errors.New("a node already exists with the same connection string.")
	}

	id, err := ctx.db.NextSequenceValueById("_noah.nodes_")
	if err != nil {
		return nil, err
	}

	node.NodeId = *id
	node.IsAlive = false

	b, err := proto.Marshal(&node)
	if err != nil {
		return nil, err
	}
	return &node, ctx.db.Set([]byte(fmt.Sprintf("%s%d", nodesPath, node.NodeId)), b)
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

func (ctx *SNode) GetAccountDistribution() (map[uint64][]uint64, error) {
	path := []byte(accountsNodesPath)
	values, err := ctx.db.GetPrefix(path)
	if err != nil {
		return nil, err
	}
	dist := map[uint64][]uint64{}
	for _, kv := range values {
		accountId, nodeId := getIdsFromAccountNodesPath(kv.Key)
		if _, ok := dist[nodeId]; !ok {
			dist[nodeId] = []uint64{accountId}
		} else {
			dist[nodeId] = append(dist[nodeId], accountId)
		}
	}
	return dist, nil
}

func getIdsFromAccountNodesPath(path []byte) (accountId uint64, nodeId uint64) {
	pathStr := string(path)
	/* This will look at the index of the items in the path and do things
	   0        1             2         3
	     /accounts_nodes/{accountId}/{nodeId}
	*/
	for index, item := range strings.Split(pathStr, "/") {
		switch index {
		case 2: // accountId
			accountId, _ = strconv.ParseUint(item, 10, 64)
		case 3:
			nodeId, _ = strconv.ParseUint(item, 10, 64)
		}
	}
	return accountId, nodeId
}
