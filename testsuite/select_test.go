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

package testsuite

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "testing"
)

func Test_Select_Simple(t *testing.T) {
    DoQueryTest(t, QueryTest{
        Query: "SELECT 1;",
        Expected: [][]interface{}{
            {int32(1),},
        },
    })
}

func Test_Select_TypeConversion(t *testing.T) {
    DoQueryTest(t, QueryTest{
        Query: "SELECT 1::text;",
        Expected: [][]interface{}{
            {"1",},
        },
    })
}

func Test_Select_CurrentUser(t *testing.T) {
    DoQueryTest(t, QueryTest{
        Query: "SELECT current_user;",
        Expected: [][]interface{}{
            {"postgres",},
        },
    })
}

func Test_Select_Numeric(t *testing.T) {
    t.Skip("numeric types are not supported by noah's client facing wire yet")
    DoQueryTest(t, QueryTest{
        Query: "SELECT 5.5;",
        Args:  nil,
        Expected: [][]interface{}{
            {5.5,},
        },
    })
}

func Test_Select_Unions(t *testing.T) {
    DoQueryTest(t, QueryTest{
        Query: "SELECT 1, 2, 3 UNION ALL SELECT 4, 5, 6 UNION ALL SELECT 7, 8, 9",
        Expected: [][]interface{}{
            {int32(1), int32(2), int32(3),},
            {int32(4), int32(5), int32(6),},
            {int32(7), int32(8), int32(9),},
        },
    })
}

func Test_Select_Math(t *testing.T) {
    DoQueryTest(t, QueryTest{
        Query: "SELECT 5-3, 1+1, 4/2, 10*2",
        Expected: [][]interface{}{
            {int32(2), int32(2), int32(2), int32(20),},
        },
    })
}

func Test_Create_And_Insert(t *testing.T) {
    DoExecTest(t, ExecTest{
        Query: "CREATE TABLE public.temp (id BIGINT);",
    })
    defer DoExecTest(t, ExecTest{
        Query: "DROP TABLE public.temp CASCADE;",
    })
    DoQueryTest(t, QueryTest{
        Query: "INSERT INTO public.temp (id) VALUES(1) RETURNING id;",
        Expected: [][]interface{}{
            {int64(1),},
        },
    })
}

func Test_Create_Global_And_Insert(t *testing.T) {
    DoExecTest(t, ExecTest{
        Query: `CREATE TABLE public.temp (id BIGINT) TABLESPACE "noah.global";`,
    })
    defer DoExecTest(t, ExecTest{
        Query: "DROP TABLE public.temp CASCADE;",
    })
    DoQueryTest(t, QueryTest{
        Query: "INSERT INTO public.temp (id) VALUES(1) RETURNING id;",
        Expected: [][]interface{}{
            {int64(1),},
        },
    })
}

func Test_Create_Global_And_Insert_Multiple(t *testing.T) {
    DoExecTest(t, ExecTest{
        Query: `CREATE TABLE public.temp (id BIGINT) TABLESPACE "noah.global";`,
    })
    defer DoExecTest(t, ExecTest{
        Query: "DROP TABLE public.temp CASCADE;",
    })

    for i := 0; i < 10; i ++ {
        DoQueryTest(t, QueryTest{
            Query: fmt.Sprintf("INSERT INTO public.temp (id) VALUES (%d) RETURNING id;", i),
            Expected: [][]interface{}{
                {int64(i),},
            },
        })
    }
}

func Test_Create_Global_And_InsertFew_ThenSelect(t *testing.T) {
    DoExecTest(t, ExecTest{
        Query: `CREATE TABLE IF NOT EXISTS public.temp (id BIGINT) TABLESPACE "noah.global";`,
    })
    // DoExecTest(t, ExecTest{
    //     Query: `DELETE FROM public.temp;`,
    // })
    defer DoExecTest(t, ExecTest{
        Query: "DROP TABLE public.temp CASCADE;",
    })

    for i := 0; i < 10; i ++ {
        DoQueryTest(t, QueryTest{
            Query: fmt.Sprintf("INSERT INTO public.temp (id) VALUES (%d) RETURNING id;", i),
            Expected: [][]interface{}{
                {int64(i),},
            },
        })
    }

    for i := 0; i < 10; i ++ {
        DoQueryTest(t, QueryTest{
            Query: fmt.Sprintf("SELECT id FROM public.temp WHERE id = %d;", i),
            Expected: [][]interface{}{
                {int64(i),},
            },
        })
    }
}

func Test_Create_And_InsertMultiRows(t *testing.T) {
    t.Skip("pg_query_go's deparser does not support deparsing multiple value rows yet")
    DoExecTest(t, ExecTest{
        Query: "CREATE TABLE IF NOT EXISTS public.temptest123 (id BIGINT);",
    })
    defer DoExecTest(t, ExecTest{
        Query: "DROP TABLE public.temptest123 CASCADE;",
    })
    DoQueryTest(t, QueryTest{
        Query: "INSERT INTO public.temptest123 (id) VALUES(1),(2) RETURNING id;",
        Expected: [][]interface{}{
            {int64(1),},
        },
    })
}

func Test_Create_And_InsertFromAmbiguousSelect(t *testing.T) {
    t.Skip("noah does not yet support inserting from an ambiguous select")

    // Create tables needed
    DoExecTest(t, ExecTest{
        Query: `CREATE TABLE IF NOT EXISTS public.tempaccounts (id BIGSERIAL PRIMARY KEY, number INT) TABLESPACE "noah.account";`,
    })
    DoExecTest(t, ExecTest{
        Query: `CREATE TABLE IF NOT EXISTS public.tempsharded (id BIGSERIAL PRIMARY KEY, account_id BIGINT NOT NULL REFERENCES public.tempaccounts (id)) TABLESPACE "noah.shard";`,
    })

    // Drop our tables at the end of the test
    defer DoExecTest(t, ExecTest{
        Query: "DROP TABLE public.tempaccounts CASCADE;",
    })
    defer DoExecTest(t, ExecTest{
        Query: "DROP TABLE IF EXISTS public.tempsharded CASCADE;",
    })
    //
    // //numberOfAccounts := 4
    //
    // account1 := DoQueryTest(t, QueryTest{
    //     Query: "INSERT INTO public.tempaccounts (number) VALUES(1) RETURNING *;",
    // })
    //
    // account2 := DoQueryTest(t, QueryTest{
    //     Query: "INSERT INTO public.tempaccounts (number) VALUES(2) RETURNING *;",
    // })
    //
    // _, _ := account1[0][0], account2[0][0]
    //
    //
}

func Test_Create_And_InsertWithSerial(t *testing.T) {
    DoExecTest(t, ExecTest{
        Query: "DROP TABLE IF EXISTS public.temp CASCADE;",
    })
    DoExecTest(t, ExecTest{
        Query: "CREATE TABLE public.temp (id BIGSERIAL, number BIGINT);",
    })
    defer DoExecTest(t, ExecTest{
        Query: "DROP TABLE public.temp CASCADE;",
    })

    for i := 0; i < 10; i++ {
        result := DoQueryTest(t, QueryTest{
            Query: fmt.Sprintf("INSERT INTO public.temp (number) VALUES(%d) RETURNING *;", i),
        })
        assert.Equal(t, int64(i), result[0][1], "the number column inserted does not match the number returned")
        assert.True(t, result[0][0].(int64) > 0, "the ID column returned is not valid")
    }
}