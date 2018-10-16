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

package npgx

import (
	"github.com/Ready-Stock/Noah/db/sql/driver"
	"github.com/Ready-Stock/Noah/db/sql/types"
	"github.com/kataras/go-errors"
	"github.com/kataras/golog"
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

// removeFromAllConnections Removes the given connection from the list.
// It returns true if the connection was found and removed or false otherwise.
func (p *ConnPool) removeFromAllConnections(conn *Conn) bool {
	for i, c := range p.allConnections {
		if conn == c {
			p.allConnections = append(p.allConnections[:i], p.allConnections[i+1:]...)
			return true
		}
	}
	return false
}

// Acquire takes exclusive use of a connection until it is released.
func (p *ConnPool) Acquire() (*Conn, error) {
	p.cond.L.Lock()
	c, err := p.acquire(nil)
	p.cond.L.Unlock()
	return c, err
}

// acquire performs acquision assuming pool is already locked
func (p *ConnPool) acquire(deadline *time.Time) (*Conn, error) {
	if p.closed {
		return nil, ErrClosedPool
	}

	// A connection is available
	if len(p.availableConnections) > 0 {
		c := p.availableConnections[len(p.availableConnections)-1]
		c.poolResetCount = p.resetCount
		p.availableConnections = p.availableConnections[:len(p.availableConnections)-1]
		return c, nil
	}

	// Set initial timeout/deadline value. If the method (acquire) happens to
	// recursively call itself the deadline should retain its value.
	if deadline == nil && p.acquireTimeout > 0 {
		tmp := time.Now().Add(p.acquireTimeout)
		deadline = &tmp
	}

	// Make sure the deadline (if it is) has not passed yet
	if p.deadlinePassed(deadline) {
		return nil, ErrAcquireTimeout
	}

	// If there is a deadline then start a timeout timer
	var timer *time.Timer
	if deadline != nil {
		timer = time.AfterFunc(deadline.Sub(time.Now()), func() {
			p.cond.Broadcast()
		})
		defer timer.Stop()
	}

	// No connections are available, but we can create more
	if len(p.allConnections)+p.inProgressConnects < p.maxConnections {
		// Create a new connection.
		// Careful here: createConnectionUnlocked() removes the current lock,
		// creates a connection and then locks it back.
		c, err := p.createConnectionUnlocked()
		if err != nil {
			return nil, err
		}
		c.poolResetCount = p.resetCount
		p.allConnections = append(p.allConnections, c)
		return c, nil
	}

	// All connections are in use and we cannot create more
	golog.Warn("waiting for available connection")


	// Wait until there is an available connection OR room to create a new connection
	for len(p.availableConnections) == 0 && len(p.allConnections)+p.inProgressConnects == p.maxConnections {
		if p.deadlinePassed(deadline) {
			return nil, ErrAcquireTimeout
		}
		p.cond.Wait()
	}

	// Stop the timer so that we do not spawn it on every acquire call.
	if timer != nil {
		timer.Stop()
	}
	return p.acquire(deadline)
}

// deadlinePassed returns true if the given deadline has passed.
func (p *ConnPool) deadlinePassed(deadline *time.Time) bool {
	return deadline != nil && time.Now().After(*deadline)
}

// createConnectionUnlocked Removes the current lock, creates a new connection, and
// then locks it back.
// Here is the point: lets say our pool dialer's OpenTimeout is set to 3 seconds.
// And we have a pool with 20 connections in it, and we try to acquire them all at
// startup.
// If it happens that the remote server is not accessible, then the first connection
// in the pool blocks all the others for 3 secs, before it gets the timeout. Then
// connection #2 holds the lock and locks everything for the next 3 secs until it
// gets OpenTimeout err, etc. And the very last 20th connection will fail only after
// 3 * 20 = 60 secs.
// To avoid this we put Connect(p.config) outside of the lock (it is thread safe)
// what would allow us to make all the 20 connection in parallel (more or less).
func (p *ConnPool) createConnectionUnlocked() (*Conn, error) {
	p.inProgressConnects++
	p.cond.L.Unlock()
	c, err := Connect(p.config)
	p.cond.L.Lock()
	p.inProgressConnects--

	if err != nil {
		return nil, err
	}
	return p.afterConnectionCreated(c)
}