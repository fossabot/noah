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
    "fmt"
    "github.com/readystock/golinq"
    "github.com/readystock/noah/db/sql/plan"
    "github.com/readystock/pg_query_go/nodes"
    "reflect"
    "strings"
)

func GetArguments(stmt interface{}) []int {
    args := make([]int, 0)
    linq.From(examineArguments(stmt, 0)).Distinct().ToSlice(&args)
    return args
}

func ReplaceArguments(stmt *pg_query.Node, args plan.QueryArguments) {
    replaceArguments(stmt, 0, args)
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

func replaceArguments(value *pg_query.Node, depth int, args plan.QueryArguments) {
    print := func(msg string, args ...interface{}) {
        fmt.Printf("%s%s\n", strings.Repeat("\t", depth), fmt.Sprintf(msg, args...))
    }

    processArray := func(values *[]pg_query.Node, depth int, args plan.QueryArguments) {
        vals := *values
        for i := 0; i < len(vals); i++ {
            print("[")
            val := vals[i]
            replaceArguments(&val, depth, args)
            vals[i] = val
            print("]")
        }
        *values = vals
    }

    if value == nil || depth > 10 {
        return
    }

    switch ntype := (*value).(type) {
    case pg_query.ParamRef:
        print("[-] [%d] Replacing param ref", depth)
        *value = &pg_query.String{
            Str: "MY HANDS ARE TYPING WORDS",
        }
    default:
        t := reflect.TypeOf(ntype)
        v := reflect.ValueOf(ntype)
        switch v.Kind() {
        case reflect.Ptr:
            if v.IsNil() {
                return
            }
            print("[-] Pointer Type {%s} Kind {%s}", v.Elem().Type().String(), v.Elem().Kind().String())
        case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
            depth--
            if v.Len() > 0 {
                print("[")
                for i := 0; i < v.Len(); i++ {
                    depth++
                    print("[%d] Type {%s} {", i, v.Index(i).Type().String())
                    //replaceArguments(v.Index(i).Addr().Interface(), depth+1, args);
                    print("},")
                    depth--
                }
                print("]")
            } else {
                print("[]")
            }
        case reflect.Struct:
            numField := v.NumField()
            for i := 0; i < numField; i++ {
                f := t.Field(i)
                print("[%d] [%d] Field {%s} Type {%s} Kind {%s} Actual {%s}", i, depth, f.Name, f.Type.String(), f.Type.Kind().String(), v.Field(i).Kind().String())
                switch v.Field(i).Type().Kind() {
                case reflect.Interface:
                    if v.Field(i).IsNil() {
                        continue
                    }
                    val, ok := v.Field(i).Interface().(pg_query.Node)
                    if !ok {
                        continue
                    }
                    replaceArguments(&val, depth+1, args)
                    reflect.ValueOf(value).Field(i).Addr().Set(reflect.ValueOf(val))
                    // if v.Field(i).CanSet() {
                    //     v.Field(i).Set(reflect.ValueOf(val))
                    // }
                case reflect.Slice:
                    if v.Field(i).Len() > 0 {
                        values := v.Field(i).Interface().([]pg_query.Node)
                        processArray(&values, depth+1, args)
                        if v.Field(i).CanSet() {
                            v.Field(i).Set(reflect.ValueOf(values))
                        }
                    }
                default:
                    val, ok := v.Field(i).Interface().(pg_query.Node)
                    if !ok {
                        continue
                    }
                    replaceArguments(&val, depth+1, args)
                }

                //replaceArguments(l.Field(i).Addr().Interface(), depth+1, args)
            }
        case reflect.Interface:
            print("snjaa")
        }
    }
}
