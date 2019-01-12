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

package queryutil

import (
    "github.com/readystock/pg_query_go/nodes"
    "reflect"
)

func GetColumns(stmt pg_query.Stmt) []pg_query.ResTarget {
    cols := make([]pg_query.ResTarget, 0)
    switch ast := stmt.(type) {
    case pg_query.SelectStmt:
        return examineColumns(ast.TargetList, 0)
    case pg_query.InsertStmt:
        return examineColumns(ast.ReturningList, 0)
    case pg_query.UpdateStmt:
        return examineColumns(ast.ReturningList, 0)
    case pg_query.DeleteStmt:
        return examineColumns(ast.ReturningList, 0)
    }
    return cols
}

func examineColumns(value interface{}, depth int) []pg_query.ResTarget {
    cols := make([]pg_query.ResTarget, 0)
    print := func(msg string, args ...interface{}) {
        // fmt.Printf("%s%s\n", strings.Repeat("\t", depth), fmt.Sprintf(msg, args...))
    }

    if value == nil {
        return cols
    }

    t := reflect.TypeOf(value)
    v := reflect.ValueOf(value)

    if v.Type() == reflect.TypeOf(pg_query.ResTarget{}) {
        col := value.(pg_query.ResTarget)
        cols = append(cols, col)
    }

    switch t.Kind() {
    case reflect.Ptr:
        if v.Elem().IsValid() {
            cols = append(cols, examineColumns(v.Elem().Interface(), depth+1)...)
        }
    case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
        depth--
        if v.Len() > 0 {
            print("[")
            for i := 0; i < v.Len(); i++ {
                depth++
                print("[%d] Type {%s} {", i, v.Index(i).Type().String())
                cols = append(cols, examineColumns(v.Index(i).Interface(), depth+1)...)
                print("},")
                depth--
            }
            print("]")
        } else {
            print("[]")
        }
    case reflect.Struct:
        for i := 0; i < t.NumField(); i++ {
            f := t.Field(i)
            print("[%d] Field {%s} Type {%s} Kind {%s}", i, f.Name, f.Type.String(), reflect.ValueOf(value).Field(i).Kind().String())
            cols = append(cols, examineColumns(reflect.ValueOf(value).Field(i).Interface(), depth+1)...)
        }
    }
    return cols
}
