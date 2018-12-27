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
    "github.com/readystock/noah/db/sql/driver"
    "github.com/readystock/noah/db/sql/driver/npgx"
    "sync"
)

type SPool struct {
    *baseContext
    sync      *sync.Mutex
    nodePools map[uint64]*npgx.ConnPool
}

func (pool *SPool) AcquireTransaction(nodeId uint64) (*npgx.Transaction, error) {
    if conn, err := pool.AcquireConnection(nodeId); err != nil {
        return nil, err
    } else {
        return conn.Begin()
    }
}

func (pool *SPool) AcquireConnection(nodeId uint64) (*npgx.Conn, error) {
    if nodePool, ok := pool.nodePools[nodeId]; !ok {
        // Init a new connection
        sNode := SNode(*pool.baseContext)
        if node, err := (&sNode).GetNode(nodeId); err != nil {
            return nil, err
        } else {
            return npgx.Connect(driver.ConnConfig{
                Host:     node.Address,
                Port:     uint16(node.Port),
                Database: node.Database,
                User:     node.User,
                Password: node.Password,
            })
        }
    } else {
        return nodePool.Acquire()
    }
}
