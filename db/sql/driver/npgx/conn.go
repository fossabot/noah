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
	"crypto/md5"
	"encoding/hex"
	"github.com/Ready-Stock/Noah/db/sql/driver"
	"github.com/Ready-Stock/Noah/db/sql/pgio"
	"github.com/Ready-Stock/Noah/db/sql/pgwire/pgproto"
	"github.com/Ready-Stock/Noah/db/sql/types"
	"github.com/Ready-Stock/pgx/pgtype"
	"github.com/kataras/go-errors"
	"github.com/kataras/golog"
	"io"
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
	wbuf					  []byte
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
		case *pgproto.BackendKeyData:
		case *pgproto.Authentication:
			if err = c.rxAuthenticationX(msg); err != nil {
				return err
			}
		case *pgproto.ReadyForQuery:
			c.rxReadyForQuery(msg)
			golog.Debugf("connection established to postgres %s", c.config.Host)

			if c.ConnInfo == minimalConnInfo {
				err = c.initConnInfo()
				if err != nil {
					return err
				}
			}

			return nil
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

func (c *Conn) initConnInfo() (err error) {
	var (
		connInfo *types.ConnInfo
	)

	if connInfo, err = initPostgresql(c); err == nil {
		c.ConnInfo = connInfo
		return err
	}

	// Check if CrateDB specific approach might still allow us to connect.
	if connInfo, err = c.crateDBTypesQuery(err); err == nil {
		c.ConnInfo = connInfo
	}

	return err
}

func initPostgresql(c *Conn) (*pgtype.ConnInfo, error) {
	const (
		namedOIDQuery = `select t.oid,
	case when nsp.nspname in ('pg_catalog', 'public') then t.typname
		else nsp.nspname||'.'||t.typname
	end
from pg_type t
left join pg_type base_type on t.typelem=base_type.oid
left join pg_namespace nsp on t.typnamespace=nsp.oid
where (
	  t.typtype in('b', 'p', 'r', 'e')
	  and (base_type.oid is null or base_type.typtype in('b', 'p', 'r'))
	)`
	)

	nameOIDs, err := connInfoFromRows(c.Query(namedOIDQuery))
	if err != nil {
		return nil, err
	}

	cinfo := pgtype.NewConnInfo()
	cinfo.InitializeDataTypes(nameOIDs)

	if err = c.initConnInfoEnumArray(cinfo); err != nil {
		return nil, err
	}

	if err = c.initConnInfoDomains(cinfo); err != nil {
		return nil, err
	}

	return cinfo, nil
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

func (c *Conn) rxAuthenticationX(msg *pgproto.Authentication) (err error) {
	switch msg.Type {
	case pgproto.AuthTypeOk:
	case pgproto.AuthTypeCleartextPassword:
		err = c.txPasswordMessage(c.config.Password)
	case pgproto.AuthTypeMD5Password:
		digestedPassword := "md5" + hexMD5(hexMD5(c.config.Password+c.config.User)+string(msg.Salt[:]))
		err = c.txPasswordMessage(digestedPassword)
	default:
		err = errors.New("Received unknown authentication message")
	}

	return
}

func hexMD5(s string) string {
	hash := md5.New()
	io.WriteString(hash, s)
	return hex.EncodeToString(hash.Sum(nil))
}

func (c *Conn) txPasswordMessage(password string) (err error) {
	buf := c.wbuf
	buf = append(buf, 'p')
	sp := len(buf)
	buf = pgio.AppendInt32(buf, -1)
	buf = append(buf, password...)
	buf = append(buf, 0)
	pgio.SetInt32(buf[sp:], int32(len(buf[sp:])))

	_, err = c.conn.Write(buf)

	return err
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
		return errors.New("notice response not supported at this time.")
	case *pgproto.NotificationResponse:
		return errors.New("notification response not supported at this time.")
	case *pgproto.ReadyForQuery:
		c.rxReadyForQuery(msg)
	case *pgproto.ParameterStatus:
		return nil
	}

	return nil
}

func connInfoFromRows(rows *Rows, err error) (map[string]types.OID, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nameOIDs := make(map[string]types.OID, 256)
	for rows.Next() {
		var oid types.OID
		var name types.Text
		if err = rows.Scan(&oid, &name); err != nil {
			return nil, err
		}

		nameOIDs[name.String] = oid
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return nameOIDs, err
}

func (c *Conn) Query(sql string, args ...interface{}) ([][]byte, error) {

}