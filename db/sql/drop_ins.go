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

package sql

import (
	"github.com/readystock/noah/db/util/queryutil"
	"github.com/readystock/pg_query_go/nodes"
)

type dropInExecutor connExecutor

func (ex *connExecutor) GetDropInFunctions() queryutil.BuiltInFunctionMap {
	dropInEx := dropInExecutor(*ex)
	return queryutil.BuiltInFunctionMap{
		"pg_catalog.current_database": dropInEx.CurrentDatabase,
		"pg_catalog.current_schema":   dropInEx.CurrentSchema,
		"pg_catalog.current_schemas":  dropInEx.CurrentSchemas,
		"CURRENT_USER":                dropInEx.CurrentUser,
	}
}

func (ex *dropInExecutor) CurrentDatabase(args ...pg_query.Node) (interface{}, error) {
	return "noah", nil
}

func (ex *dropInExecutor) CurrentSchema(args ...pg_query.Node) (interface{}, error) {
	return "public", nil
}

func (ex *dropInExecutor) CurrentSchemas(args ...pg_query.Node) (interface{}, error) {
	return []string{"public"}, nil
}

func (ex *dropInExecutor) CurrentUser(args ...pg_query.Node) (interface{}, error) {
	return "noah", nil
}
