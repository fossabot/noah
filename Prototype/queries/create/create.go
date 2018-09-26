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
 *
 * This application uses Open Source components. You can find the
 * source code of their open source projects along with license
 * information below. We acknowledge and are grateful to these
 * developers for their contributions to open source.
 *
 * Project: CockroachDB https://github.com/cockroachdb/cockroach
 * Copyright 2018 The Cockroach Authors.
 * License (Apache License 2.0) https://github.com/cockroachdb/cockroach/blob/master/LICENSE
 *
 * Project: Vitess https://github.com/vitessio/vitess
 * Copyright 2018 Google Inc.
 * License (Apache License 2.0) https://github.com/vitessio/vitess/blob/master/LICENSE
 *
 * Project: Citus https://github.com/citusdata/citus
 * Copyright 2018 Citus Data, Inc.
 * License (GNU Affero General Public License v3.0) https://github.com/citusdata/citus/blob/master/LICENSE
 *
 * Project: pg_query_go https://github.com/lfittl/pg_query_go
 * Copyright 2018 Lukas Fittl
 * License (3-Clause BSD) https://github.com/lfittl/pg_query_go/blob/master/LICENSE
 *
 * Project: pgx https://github.com/jackc/pgx
 * Copyright 2018 Jack Christensen
 * License (MIT) https://github.com/jackc/pgx/blob/master/LICENSE
 *
 * Project: BadgerDB https://github.com/dgraph-io/badger
 * Copyright 2018 Dgraph Labs, Inc. and Contributors
 * License (MIT) https://github.com/dgraph-io/badger/blob/master/LICENSE
 */

package create

import (
	"github.com/Ready-Stock/pg_query_go/nodes"
	pgq "github.com/Ready-Stock/pg_query_go"
	"github.com/Ready-Stock/Noah/Prototype/context"
	"fmt"
	"github.com/Ready-Stock/Noah/Prototype/cluster"
	"github.com/Ready-Stock/Noah/Prototype/datums"
)

type CreateStatement struct {
	Statement pg_query.CreateStmt
	Query     string
}

func CreateCreateStatment(stmt pg_query.CreateStmt, tree pgq.ParsetreeList) CreateStatement {
	return CreateStatement{
		Statement: stmt,
		Query:     tree.Query,
	}
}

func (stmt CreateStatement) HandleCreate(ctx *context.SessionContext) error {
	fmt.Printf("Preparing Create Query\n")
	j, _ := stmt.Statement.MarshalJSON()
	fmt.Println(string(j))

	var has_id bool
	var id_index int
	var id_name string

	if has_id, id_index, id_name = stmt.hasIdentityColumn(); has_id {
		fmt.Printf("Found identity column at index: (%d) name: (%s)\n", id_index, id_name)
	}

	ids := ctx.GetAllNodes()

	response := ctx.DistributeQuery(stmt.Query, ids...)

	if response.Success {
		cluster.Tables[*stmt.Statement.Relation.Relname] = datums.Table{
			TableName:*stmt.Statement.Relation.Relname,
			IsGlobal:false,
			IsTenantTable:false,
			IdentityColumn:id_name,
		}
	}

	return ctx.HandleResponse(response)
}

func (stmt CreateStatement) hasIdentityColumn() (bool, int, string) {
	identity_index := -1
	for i, telt := range stmt.Statement.TableElts.Items {
		col := telt.(pg_query.ColumnDef)
		if col.Constraints.Items != nil && len(col.Constraints.Items) > 0 {
			for _, c := range col.Constraints.Items {
				constraint := c.(pg_query.Constraint)
				if constraint.Contype == pg_query.CONSTR_PRIMARY {
					identity_index = i
					return true, identity_index, *col.Colname
				}
			}
		}
	}
	return false, identity_index, ""
}
