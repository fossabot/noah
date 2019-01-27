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
	"encoding/json"
	"github.com/readystock/golog"
	"github.com/readystock/pg_query_go"
	pg_query2 "github.com/readystock/pg_query_go/nodes"
	"testing"
)

var builtInFunctions = BuiltInFunctionMap{
	"pg_catalog.current_database": func(args ...pg_query2.Node) (i interface{}, e error) {
		return "noah", nil
	},
	"pg_catalog.current_schemas": func(args ...pg_query2.Node) (i interface{}, e error) {
		return []string{"noah", "test"}, nil
	},
	"pg_catalog.current_user": func(args ...pg_query2.Node) (i interface{}, e error) {
		return "ecourant", nil
	},
}

func Test_ReplaceFunctions(t *testing.T) {
	input := `select current_database(), current_schemas(), CURRENT_USER;`
	parsed, err := pg_query.Parse(input)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	stmt := parsed.Statements[0].(pg_query2.RawStmt).Stmt
	stmtJson, _ := json.Marshal(stmt)

	golog.Debugf("INPUT: | %s", string(stmtJson))
	newStmt, err := ReplaceFunctionCalls(stmt, builtInFunctions)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	newStmtJson, _ := json.Marshal(newStmt)
	golog.Debugf("OUTPUT: | %s", string(newStmtJson))
	compiled, err := newStmt.(pg_query2.Node).Deparse(pg_query2.Context_None)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	golog.Infof("IN  | %s", input)
	golog.Infof("OUT | %s", *compiled)
}
