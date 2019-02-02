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
	"fmt"
	"github.com/readystock/golinq"
	"github.com/readystock/noah/db/sql/plan"
	"github.com/readystock/noah/db/sql/types"
	"github.com/readystock/pg_query_go/nodes"
	"reflect"
	"strconv"
)

func GetArguments(stmt interface{}) []int {
	args := make([]int, 0)
	args = append(args, examineArguments(stmt, 0)...)
	linq.From(args).Distinct().ToSlice(&args)
	return args
}

func ReplaceArguments(stmt interface{}, args plan.QueryArguments) interface{} {
	return replaceArguments(stmt, 0, args)
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

func replaceArguments(value interface{}, depth int, args plan.QueryArguments) interface{} {
	print := func(msg string, args ...interface{}) {
		// fmt.Printf("%s%s\n", strings.Repeat("\t", depth), fmt.Sprintf(msg, args...))
	}

	if param, ok := value.(pg_query.ParamRef); ok {
		if arg, ok := args[strconv.FormatInt(int64(param.Number), 10)]; ok {
			return func() pg_query.Node {
				switch argValue := arg.Get().(type) {
				case string:
					return pg_query.A_Const{
						Val: pg_query.String{
							Str: argValue,
						},
					}
				case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
					intv, _ := strconv.ParseUint(fmt.Sprintf("%v", argValue), 10, 64)
					return pg_query.A_Const{
						Val: pg_query.Integer{
							Ival: int64(intv),
						},
					}
				case *types.Numeric:
					floatyMcFloatyFace := float64(0.0)
					argValue.AssignTo(&floatyMcFloatyFace)
					return pg_query.A_Const{
						Val: pg_query.Float{
							Str: fmt.Sprintf("%v", floatyMcFloatyFace),
						},
					}
				case float32, float64:
					return pg_query.A_Const{
						Val: pg_query.Float{
							Str: fmt.Sprintf("%v", argValue),
						},
					}
				case bool:
					boolVal := string([]rune(fmt.Sprintf("%v", argValue))[0])
					return pg_query.TypeCast{
						Arg: pg_query.A_Const{
							Val: pg_query.String{
								Str: boolVal,
							},
						},
						TypeName: &pg_query.TypeName{
							Names: pg_query.List{
								Items: []pg_query.Node{
									pg_query.String{
										Str: "pg_catalog",
									},
									pg_query.String{
										Str: "bool",
									},
								},
							},
						},
					}
				case nil:
					return &pg_query.A_Const{
						Val: pg_query.Null{},
					}
				default:
					panic(fmt.Sprintf("unsupported type %+v", argValue))
				}
			}()
		} else {
			panic("parameter is not a param reference")
		}

	}

	if value == nil {
		return nil
	}

	typ := reflect.TypeOf(value)
	val := reflect.ValueOf(value)

	print("[-] Parent Type <%s> Kind <%s>", typ.Name(), typ.Kind().String())
	depth++
	switch typ.Kind() {
	case reflect.Ptr:
		return value
	case reflect.Slice:
		if val.Len() > 0 {
			print("[-] Slice Type <%s> Size: %d", val.Type().String(), val.Len())
			copySlice := reflect.MakeSlice(reflect.SliceOf(typ.Elem()), val.Len(), (val.Cap()+1)*2)
			reflect.Copy(copySlice, val)
			print("[")
			depth++
			for i := 0; i < val.Len(); i++ {
				item := val.Index(i)
				copy := reflect.New(item.Type())
				print("[%d] Copying Item Type <%s> Kind <%s> Actual %s", i, copy.Type().String(), copy.Kind().String(), item.Type().String())
				result := replaceArguments(item.Interface(), depth+1, args)
				if result != nil {
					copySlice.Index(i).Set(reflect.ValueOf(result))
				}
			}
			depth--
			print("]")
			return copySlice.Interface()
		} else {
			print("[-] Slice Type <%s> Size: Empty", val.Type().Name())
		}
	case reflect.Struct:
		copy := reflect.New(typ).Elem()
		copyType := reflect.TypeOf(copy.Interface())
		print("[-] Copied Struct Type <%s> Kind <%s> Fields: %d", copy.Type().Name(), copy.Kind().String(), copy.NumField())
		depth++
		for i := 0; i < copy.NumField(); i++ {
			field := copy.Field(i)
			fieldType := copyType.Field(i)
			actual := val.Field(i)
			print("[%d] Copying Field <%s> Type <%s> Kind <%s> Actual %s", i, fieldType.Name, field.Type().String(), field.Kind().String(), actual.Type().String())
			result := replaceArguments(actual.Interface(), depth+1, args)
			if result != nil {
				copy.Field(i).Set(reflect.ValueOf(result))
			}
		}
		return copy.Interface()
	default:
		return value
	}
	return value
}
