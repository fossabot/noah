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
	"database/sql"
	"fmt"
	"github.com/Ready-Stock/Noah/db/sql/pgwire/pgproto"
	"github.com/Ready-Stock/Noah/db/sql/types"
	"github.com/pkg/errors"
	"reflect"
	"time"
)

type Row Rows

type Rows struct {
	conn       *Conn
	connPool   *ConnPool
	values     [][]byte
	fields     []FieldDescription
	rowCount   int
	columnIdx  int
	err        error
	startTime  time.Time
	sql        string
	args       []interface{}
	unlockConn bool
	closed     bool
}

func (rows *Rows) FieldDescriptions() []FieldDescription {
	return rows.fields
}

// Close closes the rows, making the connection ready for use again. It is safe
// to call Close after rows is already closed.
func (rows *Rows) Close() {
	if rows.closed {
		return
	}

	if rows.unlockConn {
		rows.conn.unlock()
		rows.unlockConn = false
	}

	rows.closed = true

	rows.err = rows.conn.termContext(rows.err)

	if rows.connPool != nil {
		rows.connPool.Release(rows.conn)
	}
}

func (rows *Rows) Err() error {
	return rows.err
}

// fatal signals an error occurred after the query was sent to the server. It
// closes the rows automatically.
func (rows *Rows) fatal(err error) {
	if rows.err != nil {
		return
	}

	rows.err = err
	rows.Close()
}

// Next prepares the next row for reading. It returns true if there is another
// row and false if no more rows are available. It automatically closes rows
// when all rows are read.
func (rows *Rows) Next() bool {
	if rows.closed {
		return false
	}

	rows.rowCount++
	rows.columnIdx = 0

	for {
		msg, err := rows.conn.rxMsg()
		if err != nil {
			rows.fatal(err)
			return false
		}

		switch msg := msg.(type) {
		case *pgproto.RowDescription:
			rows.fields = rows.conn.rxRowDescription(msg)
			for i := range rows.fields {
				if dt, ok := rows.conn.ConnInfo.DataTypeForOID(rows.fields[i].DataType); ok {
					rows.fields[i].DataTypeName = dt.Name
					rows.fields[i].FormatCode = TextFormatCode
				} else {
					rows.fields[i].DataTypeName = "text"
					rows.fields[i].FormatCode = TextFormatCode
					//rows.fatal(errors.Errorf("unknown oid: %d", rows.fields[i].DataType))
					//return false
				}
			}
		case *pgproto.DataRow:
			if len(msg.Values) != len(rows.fields) {
				rows.fatal(ProtocolError(fmt.Sprintf("Row description field count (%v) and data row field count (%v) do not match", len(rows.fields), len(msg.Values))))
				return false
			}

			rows.values = msg.Values
			return true
		case *pgproto.CommandComplete:
			rows.Close()
			return false

		default:
			err = rows.conn.processContextFreeMsg(msg)
			if err != nil {
				rows.fatal(err)
				return false
			}
		}
	}
}

func (rows *Rows) nextColumn() ([]byte, *FieldDescription, bool) {
	if rows.closed {
		return nil, nil, false
	}
	if len(rows.fields) <= rows.columnIdx {
		rows.fatal(ProtocolError("No next column available"))
		return nil, nil, false
	}

	buf := rows.values[rows.columnIdx]
	fd := &rows.fields[rows.columnIdx]
	rows.columnIdx++
	return buf, fd, true
}

type scanArgError struct {
	col int
	err error
}

func (e scanArgError) Error() string {
	return fmt.Sprintf("can't scan into dest[%d]: %v", e.col, e.err)
}

// Scan works the same as (*Rows Scan) with the following exceptions. If no
// rows were found it returns ErrNoRows. If multiple rows are returned it
// ignores all but the first.
func (r *Row) Scan(dest ...interface{}) (err error) {
	rows := (*Rows)(r)

	if rows.Err() != nil {
		return rows.Err()
	}

	if !rows.Next() {
		if rows.Err() == nil {
			return ErrNoRows
		}
		return rows.Err()
	}

	rows.Scan(dest...)
	rows.Close()
	return rows.Err()
}

// Scan reads the values from the current row into dest values positionally.
// dest can include pointers to core types, values implementing the Scanner
// interface, []byte, and nil. []byte will skip the decoding process and directly
// copy the raw bytes received from PostgreSQL. nil will skip the value entirely.
func (rows *Rows) Scan(dest ...interface{}) (err error) {
	if len(rows.fields) != len(dest) {
		err = errors.Errorf("Scan received wrong number of arguments, got %d but expected %d", len(dest), len(rows.fields))
		rows.fatal(err)
		return err
	}

	for i, d := range dest {
		buf, fd, _ := rows.nextColumn()

		if d == nil {
			continue
		}

		if s, ok := d.(types.BinaryDecoder); ok && fd.FormatCode == BinaryFormatCode {
			err = s.DecodeBinary(rows.conn.ConnInfo, buf)
			if err != nil {
				rows.fatal(scanArgError{col: i, err: err})
			}
		} else if s, ok := d.(types.TextDecoder); ok && fd.FormatCode == TextFormatCode {
			err = s.DecodeText(rows.conn.ConnInfo, buf)
			if err != nil {
				rows.fatal(scanArgError{col: i, err: err})
			}
		} else {
			if dt, ok := rows.conn.ConnInfo.DataTypeForOID(fd.DataType); ok {
				value := dt.Value
				switch fd.FormatCode {
				case TextFormatCode:
					if textDecoder, ok := value.(types.TextDecoder); ok {
						err = textDecoder.DecodeText(rows.conn.ConnInfo, buf)
						if err != nil {
							rows.fatal(scanArgError{col: i, err: err})
						}
					} else {
						rows.fatal(scanArgError{col: i, err: errors.Errorf("%T is not a pgtype.TextDecoder", value)})
					}
				case BinaryFormatCode:
					if binaryDecoder, ok := value.(types.BinaryDecoder); ok {
						err = binaryDecoder.DecodeBinary(rows.conn.ConnInfo, buf)
						if err != nil {
							rows.fatal(scanArgError{col: i, err: err})
						}
					} else {
						rows.fatal(scanArgError{col: i, err: errors.Errorf("%T is not a pgtype.BinaryDecoder", value)})
					}
				default:
					rows.fatal(scanArgError{col: i, err: errors.Errorf("unknown format code: %v", fd.FormatCode)})
				}

				if rows.Err() == nil {
					if scanner, ok := d.(sql.Scanner); ok {
						sqlSrc, err := types.DatabaseSQLValue(rows.conn.ConnInfo, value)
						if err != nil {
							rows.fatal(err)
						}
						err = scanner.Scan(sqlSrc)
						if err != nil {
							rows.fatal(scanArgError{col: i, err: err})
						}
					} else if err := value.AssignTo(d); err != nil {
						rows.fatal(scanArgError{col: i, err: err})
					}
				}
			} else {
				rows.fatal(scanArgError{col: i, err: errors.Errorf("unknown oid: %v", fd.DataType)})
			}
		}

		if rows.Err() != nil {
			return rows.Err()
		}
	}

	return nil
}

func (c *Conn) QueryEx(ctx context.Context, sql string, options *QueryExOptions) (rows *Rows, err error) {
	c.lastActivityTime = time.Now()
	rows = c.getRows(sql)

	err = c.waitForPreviousCancelQuery(ctx)
	if err != nil {
		rows.fatal(err)
		return rows, err
	}

	if err := c.ensureConnectionReadyForQuery(); err != nil {
		rows.fatal(err)
		return rows, err
	}

	if err := c.lock(); err != nil {
		rows.fatal(err)
		return rows, err
	}
	rows.unlockConn = true

	err = c.initContext(ctx)
	if err != nil {
		rows.fatal(err)
		return rows, rows.err
	}

	err = c.sendSimpleQuery(sql)
	if err != nil {
		rows.fatal(err)
		return rows, err
	}
	return rows, nil
}

// Values returns an array of the row values
func (rows *Rows) Values() ([]interface{}, error) {
	if rows.closed {
		return nil, errors.New("rows is closed")
	}

	values := make([]interface{}, 0, len(rows.fields))

	for range rows.fields {
		buf, fd, _ := rows.nextColumn()

		if buf == nil {
			values = append(values, nil)
			continue
		}

		if dt, ok := rows.conn.ConnInfo.DataTypeForOID(fd.DataType); ok {
			value := reflect.New(reflect.ValueOf(dt.Value).Elem().Type()).Interface().(types.Value)
			switch fd.FormatCode {
			case TextFormatCode:
				decoder := value.(types.TextDecoder)
				if decoder == nil {
					decoder = &types.GenericText{}
				}
				err := decoder.DecodeText(rows.conn.ConnInfo, buf)
				if err != nil {
					rows.fatal(err)
				}
				values = append(values, decoder.(types.Value).Get())
			case BinaryFormatCode:
				decoder := value.(types.BinaryDecoder)
				if decoder == nil {
					decoder = &types.GenericBinary{}
				}
				err := decoder.DecodeBinary(rows.conn.ConnInfo, buf)
				if err != nil {
					rows.fatal(err)
				}
				values = append(values, value.Get())
			default:
				rows.fatal(errors.New("Unknown format code"))
			}
		} else {
			value := types.Text{}
			switch fd.FormatCode {
			case TextFormatCode:
				err := value.DecodeText(rows.conn.ConnInfo, buf)
				if err != nil {
					rows.fatal(err)
				}
				values = append(values, value.Get())
			case BinaryFormatCode:
				err := value.DecodeBinary(rows.conn.ConnInfo, buf)
				if err != nil {
					rows.fatal(err)
				}
				values = append(values, value.Get())
			default:
				rows.fatal(errors.New("Unknown format code"))
			}
			//rows.fatal(errors.New("Unknown type"))
		}

		if rows.Err() != nil {
			return nil, rows.Err()
		}
	}

	return values, rows.Err()
}

func (c *Conn) getRows(sql string) *Rows {
	if len(c.preallocatedRows) == 0 {
		c.preallocatedRows = make([]Rows, 64)
	}

	r := &c.preallocatedRows[len(c.preallocatedRows)-1]
	c.preallocatedRows = c.preallocatedRows[0 : len(c.preallocatedRows)-1]

	r.conn = c
	r.startTime = c.lastActivityTime
	r.sql = sql
	r.args = make([]interface{}, 0)

	return r
}
