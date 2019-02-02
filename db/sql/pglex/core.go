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

package pglex

import (
	"fmt"
	node "github.com/readystock/pg_query_go/nodes"
)

func HandleRawStmt(stmt node.RawStmt) error {
	switch tstmt := stmt.Stmt.(type) {
	case node.VariableSetStmt:
		fmt.Printf("Changing Setting (%s)\n", *tstmt.Name)
	}
	return nil
}
