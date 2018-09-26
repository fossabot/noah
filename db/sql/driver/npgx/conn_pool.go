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

package npgx

import (
	"github.com/Ready-Stock/Noah/db/sql/driver"
	"github.com/Ready-Stock/Noah/db/sql/types"
	"github.com/kataras/go-errors"
	"sync"
	"time"
)

var (
	// Errors

	// ErrAcquireTimeout occurs when an attempt to acquire a connection times out.
	ErrAcquireTimeout = errors.New("timeout acquiring connection from pool")
	// ErrClosedPool occurs on an attempt to acquire a connection from a closed pool.
	ErrClosedPool = errors.New("cannot acquire from closed pool")
)

type ConnPoolConfig struct {
	driver.ConnConfig
	MaxConnections int               // max simultaneous connections to use, default 5, must be at least 2
	AfterConnect   func(*Conn) error // function to call on every new connection
	AcquireTimeout time.Duration     // max wait time when all connections are busy (0 means no timeout)
}

type ConnPool struct {
	allConnections       []*Conn
	availableConnections []*Conn
	cond                 *sync.Cond
	config               driver.ConnConfig // config used when establishing connection
	inProgressConnects   int
	maxConnections       int
	resetCount           int
	afterConnect         func(*Conn) error
	closed               bool
	preparedStatements   map[string]*PreparedStatement
	acquireTimeout       time.Duration
	connInfo             *types.ConnInfo
}

type ConnPoolStat struct {
	MaxConnections       int // max simultaneous connections to use
	CurrentConnections   int // current live connections
	AvailableConnections int // unused live connections
}


// NewConnPool creates a new ConnPool. config.ConnConfig is passed through to
// Connect directly.
func NewConnPool(config ConnPoolConfig) (p *ConnPool, err error) {
	p = new(ConnPool)
	p.config = config.ConnConfig
	p.connInfo = minimalConnInfo
	p.maxConnections = config.MaxConnections
	if p.maxConnections == 0 {
		p.maxConnections = 5
	}
	if p.maxConnections < 1 {
		return nil, errors.New("MaxConnections must be at least 1")
	}
	p.acquireTimeout = config.AcquireTimeout
	if p.acquireTimeout < 0 {
		return nil, errors.New("AcquireTimeout must be equal to or greater than 0")
	}

	p.afterConnect = config.AfterConnect

	p.allConnections = make([]*Conn, 0, p.maxConnections)
	p.availableConnections = make([]*Conn, 0, p.maxConnections)
	p.preparedStatements = make(map[string]*PreparedStatement)
	p.cond = sync.NewCond(new(sync.Mutex))

	// Initially establish one connection
	var c *Conn
	c, err = p.createConnection()
	if err != nil {
		return
	}
	p.allConnections = append(p.allConnections, c)
	p.availableConnections = append(p.availableConnections, c)
	p.connInfo = c.ConnInfo.DeepCopy()

	return
}

func (p *ConnPool) createConnection() (*Conn, error) {
	c, err := connect(p.config, p.connInfo)
	if err != nil {
		return nil, err
	}
	return p.afterConnectionCreated(c)
}

// afterConnectionCreated executes (if it is) afterConnect() callback and prepares
// all the known statements for the new connection.
func (p *ConnPool) afterConnectionCreated(c *Conn) (*Conn, error) {
	if p.afterConnect != nil {
		err := p.afterConnect(c)
		if err != nil {
			c.die(err)
			return nil, err
		}
	}

	for _, ps := range p.preparedStatements {
		if _, err := c.Prepare(ps.Name, ps.SQL); err != nil {
			c.die(err)
			return nil, err
		}
	}

	return c, nil
}


// Release gives up use of a connection.
func (p *ConnPool) Release(conn *Conn) {
	if conn.ctxInProgress {
		panic("should never release when context is in progress")
	}

	if conn.txStatus != 'I' {
		conn.Exec("rollback")
	}


	p.cond.L.Lock()

	if conn.poolResetCount != p.resetCount {
		conn.Close()
		p.cond.L.Unlock()
		p.cond.Signal()
		return
	}

	if conn.IsAlive() {
		p.availableConnections = append(p.availableConnections, conn)
	} else {
		p.removeFromAllConnections(conn)
	}
	p.cond.L.Unlock()
	p.cond.Signal()
}
