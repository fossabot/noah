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

package pgwire

import (
    "github.com/pkg/errors"
    "github.com/readystock/golog"
    "github.com/readystock/noah/db/sql"
    "github.com/readystock/noah/db/sql/oid"
    "github.com/readystock/noah/db/sql/pgwire/arguments"
    "github.com/readystock/noah/db/sql/pgwire/pgerror"
    "github.com/readystock/noah/db/sql/pgwire/pgwirebase"
    "github.com/readystock/pg_query_go"
    nodes "github.com/readystock/pg_query_go/nodes"
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
            // If the number of arguments so far is 0, we want to check with our own function
            // to double check.
            golog.Infof("found %d numQArgTypes in query", numQArgTypes)
            if numQArgTypes == 0 {
                args := arguments.GetArguments(stmt)
                if args > 0 {
                    golog.Infof("found %d arguments in query", args)
                }
            }

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
    //     if t == 0 {
    //         continue
    //     }
    //     v, ok := types.OidToType[t]
    //     if !ok {
    //         err := pgwirebase.NewProtocolViolationErrorf("unknown oid type: %v", t)
    //         return c.stmtBuf.Push(ctx, sql.SendError{Err: err})
    //     }
    //     sqlTypeHints[strconv.Itoa(i+1)] = v
    // }

}
