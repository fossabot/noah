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
    "github.com/readystock/noah/db/sql/pgwire/pgproto"
    "github.com/readystock/pg_query_go/nodes"
)

// getPlanColumns implements the logic for the
// planColumns/planMutableColumns functions. The mut argument
// indicates whether the slice should be mutable (mut=true) or not.
func getPlanColumns(stmt pg_query.Stmt) []pgproto.FieldDescription {
    return nil
}