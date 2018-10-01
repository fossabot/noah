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
 */

package system

import (
	"github.com/Ready-Stock/Noah/db/sql/driver"
	"github.com/Ready-Stock/Noah/db/sql/driver/npgx"
	"sync"
)

type NodePool struct {
	base      *BaseContext
	sync      sync.Mutex
	nodePools map[uint64]*npgx.ConnPool
}

func (pool *NodePool) AcquireTransaction(nodeId uint64) (*npgx.Transaction, error) {
	if conn, err := pool.AcquireConnection(nodeId); err != nil {
		return nil, err
	} else {
		return conn.Begin()
	}
}

func (pool *NodePool) ReleaseTransaction(nodeId uint64, tx *npgx.Transaction) {

}

func (pool *NodePool) AcquireConnection(nodeId uint64) (*npgx.Conn, error) {
	if nodePool, ok := pool.nodePools[nodeId]; !ok {
		// Init a new connection
		if node, err := pool.base.GetNode(nodeId); err != nil {
			return nil, err
		} else {
			return npgx.Connect(driver.ConnConfig{
				Host:     node.IPAddress,
				Port:     node.Port,
				Database: node.Database,
				User:     node.User,
				Password: node.Password,
			})
		}
	} else {
		return nodePool.Acquire()
	}
}

func (pool *NodePool) ReleaseConnection(nodeId uint64, conn *npgx.Conn) {
	if nodePool, ok := pool.nodePools[nodeId]; !ok {
		conn.Close()
	} else {
		nodePool.Release(conn)
	}
}