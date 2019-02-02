/*
 * Copyright (c) 2019 Ready Stock
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
	"github.com/readystock/noah/db/sql/driver/npgx"
	"sync"
	"time"
)

type SPool struct {
	*baseContext
	sync      *sync.Mutex
	nodePools map[uint64]*npgx.ConnPool
}

func (pool *SPool) Acquire(nodeId uint64) (*npgx.Conn, error) {
	nodePool, ok := pool.nodePools[nodeId]
	if !ok {
		sNode := SNode(*pool.baseContext)
		if node, err := (&sNode).GetNode(nodeId); err != nil {
			return nil, err
		} else {
			p, err := npgx.NewConnPool(npgx.ConnPoolConfig{
				ConnConfig: npgx.ConnConfig{
					NodeId:   node.NodeId,
					Host:     node.Address,
					Port:     uint16(node.Port),
					Database: node.Database,
					User:     node.User,
					Password: node.Password,
				},
				MaxConnections: 5,
				AcquireTimeout: time.Second * 1,
			})
			if err != nil {
				return nil, err
			}

			pool.nodePools[nodeId] = p
			nodePool = p
		}
	}
	return nodePool.Acquire()
}

func (pool *SPool) Release(conn *npgx.Conn) error {
	nodePool, ok := pool.nodePools[conn.GetNodeId()]
	if !ok {
		return conn.Close()
	}

	nodePool.Release(conn)
	return nil
}
