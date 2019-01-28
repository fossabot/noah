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

func FindAccountIds(stmt interface{}, shardColumnName string) []uint64 {
	return examineWhereClause(stmt, 0, shardColumnName)
}

func examineWhereClause(value interface{}, depth int, shardColumnName string) []uint64 {
	ids := make([]uint64, 0)
	print := func(msg string, args ...interface{}) {
		// fmt.Printf("%s%s\n", strings.Repeat("\t", depth), fmt.Sprintf(msg, args...))
	}

	if value == nil {
		return ids
	}

	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)

	if v.Type() == reflect.TypeOf(pg_query.A_Expr{}) {
		expr := value.(pg_query.A_Expr)
		colRef, ok := expr.Lexpr.(pg_query.ColumnRef)
		if !ok {
			goto Continue
		}

		if colRef.Fields.Items[len(colRef.Fields.Items)-1].(pg_query.String).Str != shardColumnName {
			goto Continue
		}

		valueConst, ok := expr.Rexpr.(pg_query.A_Const)
		if !ok {
			goto Continue
		}

		numericValue, ok := valueConst.Val.(pg_query.Integer)
		if !ok {
			goto Continue
		}
		ids = append(ids, uint64(numericValue.Ival))
		return ids
	}

Continue:
	switch t.Kind() {
	case reflect.Ptr:
		if v.Elem().IsValid() {
			ids = append(ids, examineWhereClause(v.Elem().Interface(), depth+1, shardColumnName)...)
		}
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		depth--
		if v.Len() > 0 {
			print("[")
			for i := 0; i < v.Len(); i++ {
				depth++
				print("[%d] Type {%s} {", i, v.Index(i).Type().String())
				ids = append(ids, examineWhereClause(v.Index(i).Interface(), depth+1, shardColumnName)...)
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
			ids = append(ids, examineWhereClause(reflect.ValueOf(value).Field(i).Interface(), depth+1, shardColumnName)...)
		}
	}
	return ids
}
