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
 *
 * Project: pq github.com/lib/pq
 * Copyright 2018  'pq' Contributors Portions Copyright (C) 2018 Blake Mizerany
 * License https://github.com/lib/pq/blob/master/LICENSE.md
 *
 * Project: go-linq github.com/ahmetb/go-linq
 * Copyright 2018 Ahmet Alp Balkan
 * License https://github.com/ahmetb/go-linq/blob/master/LICENSE
 */

package sql

import (
    "github.com/magiconair/properties/assert"
    "github.com/readystock/noah/db/system"
    "github.com/readystock/noah/testutils"
    "github.com/readystock/pg_query_go"
    pg_query2 "github.com/readystock/pg_query_go/nodes"
    "os"
    "testing"
)

var (
    SystemCtx    *system.SContext
    ConnExecutor *connExecutor
    Nodes        = []system.NNode{
        {
            NodeId:    1,
            Address:   "127.0.0.1:0",
            Port:      5432,
            Database:  "postgres",
            User:      "postgres",
            Password:  "",
            ReplicaOf: 0,
            Region:    "",
            Zone:      "",
            IsAlive:   true,
        },
    }
)

func TestMain(m *testing.M) {
    tempFolder := testutils.CreateTempFolder()
    defer testutils.DeleteTempFolder(tempFolder)

    SystemCtx, err := system.NewSystemContext(tempFolder, "127.0.0.1:0", "", "")
    if err != nil {
        panic(err)
    }
    defer SystemCtx.Close()

    for _, node := range Nodes {
        if err := SystemCtx.Nodes.AddNode(node); err != nil {
            panic(err)
        }
    }

    ConnExecutor = CreateConnExecutor(SystemCtx)

    retCode := m.Run()
    os.Exit(retCode)
}

func Test_Create_GetTargetNodes(t *testing.T) {
    sql := `CREATE TABLE test (id bigserial, email text);`
    parsed, err := pg_query.Parse(sql)
    if err != nil {
        panic(err)
    }

    stmt := CreateCreateStatement(parsed.Statements[0].(pg_query2.RawStmt).Stmt.(pg_query2.CreateStmt))

    nodes, err := stmt.getTargetNodes(ConnExecutor)
    if err != nil {
        panic(err)
    }

    // For create statements we want to run the query on ALL nodes that we can, so if the number of
    // nodes that are passed to compile query does not match the number of plans returned then that
    // means a node is missing or the plans were not generated completely.
    assert.Equal(t, len(nodes), len(nodes),
        "the number of nodes returned did not match the number of nodes that this query should target.")
}

func Test_Create_CompilePlan_Default(t *testing.T) {
    sql := `CREATE TABLE test (id bigserial, email text);`
    parsed, err := pg_query.Parse(sql)
    if err != nil {
        panic(err)
    }

    stmt := CreateCreateStatement(parsed.Statements[0].(pg_query2.RawStmt).Stmt.(pg_query2.CreateStmt))

    plans, err := stmt.compilePlan(ConnExecutor, Nodes)
    if err != nil {
        panic(err)
    }

    // For create statements we want to run the query on ALL nodes that we can, so if the number of
    // nodes that are passed to compile query does not match the number of plans returned then that
    // means a node is missing or the plans were not generated completely.
    assert.Equal(t, len(plans), len(Nodes),
        "the number of plans returned did not match the number of nodes that this query should target.")

    // This is a simple rewrite, we want to make sure that bigserial is being changed to bigint when
    // it is found in a create statement. We also want to make sure that the text matches a deparsed
    // query from pg_query_go.
    assert.Equal(t, plans[0].CompiledQuery, `CREATE TABLE "test" (id bigint, email text)`,
        "the resulting query plan did not equal the expected query plan, did something change with how queries were recompiled?")
}

func Test_Create_CompilePlan_AccountNoPrimaryKey(t *testing.T) {
    sql := `CREATE TABLE test (id bigserial, email text) TABLESPACE "noah.account";`
    parsed, err := pg_query.Parse(sql)
    if err != nil {
        panic(err)
    }

    stmt := CreateCreateStatement(parsed.Statements[0].(pg_query2.RawStmt).Stmt.(pg_query2.CreateStmt))

    _, err = stmt.compilePlan(ConnExecutor, Nodes)
    if err == nil {
        panic("plan should have failed to compile, account tables require a primary key.")
    }
}

func Test_Create_CompilePlan_AccountMultiColumnNamedPrimaryKey(t *testing.T) {
    sql := `CREATE TABLE test (temp text, id bigserial, email text, CONSTRAINT pk_test PRIMARY KEY (temp, id)) TABLESPACE "noah.account";`
    parsed, err := pg_query.Parse(sql)
    if err != nil {
        panic(err)
    }

    stmt := CreateCreateStatement(parsed.Statements[0].(pg_query2.RawStmt).Stmt.(pg_query2.CreateStmt))

    _, err = stmt.compilePlan(ConnExecutor, Nodes)
    if err == nil {
        panic("multi column primary keys should produce an error")
    }
}

func Test_Create_CompilePlan_AccountMissingNamedPrimaryKey(t *testing.T) {
    sql := `CREATE TABLE test (temp text, id bigserial, email text, CONSTRAINT pk_test PRIMARY KEY (account_id)) TABLESPACE "noah.account";`
    parsed, err := pg_query.Parse(sql)
    if err != nil {
        panic(err)
    }

    stmt := CreateCreateStatement(parsed.Statements[0].(pg_query2.RawStmt).Stmt.(pg_query2.CreateStmt))

    _, err = stmt.compilePlan(ConnExecutor, Nodes)
    if err == nil {
        panic("named primary key columns that do not exist should produce an error")
    }
}

func Test_Create_CompilePlan_AccountNamedPrimaryKey(t *testing.T) {
    sql := `CREATE TABLE test (temp text, id bigserial, email text, CONSTRAINT pk_test PRIMARY KEY (id)) TABLESPACE "noah.account";`
    parsed, err := pg_query.Parse(sql)
    if err != nil {
        panic(err)
    }

    stmt := CreateCreateStatement(parsed.Statements[0].(pg_query2.RawStmt).Stmt.(pg_query2.CreateStmt))

    plans, err := stmt.compilePlan(ConnExecutor, Nodes)
    if err != nil {
        panic(err)
    }

    assert.Equal(t, len(plans), len(Nodes),
        "the number of plans returned did not match the number of nodes that this query should target.")

    assert.Equal(t, plans[0].CompiledQuery, `CREATE TABLE "test" (temp text, id bigint, email text, CONSTRAINT pk_test PRIMARY KEY ("id"))`,
        "the resulting query plan did not equal the expected query plan, did something change with how queries were recompiled?")
}

func Test_Create_CompilePlan_AccountUUIDPrimaryKey(t *testing.T) {
    sql := `CREATE TABLE test (id uuid PRIMARY KEY, email text) TABLESPACE "noah.account";`
    parsed, err := pg_query.Parse(sql)
    if err != nil {
        panic(err)
    }

    stmt := CreateCreateStatement(parsed.Statements[0].(pg_query2.RawStmt).Stmt.(pg_query2.CreateStmt))

    _, err = stmt.compilePlan(ConnExecutor, Nodes)
    if err == nil {
        panic("plan should have failed to compile, tables must have a numeric primary key.")
    }
}

func Test_Create_CompilePlan_Account(t *testing.T) {
    sql := `CREATE TABLE test (id bigserial PRIMARY KEY, email text) TABLESPACE "noah.account";`
    parsed, err := pg_query.Parse(sql)
    if err != nil {
        panic(err)
    }

    stmt := CreateCreateStatement(parsed.Statements[0].(pg_query2.RawStmt).Stmt.(pg_query2.CreateStmt))

    plans, err := stmt.compilePlan(ConnExecutor, Nodes)
    if err != nil {
        panic(err)
    }

    assert.Equal(t, len(plans), len(Nodes),
        "the number of plans returned did not match the number of nodes that this query should target.")

    assert.Equal(t, plans[0].CompiledQuery, `CREATE TABLE "test" (id bigint PRIMARY KEY, email text)`,
        "the resulting query plan did not equal the expected query plan, did something change with how queries were recompiled?")
}

func Test_Create_CompilePlan_MultiplePrimaryKeys(t *testing.T) {
    sql := `CREATE TABLE test (id bigserial PRIMARY KEY, email text PRIMARY KEY);`
    parsed, err := pg_query.Parse(sql)
    if err != nil {
        panic(err)
    }

    stmt := CreateCreateStatement(parsed.Statements[0].(pg_query2.RawStmt).Stmt.(pg_query2.CreateStmt))

    _, err = stmt.compilePlan(ConnExecutor, Nodes)
    if err == nil {
        panic("the compiler should not allow you to create a table with multiple primary keys")
    }
}

func Test_Create_CompilePlan_ReplacementTypes(t *testing.T) {
    sql := `CREATE TABLE test (id bigserial, tinyid serial, flake snowflake);`
    parsed, err := pg_query.Parse(sql)
    if err != nil {
        panic(err)
    }

    stmt := CreateCreateStatement(parsed.Statements[0].(pg_query2.RawStmt).Stmt.(pg_query2.CreateStmt))

    plans, err := stmt.compilePlan(ConnExecutor, Nodes)
    if err != nil {
        panic(err)
    }

    assert.Equal(t, len(plans), len(Nodes),
        "the number of plans returned did not match the number of nodes that this query should target.")

    assert.Equal(t, plans[0].CompiledQuery, `CREATE TABLE "test" (id bigint, tinyid int, flake bigint)`,
        "the resulting query plan did not equal the expected query plan, did something change with how queries were recompiled?")
}
