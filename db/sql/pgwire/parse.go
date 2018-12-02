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

package pgwire

import (
	"github.com/Ready-Stock/Noah/db/sql"
	"github.com/Ready-Stock/Noah/db/sql/oid"
	"github.com/Ready-Stock/Noah/db/sql/pgwire/pgerror"
	"github.com/Ready-Stock/Noah/db/sql/pgwire/pgwirebase"
	"github.com/Ready-Stock/pg_query_go"
	nodes "github.com/Ready-Stock/pg_query_go/nodes"
	"github.com/kataras/golog"
	"github.com/pkg/errors"
	"reflect"
	"time"
)

// An error is returned iff the statement buffer has been closed. In that case,
// the connection should be considered toast.
func (c *conn) handleParse(buf *pgwirebase.ReadBuffer) error {
	// protocolErr is set if a protocol error has to be sent to the client. A
	// stanza at the bottom of the function pushes instructions for sending this
	// error.
	startParse := time.Now().UTC()
	var protocolErr *pgerror.Error

	name, err := buf.GetString()
	if protocolErr != nil {
		return c.stmtBuf.Push(sql.SendError{Err: err})
	}
	query, err := buf.GetString()
	if err != nil {
		return c.stmtBuf.Push(sql.SendError{Err: err})
	}
	// The client may provide type information for (some of) the placeholders.
	numQArgTypes, err := buf.GetUint16()
	if err != nil {
		return err
	}
	inTypeHints := make([]oid.Oid, numQArgTypes)
	for i := range inTypeHints {
		typ, err := buf.GetUint32()
		if err != nil {
			return c.stmtBuf.Push(sql.SendError{Err: err})
		}
		inTypeHints[i] = oid.Oid(typ)
	}
	golog.Infof("[%s] Query: `%s`", c.conn.RemoteAddr().String(), query)
	p, err := pg_query.Parse(query)
	endParse := time.Now().UTC()
	if err != nil {
		golog.Errorf("[%s] %s", c.conn.RemoteAddr().String(), err.Error())
		return c.stmtBuf.Push(sql.SendError{Err: err})
	}

	j, _ := p.MarshalJSON()
	// golog.Debugf("[%s] Tree: %s", c.conn.RemoteAddr().String(), string(j))
	if len(p.Statements) > 0 {
		if stmt, ok := p.Statements[0].(nodes.RawStmt).Stmt.(nodes.Stmt); !ok {
			return c.stmtBuf.Push(sql.SendError{Err: errors.Errorf("error, cannot currently handle statements of type: %s, json: %s", reflect.TypeOf(p.Statements[0].(nodes.RawStmt).Stmt).Name(), string(j))})
		} else {
			return c.stmtBuf.Push(sql.PrepareStmt{
				Name:         name,
				RawTypeHints: inTypeHints,
				ParseStart:   startParse,
				ParseEnd:     endParse,
				PGQuery:      stmt,
			})
		}
	} else {
		return c.stmtBuf.Push(sql.SendError{Err: errors.Errorf("error, no statement to execute")})
	}


	// Prepare the mapping of SQL placeholder names to types. Pre-populate it with
	// the type hints received from the client, if any.
	// sqlTypeHints := make(tree.PlaceholderTypes)
	// for i, t := range inTypeHints {
	// 	if t == 0 {
	// 		continue
	// 	}
	// 	v, ok := types.OidToType[t]
	// 	if !ok {
	// 		err := pgwirebase.NewProtocolViolationErrorf("unknown oid type: %v", t)
	// 		return c.stmtBuf.Push(ctx, sql.SendError{Err: err})
	// 	}
	// 	sqlTypeHints[strconv.Itoa(i+1)] = v
	// }

}
