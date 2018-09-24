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
 * License (Apache Public License 2.0) https://github.com/cockroachdb/cockroach/blob/master/LICENSE
 *
 *
 *
 *
 */

package npgx

import (
	"github.com/Ready-Stock/Noah/db/sql/driver"
	"github.com/Ready-Stock/Noah/db/sql/pgwire/pgproto"
	"github.com/Ready-Stock/Noah/db/sql/types"
	"github.com/kataras/go-errors"
	"net"
	"sync"
	"time"
)

var (
	// Errors

	// ErrNoRows occurs when rows are expected but none are returned.
	ErrNoRows = errors.New("no rows in result set")
	// ErrDeadConn occurs on an attempt to use a dead connection
	ErrDeadConn = errors.New("conn is dead")
	// ErrConnBusy occurs when the connection is busy (for example, in the middle of
	// reading query results) and another action is attempted.
	ErrConnBusy = errors.New("conn is busy")
)

var (
	// Defaults
	minimalConnInfo *types.ConnInfo
)

const (
	connStatusUnitialized = iota
	connStatusClosed
	connStatusIdle
	connStatusBusy
)

func init() {
	minimalConnInfo = types.NewConnInfo()
	minimalConnInfo.InitializeDataTypes(map[string]types.OID{
		"int4":    types.Int4OID,
		"name":    types.NameOID,
		"oid":     types.OIDOID,
		"text":    types.TextOID,
		"varchar": types.VarcharOID,
	})
}

type DialFunc func(network, addr string) (net.Conn, error)

type Conn struct {
	conn                      net.Conn
	lastActivityTime          time.Time // the last time the connection was used
	config                    driver.ConnConfig
	txStatus                  byte
	status                    byte
	mux                       sync.Mutex
	frontend                  *pgproto.Frontend
	pendingReadyForQueryCount int // numer of ReadyForQuery messages expected
	causeOfDeath              error
	// Public Vars
	ConnInfo *types.ConnInfo
}

func Connect(config driver.ConnConfig) (c *Conn, err error) {
	return connect(config, minimalConnInfo)
}

func connect(config driver.ConnConfig, connInfo *types.ConnInfo) (c *Conn, err error) {
	c = new(Conn)

	c.config = config
	c.ConnInfo = connInfo

	if c.config.User == "" {
		c.config.User = "postgres"
	}

	if c.config.Port == 0 {
		c.config.Port = 5432
	}

	if c.config.Database != "" {
		c.config.Database = "postgres"
	}

	network, address := c.config.NetworkAddress()
	d := defaultDialer()
	if err := c.connect(config, network, address, d.Dial); err != nil {
		return nil, err
	} else {
		return c, nil
	}
}

func (c *Conn) connect(config driver.ConnConfig, network, address string, dial DialFunc) (err error) {
	c.conn, err = dial(network, address)
	if err != nil {
		return err
	}
	defer func() {
		if c != nil && err != nil {
			c.conn.Close()
			c.mux.Lock()
			c.status = connStatusClosed
			c.mux.Unlock()
		}
	}()

	c.lastActivityTime = time.Now()

	c.mux.Lock()
	c.status = connStatusIdle
	c.mux.Unlock()

	c.frontend, err = pgproto.NewFrontend(c.conn, c.conn)
	if err != nil {
		return err
	}

	// Start building the startup message for connecting to the database.00
	startupMsg := pgproto.StartupMessage{
		ProtocolVersion: pgproto.ProtocolVersionNumber,
		Parameters: map[string]string{
			"user":     c.config.User,
			"database": c.config.Database,
		},
	}

	if _, err := c.conn.Write(startupMsg.Encode(nil)); err != nil {
		return err
	}
	c.pendingReadyForQueryCount = 1

	for {
		msg, err := c.rxMsg()
		if err != nil {
			return err
		}

		switch msg := msg.(type) {
		case pgproto.BackendKeyData:
		case pgproto.Authentication:
		case pgproto.ReadyForQuery:
		default:
			if err = c.processContextFreeMsg(msg); err != nil {
				return err
			}
		}
	}
}

func defaultDialer() *net.Dialer {
	return &net.Dialer{KeepAlive: 5 * time.Minute}
}

func (c *Conn) IsAlive() bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.status >= connStatusIdle
}

func (c *Conn) rxMsg() (pgproto.BackendMessage, error) {
	if !c.IsAlive() {
		return nil, ErrDeadConn
	}

	msg, err := c.frontend.Receive()
	if err != nil {
		if netErr, ok := err.(net.Error); !(ok && netErr.Timeout()) {
			c.die(err)
		}
		return nil, err
	}

	c.lastActivityTime = time.Now()

	// fmt.Printf("rxMsg: %#v\n", msg)

	return msg, nil
}

func (c *Conn) die(err error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.status == connStatusClosed {
		return
	}

	c.status = connStatusClosed
	c.causeOfDeath = err
	c.conn.Close()
}


func (c *Conn) processContextFreeMsg(msg pgproto.BackendMessage) (err error) {
	switch msg := msg.(type) {
	case *pgproto.ErrorResponse:
		return c.rxErrorResponse(msg)
	case *pgproto.NoticeResponse:
		c.rxNoticeResponse(msg)
	case *pgproto.NotificationResponse:
		c.rxNotificationResponse(msg)
	case *pgproto.ReadyForQuery:
		c.rxReadyForQuery(msg)
	case *pgproto.ParameterStatus:
		c.rxParameterStatus(msg)
	}

	return nil
}

