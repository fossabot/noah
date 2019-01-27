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
	"github.com/ahmetb/go-linq"
	"github.com/readystock/noah/db/sql/types"
	"github.com/readystock/pg_query_go/nodes"
	"reflect"
	"strconv"
	"strings"
)

type BuiltInFunction func(args ...pg_query.Node) (interface{}, error)
type BuiltInFunctionMap map[string]BuiltInFunction

func ReplaceFunctionCalls(stmt interface{}, builtIns BuiltInFunctionMap) (interface{}, error) {
	return replaceFunctions(stmt, 0, builtIns)
}

func replaceFunctions(value interface{}, depth int, builtIns BuiltInFunctionMap) (interface{}, error) {
	print := func(msg string, args ...interface{}) {
		// fmt.Printf("%s%s\n", strings.Repeat("\t", depth), fmt.Sprintf(msg, args...))
	}

	if funcCall, ok := value.(pg_query.FuncCall); ok {
		functionName, err := funcCall.Name()
		if err != nil {
			return funcCall, err
		}
		print("found function [%s] checking for drop-ins", functionName)
		if builtInFunction, ok := builtIns[functionName]; ok {
			result, err := builtInFunction(funcCall.Args.Items...)
			if err != nil {
				return funcCall, err
			}
			return convertObjectToQueryLiteral(result)
		}
	}

	if value == nil {
		return nil, nil
	}

	typ := reflect.TypeOf(value)
	val := reflect.ValueOf(value)

	print("[-] Parent Type <%s> Kind <%s>", typ.Name(), typ.Kind().String())
	depth++
	switch typ.Kind() {
	case reflect.Ptr:
		return value, nil
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
				result, err := replaceFunctions(item.Interface(), depth+1, builtIns)
				if err != nil {
					return nil, err
				}
				if result != nil {
					copySlice.Index(i).Set(reflect.ValueOf(result))
				}
			}
			depth--
			print("]")
			return copySlice.Interface(), nil
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
			result, err := replaceFunctions(actual.Interface(), depth+1, builtIns)
			if err != nil {
				return nil, err
			}
			if result != nil {
				copy.Field(i).Set(reflect.ValueOf(result))
			}
		}
		return copy.Interface(), nil
	default:
		return value, nil
	}
	return value, nil
}

func convertObjectToQueryLiteral(obj interface{}) (pg_query.Node, error) {
	switch argValue := obj.(type) {
	case string:
		return pg_query.A_Const{
			Val: pg_query.String{
				Str: argValue,
			},
		}, nil
	case []string:
		items := make([]string, len(argValue))
		linq.From(argValue).SelectT(func(val string) string {
			return fmt.Sprintf(`"%s"`, strings.Replace(val, `"`, `\"`, 0))
		}).ToSlice(&items)
		return pg_query.TypeCast{
			Arg: pg_query.A_Const{
				Val: pg_query.String{
					Str: fmt.Sprintf("{%s}", strings.Join(items, ", ")),
				},
			},
			TypeName: &pg_query.TypeName{
				Names: pg_query.List{
					Items: []pg_query.Node{
						pg_query.String{
							Str: "text",
						},
					},
				},
				ArrayBounds: pg_query.List{
					Items: []pg_query.Node{
						pg_query.Integer{
							Ival: -1,
						},
					},
				},
			},
		}, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		intv, _ := strconv.ParseUint(fmt.Sprintf("%v", argValue), 10, 64)
		return pg_query.A_Const{
			Val: pg_query.Integer{
				Ival: int64(intv),
			},
		}, nil
	case *types.Numeric:
		floatyMcFloatyFace := float64(0.0)
		err := argValue.AssignTo(&floatyMcFloatyFace)
		if err != nil {
			return nil, err
		}
		return pg_query.A_Const{
			Val: pg_query.Float{
				Str: fmt.Sprintf("%v", floatyMcFloatyFace),
			},
		}, nil
	case float32, float64:
		return pg_query.A_Const{
			Val: pg_query.Float{
				Str: fmt.Sprintf("%v", argValue),
			},
		}, nil
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
		}, nil
	case nil:
		return &pg_query.A_Const{
			Val: pg_query.Null{},
		}, nil
	default:
		panic(fmt.Sprintf("unsupported type %+v", argValue))
	}
}
