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

// func Test_Select_MultiColumn(t *testing.T) {
//     test := QueryTest{
//         Query: "SELECT true::bool, 5.5;",
//         Args: nil,
//         Expected: [][]interface{}{
//             {
//                 true,
//                 5.5,
//             },
//         },
//     }
//     DoQueryTest(t, test)
// }

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

func Test_Create_AndSelect(t *testing.T) {
    DoExecTest(t, ExecTest{
        Query:    "CREATE TABLE IF NOT EXISTS public.temptest123 (id BIGINT);",
        Expected: 0,
    })
    defer DoExecTest(t, ExecTest{
        Query:    "DROP TABLE public.temptest123 CASCADE;",
        Expected: 0,
    })
    DoQueryTest(t, QueryTest{
        Query: "INSERT INTO public.temptest123 (id) VALUES(1) RETURNING id;",
        Expected: [][]interface{}{
            {int64(1),},
        },
    })
}
