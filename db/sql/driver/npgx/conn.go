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
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"github.com/Ready-Stock/Noah/db/sql/driver"
	"github.com/Ready-Stock/Noah/db/sql/pgio"
	"github.com/Ready-Stock/Noah/db/sql/pgwire/pgproto"
	"github.com/Ready-Stock/Noah/db/sql/types"
	"github.com/kataras/go-errors"
	"github.com/kataras/golog"
	"io"
	"net"
	"sync"
	"sync/atomic"
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

// ProtocolError occurs when unexpected data is received from PostgreSQL
type ProtocolError string

func (e ProtocolError) Error() string {
	return string(e)
}

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
	wbuf                      []byte
	lastActivityTime          time.Time // the last time the connection was used
	config                    driver.ConnConfig
	txStatus                  byte
	status                    byte
	mux                       sync.Mutex
	frontend                  *pgproto.Frontend
	causeOfDeath              error
	pid                       uint32 // backend pid
	secretKey                 uint32 // key to use to send a cancel query message to the server
	pendingReadyForQueryCount int    // numer of ReadyForQuery messages expected
	cancelQueryInProgress     int32
	cancelQueryCompleted      chan struct{}

	// context support
	ctxInProgress bool
	doneChan      chan struct{}
	closedChan    chan error

	// Public Vars
	ConnInfo *types.ConnInfo
}

func (c *Conn) Query(sql string, args ...interface{}) (*Rows, error) {
	return nil, nil
}

func (c *Conn) QueryBytes(sql string, args ...interface{}) ([][]byte, error) {
	return nil, nil
}

// Prepare creates a prepared statement with name and sql. sql can contain placeholders
// for bound parameters. These placeholders are referenced positional as $1, $2, etc.
//
// Prepare is idempotent; i.e. it is safe to call Prepare multiple times with the same
// name and sql arguments. This allows a code path to Prepare and Query/Exec without
// concern for if the statement has already been prepared.
func (c *Conn) Prepare(name, sql string) (ps *PreparedStatement, err error) {
	return c.PrepareEx(context.Background(), name, sql, nil)
}

// PrepareEx creates a prepared statement with name and sql. sql can contain placeholders
// for bound parameters. These placeholders are referenced positional as $1, $2, etc.
// It defers from Prepare as it allows additional options (such as parameter OIDs) to be passed via struct
//
// PrepareEx is idempotent; i.e. it is safe to call PrepareEx multiple times with the same
// name and sql arguments. This allows a code path to PrepareEx and Query/Exec without
// concern for if the statement has already been prepared.
func (c *Conn) PrepareEx(ctx context.Context, name, sql string, opts *PrepareExOptions) (ps *PreparedStatement, err error) {
	err = c.waitForPreviousCancelQuery(ctx)
	if err != nil {
		return nil, err
	}

	err = c.initContext(ctx)
	if err != nil {
		return nil, err
	}

	ps, err = c.prepareEx(name, sql, opts)
	err = c.termContext(err)
	return ps, err
}

// PreparedStatement is a description of a prepared statement
type PreparedStatement struct {
	Name              string
	SQL               string
	FieldDescriptions []FieldDescription
	ParameterOIDs     []types.OID
}

// PrepareExOptions is an option struct that can be passed to PrepareEx
type PrepareExOptions struct {
	ParameterOIDs []types.OID
}

func Connect(config driver.ConnConfig) (c *Conn, err error) {
	return connect(config, minimalConnInfo)
}

// Close closes a connection. It is safe to call Close on a already closed
// connection.
func (c *Conn) Close() (err error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.status < connStatusIdle {
		return nil
	}
	c.status = connStatusClosed

	defer func() {
		c.conn.Close()
		c.causeOfDeath = errors.New("Closed")
	}()

	err = c.conn.SetDeadline(time.Time{})
	if err != nil {
		golog.Warn("failed to clear deadlines to send close message. %s", err.Error())
		return err
	}

	_, err = c.conn.Write([]byte{'X', 0, 0, 0, 4})
	if err != nil {
		golog.Warn("failed to send terminate message. %s", err.Error())
		return err
	}

	err = c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		golog.Warn("failed to set read deadline to finish closing. %s", err.Error())
		return err
	}

	_, err = c.conn.Read(make([]byte, 1))
	if err != io.EOF {
		return err
	}

	return nil
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

func initPostgresql(c *Conn) (*types.ConnInfo, error) {
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

	cinfo := types.NewConnInfo()
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

func (c *Conn) ensureConnectionReadyForQuery() error {
	for c.pendingReadyForQueryCount > 0 {
		msg, err := c.rxMsg()
		if err != nil {
			return err
		}

		switch msg := msg.(type) {
		case *pgproto.ErrorResponse:
			pgErr := c.rxErrorResponse(msg)
			if pgErr.Severity == "FATAL" {
				return pgErr
			}
		default:
			err = c.processContextFreeMsg(msg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Conn) sendSimpleQuery(sql string, args ...interface{}) error {
	if err := c.ensureConnectionReadyForQuery(); err != nil {
		return err
	}

	if len(args) == 0 {
		buf := appendQuery(c.wbuf, sql)

		_, err := c.conn.Write(buf)
		if err != nil {
			c.die(err)
			return err
		}
		c.pendingReadyForQueryCount++

		return nil
	}

	ps, err := c.Prepare("", sql)
	if err != nil {
		return err
	}

	return c.sendPreparedQuery(ps, args...)
}

func (c *Conn) initContext(ctx context.Context) error {
	if c.ctxInProgress {
		return errors.New("ctx already in progress")
	}

	if ctx.Done() == nil {
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	c.ctxInProgress = true

	go c.contextHandler(ctx)

	return nil
}

func (c *Conn) termContext(opErr error) error {
	if !c.ctxInProgress {
		return opErr
	}

	var err error

	select {
	case err = <-c.closedChan:
		if opErr == nil {
			err = nil
		}
	case c.doneChan <- struct{}{}:
		err = opErr
	}

	c.ctxInProgress = false
	return err
}

func (c *Conn) sendPreparedQuery(ps *PreparedStatement, arguments ...interface{}) (err error) {
	if len(ps.ParameterOIDs) != len(arguments) {
		return errors.New("Prepared statement \"%v\" requires %d parameters, but %d were provided").Format(ps.Name, len(ps.ParameterOIDs), len(arguments))
	}

	resultFormatCodes := make([]int16, len(ps.FieldDescriptions))
	for i, fd := range ps.FieldDescriptions {
		resultFormatCodes[i] = fd.FormatCode
	}

	buf, err := appendBind(c.wbuf, "", ps.Name, c.ConnInfo, ps.ParameterOIDs, arguments, resultFormatCodes)
	if err != nil {
		return err
	}

	buf = appendExecute(buf, "", 0)
	buf = appendSync(buf)
	n, err := c.conn.Write(buf)
	if err != nil {
		if fatalWriteErr(n, err) {
			c.die(err)
		}
		return err
	}
	c.pendingReadyForQueryCount++
	return nil
}

// fatalWriteError takes the response of a net.Conn.Write and determines if it is fatal
func fatalWriteErr(bytesWritten int, err error) bool {
	// Partial writes break the connection
	if bytesWritten > 0 {
		return true
	}

	netErr, is := err.(net.Error)
	return !(is && netErr.Timeout())
}

func (c *Conn) contextHandler(ctx context.Context) {
	select {
	case <-ctx.Done():
		c.cancelQuery()
		c.closedChan <- ctx.Err()
	case <-c.doneChan:
	}
}

// cancelQuery sends a cancel request to the PostgreSQL server. It returns an
// error if unable to deliver the cancel request, but lack of an error does not
// ensure that the query was canceled. As specified in the documentation, there
// is no way to be sure a query was canceled. See
// https://www.postgresql.org/docs/current/static/protocol-flow.html#AEN112861
func (c *Conn) cancelQuery() {
	if !atomic.CompareAndSwapInt32(&c.cancelQueryInProgress, 0, 1) {
		panic("cancelQuery when cancelQueryInProgress")
	}

	if err := c.conn.SetDeadline(time.Now()); err != nil {
		c.Close() // Close connection if unable to set deadline
		return
	}

	doCancel := func() error {
		network, address := c.config.NetworkAddress()
		cancelConn, err := defaultDialer().Dial(network, address)
		if err != nil {
			return err
		}
		defer cancelConn.Close()

		// If server doesn't process cancellation request in bounded time then abort.
		err = cancelConn.SetDeadline(time.Now().Add(15 * time.Second))
		if err != nil {
			return err
		}

		buf := make([]byte, 16)
		binary.BigEndian.PutUint32(buf[0:4], 16)
		binary.BigEndian.PutUint32(buf[4:8], 80877102)
		binary.BigEndian.PutUint32(buf[8:12], uint32(c.pid))
		binary.BigEndian.PutUint32(buf[12:16], uint32(c.secretKey))
		_, err = cancelConn.Write(buf)
		if err != nil {
			return err
		}

		_, err = cancelConn.Read(buf)
		if err != io.EOF {
			return errors.New("Server failed to close connection after cancel query request: %v %v").Format(err, buf)
		}

		return nil
	}

	go func() {
		err := doCancel()
		if err != nil {
			c.Close() // Something is very wrong. Terminate the connection.
		}
		c.cancelQueryCompleted <- struct{}{}
	}()
}

func (c *Conn) waitForPreviousCancelQuery(ctx context.Context) error {
	if atomic.LoadInt32(&c.cancelQueryInProgress) == 0 {
		return nil
	}

	select {
	case <-c.cancelQueryCompleted:
		atomic.StoreInt32(&c.cancelQueryInProgress, 0)
		if err := c.conn.SetDeadline(time.Time{}); err != nil {
			c.Close() // Close connection if unable to disable deadline
			return err
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
