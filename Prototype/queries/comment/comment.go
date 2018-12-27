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
 *
 * Project: Sonyflake https://github.com/sony/sonyflake
 * Copyright 2018 Sony Corporation
 * License (MIT) https://github.com/sony/sonyflake/blob/master/LICENSE
 *
 * Project: Raft https://github.com/hashicorp/raft
 * Copyright 2018 HashiCorp
 * License (MPL-2.0) https://github.com/hashicorp/raft/blob/master/LICENSE
 *
 * Project: pq github.com/lib/pq
 * Copyright 2018  'pq' Contributors Portions Copyright (C) 2018 Blake Mizerany
 * License https://github.com/lib/pq/blob/master/LICENSE.md
 *
 * Project: go-linq github.com/ahmetb/go-linq
 * Copyright 2018 Ahmet Alp Balkan
 * License https://github.com/ahmetb/go-linq/blob/master/LICENSE
 */

package comment

import (
	"fmt"
	"github.com/kataras/go-errors"
	"github.com/readystock/noah/Prototype/cluster"
	"github.com/readystock/noah/Prototype/context"
	pgq "github.com/readystock/pg_query_go"
	"github.com/readystock/pg_query_go/nodes"
	"strings"
)

var (
	errorTableNotFound     = errors.New("table (%s) cannot be found in metadata")
	errorTableTypeNotValid = errors.New("table type (%s) is not valid")
)

type CommentStatement struct {
	Statement pg_query.CommentStmt
	Query     string
}

func CreateCommentStatment(stmt pg_query.CommentStmt, tree pgq.ParsetreeList) CommentStatement {
	return CommentStatement{
		Statement: stmt,
		Query:     tree.Query,
	}
}

func (stmt CommentStatement) HandleComment(ctx *context.SessionContext) error {
	fmt.Printf("Preparing Comment Query\n")
	j, _ := stmt.Statement.MarshalJSON()
	fmt.Println(string(j))

	if stmt.Statement.Object != nil && len(stmt.Statement.Object.(pg_query.List).Items) > 0 {
		lst := stmt.Statement.Object.(pg_query.List)
		table := lst.Items[len(lst.Items)-1].(pg_query.String).Str
		if tblmeta, ok := cluster.Tables[table]; !ok {
			return errorTableNotFound.Format(table)
		} else {
			switch strings.ToUpper(*stmt.Statement.Comment) {
			case "GLOBAL":
				tblmeta.IsGlobal = true
				tblmeta.IsTenantTable = false
			case "TENANTS":
				tblmeta.IsGlobal = true
				tblmeta.IsTenantTable = true
			case "LOCAL":
				tblmeta.IsGlobal = false
				tblmeta.IsTenantTable = false
			default:
				return errorTableTypeNotValid.Format(*stmt.Statement.Comment)
			}
			cluster.Tables[table] = tblmeta
		}
	}
	ids := ctx.GetAllNodes()
	response := ctx.DistributeQuery(stmt.Query, ids...)
	return ctx.HandleResponse(response)
}
