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

package npgx

import (
	"context"
	"github.com/kataras/go-errors"
	"github.com/readystock/golog"
	"github.com/readystock/noah/db/sql/types"
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
	ConnConfig
	MaxConnections int               // max simultaneous connections to use, default 5, must be at least 2
	AfterConnect   func(*Conn) error // function to call on every new connection
	AcquireTimeout time.Duration     // max wait time when all connections are busy (0 means no timeout)
}

type ConnPool struct {
	allConnections       []*Conn
	availableConnections []*Conn
	cond                 *sync.Cond
	config               ConnConfig // config used when establishing connection
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
	// golog.Warnf("releasing connection for node [%d]", conn.GetNodeId())
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
	// golog.Verbosef("acquiring connection for node [%d]", p.config.NodeId)
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

// Begin acquires a connection and begins a transaction on it. When the
// transaction is closed the connection will be automatically released.
func (p *ConnPool) Begin() (*Transaction, error) {
	return p.BeginEx(context.Background(), nil)
}

// BeginEx acquires a connection and starts a transaction with txOptions
// determining the transaction mode. When the transaction is closed the
// connection will be automatically released.
func (p *ConnPool) BeginEx(ctx context.Context, txOptions *TransactionOptions) (*Transaction, error) {
	for {
		c, err := p.Acquire()
		if err != nil {
			return nil, err
		}

		tx, err := c.BeginEx(ctx, txOptions)
		if err != nil {
			alive := c.IsAlive()
			p.Release(c)

			// If connection is still alive then the error is not something trying
			// again on a new connection would fix, so just return the error. But
			// if the connection is dead try to acquire a new connection and try
			// again.
			if alive {
				return nil, err
			}
			continue
		}

		tx.connPool = p
		return tx, nil
	}
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
