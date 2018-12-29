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

package queryutil

import (
    "github.com/readystock/noah/db/sql/plan"
    "github.com/readystock/noah/db/sql/types"
    "github.com/readystock/pg_query_go"
    pg_query2 "github.com/readystock/pg_query_go/nodes"
    "github.com/stretchr/testify/assert"
    "testing"
)

var (
    testQueries = []struct {
        Query    string
        ArgCount int
    }{
        {
            Query:    "SELECT $1::text;",
            ArgCount: 1,
        },
        {
            Query:    "SELECT e.typdelim FROM pg_catalog.pg_type t, pg_catalog.pg_type e WHERE t.oid = $1 and t.typelem = e.oid",
            ArgCount: 1,
        },
        {
            Query:    "SELECT e.typdelim FROM pg_catalog.pg_type t, pg_catalog.pg_type e WHERE t.oid = $1 and t.typelem = e.oid AND $2=$3",
            ArgCount: 3,
        },
        {
            Query:    "SELECT e.typdelim FROM pg_catalog.pg_type t, pg_catalog.pg_type e WHERE t.oid = $1 and t.typelem = e.oid AND $2=$1",
            ArgCount: 2,
        },
    }
)

func Test_GetArguments(t *testing.T) {
    for _, item := range testQueries {
        parsed, err := pg_query.Parse(item.Query)
        if err != nil {
            t.Error(err)
            t.FailNow()
        }

        stmt := parsed.Statements[0].(pg_query2.RawStmt).Stmt

        argCount := GetArguments(stmt)

        assert.Equal(t, item.ArgCount, argCount, "number of arguments does not match expected")
    }
}

var (
    testReplacements = []struct {
        Query     string
        ArgCount  int
        Arguments plan.QueryArguments
    }{
        {
            Query:    "SELECT $1",
            ArgCount: 1,
            Arguments: map[string]types.Value{
                "1": &types.Int4{
                    Status: types.Present,
                    Int:    1,
                },
            },
        },
    }
)

func Test_ReplaceArguments(t *testing.T) {
    for _, item := range testReplacements {
        parsed, err := pg_query.Parse(item.Query)
        if err != nil {
            t.Error(err)
            t.FailNow()
        }

        stmt := parsed.Statements[0].(pg_query2.RawStmt).Stmt

        argCount := GetArguments(stmt)

        assert.Equal(t, item.ArgCount, len(argCount), "number of arguments does not match expected")

        // Now we will replace the arguments, and there should be 0 after
        ReplaceArguments(stmt, item.Arguments)

        argCount = GetArguments(stmt)
        assert.Equal(t, 0, argCount, "number of arguments should now be 0")
    }
}
