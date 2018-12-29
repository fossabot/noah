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
    "github.com/readystock/golinq"
    "github.com/readystock/pg_query_go/nodes"
    "reflect"
)

func GetArguments(stmt interface{}) (numArgs int) {
    return linq.From(examineArguments(stmt, 0)).Distinct().Count()
}

func examineArguments(value interface{}, depth int) []int {
    args := make([]int, 0)
    print := func(msg string, args ...interface{}) {
        // fmt.Printf("%s%s\n", strings.Repeat("\t", depth), fmt.Sprintf(msg, args...))
    }

    if value == nil {
        return args
    }

    t := reflect.TypeOf(value)
    v := reflect.ValueOf(value)

    if v.Type() == reflect.TypeOf(pg_query.ParamRef{}) {
        param := value.(pg_query.ParamRef)
        args = append(args, param.Number)
    }

    switch t.Kind() {
    case reflect.Ptr:
        if v.Elem().IsValid() {
            args = append(args, examineArguments(v.Elem().Interface(), depth+1)...)
        }
    case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
        depth--
        if v.Len() > 0 {
            print("[")
            for i := 0; i < v.Len(); i++ {
                depth++
                print("[%d] Type {%s} {", i, v.Index(i).Type().String())
                args = append(args, examineArguments(v.Index(i).Interface(), depth+1)...)
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
            args = append(args, examineArguments(reflect.ValueOf(value).Field(i).Interface(), depth+1)...)
        }
    }
    return args
}
