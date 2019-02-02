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
	"github.com/readystock/pg_query_go/nodes"
	"reflect"
	"strings"
)

func GetTables(stmt interface{}) []string {
	tables := make([]string, 0)
	linq.From(examineTables(stmt, 0)).Distinct().ToSlice(&tables)
	return tables
}

func examineTables(value interface{}, depth int) []string {
	args := make([]string, 0)
	print := func(msg string, args ...interface{}) {
		// fmt.Printf("%s%s\n", strings.Repeat("\t", depth), fmt.Sprintf(msg, args...))
	}

	if value == nil {
		return args
	}

	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)

	if v.Type() == reflect.TypeOf(pg_query.RangeVar{}) {
		rangeVar := value.(pg_query.RangeVar)
		args = append(args, *rangeVar.Relname)
	}

	switch t.Kind() {
	case reflect.Ptr:
		if v.Elem().IsValid() {
			args = append(args, examineTables(v.Elem().Interface(), depth+1)...)
		}
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		depth--
		if v.Len() > 0 {
			print("[")
			for i := 0; i < v.Len(); i++ {
				depth++
				print("[%d] Type {%s} {", i, v.Index(i).Type().String())
				args = append(args, examineTables(v.Index(i).Interface(), depth+1)...)
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
			args = append(args, examineTables(reflect.ValueOf(value).Field(i).Interface(), depth+1)...)
		}
	}
	return args
}

type DropInTableDefinition struct {
	Rows    [][]interface{}
	Columns []string
}

type DropInTable func() (DropInTableDefinition, error)

type DropInTableMap map[string]DropInTable

func ReplaceTables(stmt interface{}, dropIns DropInTableMap) (interface{}, error) {
	return replaceTables(stmt, 0, dropIns)
}

func replaceTables(value interface{}, depth int, dropIns DropInTableMap) (interface{}, error) {
	print := func(msg string, args ...interface{}) {
		fmt.Printf("%s%s\n", strings.Repeat("\t", depth), fmt.Sprintf(msg, args...))
	}

	if tableItem, ok := value.(pg_query.RangeVar); ok {
		if dropInTableFunc, ok := dropIns[*tableItem.Relname]; ok {
			print("[#] Dropping in table |%s|", *tableItem.Relname)
			_, err := dropInTableFunc()
			if err != nil {
				return nil, err
			}

		}
		return tableItem, nil
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
				result, err := replaceTables(item.Interface(), depth+1, dropIns)
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
			result, err := replaceTables(actual.Interface(), depth+1, dropIns)
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
