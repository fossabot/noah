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

package sql

import (
    "encoding/json"
    "fmt"
    "github.com/readystock/noah/db/system"
    "github.com/readystock/pg_query_go"
    pg_query2 "github.com/readystock/pg_query_go/nodes"
    "github.com/stretchr/testify/assert"
    "testing"
)

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
    assert.Equal(t,`CREATE TABLE "test" (id bigint, email text)`, plans[0].CompiledQuery,
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

    assert.Equal(t,`CREATE TABLE "test" (temp text, id bigint, email text, CONSTRAINT pk_test PRIMARY KEY ("id"))`, plans[0].CompiledQuery,
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

    assert.Equal(t,`CREATE TABLE "test" (id bigint PRIMARY KEY, email text)`, plans[0].CompiledQuery,
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

    assert.Equal(t,`CREATE TABLE "test" (id bigint, tinyid int, flake bigint)`, plans[0].CompiledQuery,
        "the resulting query plan did not equal the expected query plan, did something change with how queries were recompiled?")
}

func Test_Create_CompilePlan_Sharded_ReferencedForeignKey(t *testing.T) {
    accountSql := `CREATE TABLE accounts (account_id BIGSERIAL PRIMARY KEY, name TEXT) TABLESPACE "noah.account";`
    parsedAccount, err := pg_query.Parse(accountSql)
    if err != nil {
        panic(err)
    }

    accountStmt := CreateCreateStatement(parsedAccount.Statements[0].(pg_query2.RawStmt).Stmt.(pg_query2.CreateStmt))

    _, err = accountStmt.compilePlan(ConnExecutor, Nodes)
    if err != nil {
        panic(err)
    }

    if err := SystemCtx.Schema.CreateTable(accountStmt.table); err != nil {
        panic(err)
    }
    defer SystemCtx.Schema.DropTable(accountStmt.table.TableName)

    sql := `CREATE TABLE products (id BIGSERIAL PRIMARY KEY, account_id BIGINT NOT NULL REFERENCES accounts (account_id)) TABLESPACE "noah.shard";`
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

    assert.Equal(t,`CREATE TABLE "products" (id bigint PRIMARY KEY, account_id int8 NOT NULL FOREIGN KEY REFERENCES "accounts" ("account_id"))`, plans[0].CompiledQuery,
        "the resulting query plan did not equal the expected query plan, did something change with how queries were recompiled?")

    assert.NotNil(t, stmt.table.ShardKey, "table's shard key should not be null")

    assert.Equal(t, accountStmt.table.TableName, stmt.table.ShardKey.(*system.NTable_SKey).SKey.ForeignKey.(*system.NColumn_FKey).FKey.TableName, "the table name from the accounts table does not match the table name of the shard column foreign key")

    assert.Equal(t, accountStmt.table.PrimaryKey.(*system.NTable_PKey).PKey.ColumnName, stmt.table.ShardKey.(*system.NTable_SKey).SKey.ForeignKey.(*system.NColumn_FKey).FKey.ColumnName, "the column name of the referenced column does not match the accounts table's primary key column name")
}

func Test_Create_CompilePlan_Sharded_NamedForeignKey(t *testing.T) {
    accountSql := `CREATE TABLE accounts (account_id BIGSERIAL PRIMARY KEY, name TEXT) TABLESPACE "noah.account";`
    parsedAccount, err := pg_query.Parse(accountSql)
    if err != nil {
        panic(err)
    }

    accountStmt := CreateCreateStatement(parsedAccount.Statements[0].(pg_query2.RawStmt).Stmt.(pg_query2.CreateStmt))

    _, err = accountStmt.compilePlan(ConnExecutor, Nodes)
    if err != nil {
        panic(err)
    }

    if err := SystemCtx.Schema.CreateTable(accountStmt.table); err != nil {
        panic(err)
    }
    defer SystemCtx.Schema.DropTable(accountStmt.table.TableName)

    sql := `CREATE TABLE products (id BIGSERIAL PRIMARY KEY, account_id BIGINT NOT NULL, CONSTRAINT fk_products_account FOREIGN KEY (account_id) REFERENCES accounts (account_id)) TABLESPACE "noah.shard";`
    parsed, err := pg_query.Parse(sql)
    if err != nil {
        panic(err)
    }
    j, _ := json.Marshal(parsed)
    fmt.Println(string(j))

    stmt := CreateCreateStatement(parsed.Statements[0].(pg_query2.RawStmt).Stmt.(pg_query2.CreateStmt))

    plans, err := stmt.compilePlan(ConnExecutor, Nodes)
    if err != nil {
        panic(err)
    }

    assert.Equal(t, len(plans), len(Nodes),
        "the number of plans returned did not match the number of nodes that this query should target.")

    assert.Equal(t, `CREATE TABLE "products" (id bigint PRIMARY KEY, account_id int8 NOT NULL, CONSTRAINT fk_products_account FOREIGN KEY ("account_id") REFERENCES "accounts" ("account_id"))`, plans[0].CompiledQuery,
        "the resulting query plan did not equal the expected query plan, did something change with how queries were recompiled?")

    assert.NotNil(t, stmt.table.ShardKey, "table's shard key should not be null")
}