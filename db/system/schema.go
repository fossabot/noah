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
 * Copyright 2018  'pq' Contributors Portions Copyright (C) 2011 Blake Mizerany
 * License https://github.com/lib/pq/blob/master/LICENSE.md
 */

package system

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/golang/protobuf/proto"
	"github.com/kataras/go-errors"
	"strings"
)

var (
	ErrTableAlreadyExists = errors.New("table with name [%s] already exists")
)

type SSchema baseContext

func (ctx *SSchema) CreateTable(table NTable) error {
	if existingTable, err := ctx.GetTable(table.TableName); err != nil {
		return err
	} else if existingTable != nil {
		return ErrTableAlreadyExists.Format(table.TableName)
	}
	b, err := proto.Marshal(&table)
	if err != nil {
		return err
	}
	return ctx.db.Set([]byte(fmt.Sprintf("%s%s", tablesPath, strings.ToLower(table.TableName))), b)
}

func (ctx *SSchema) GetTable(tableName string) (*NTable, error) {
	bytes, err := ctx.db.Get([]byte(fmt.Sprintf("%s%s", tablesPath, strings.ToLower(tableName))))
	if err != nil {
		return nil, err
	}
	if len(bytes) == 0 {
		return nil, nil
	}
	table := NTable{}
	err = proto.Unmarshal(bytes, &table)
	if err != nil {
		return nil, err
	}
	return &table, nil
}

func (ctx *SSchema) GetTables() ([]NTable, error) {
	tableBytes, err := ctx.db.GetPrefix([]byte(tablesPath))
	if err != nil {
		return nil, err
	}
	tables := make([]NTable, len(tableBytes))
	for i, kv := range tableBytes {
		table := NTable{}
		err := proto.Unmarshal(kv.Value, &table)
		if err != nil {
			return nil, err
		}
		tables[i] = table
	}
	return tables, nil
}

func (ctx *SSchema) GetAccountsTable() (*NTable, error) {
	tables, err := ctx.GetTables()
	if err != nil {
		return nil, err
	}
	table := linq.From(tables).FirstWithT(func(t NTable) bool {
		return t.TableType == NTableType_ACCOUNT
	}).(NTable)
	return &table, nil
}

func (ctx *SSchema) DropTable(tableName string) (*NTable, error) {
	return nil, nil
}
