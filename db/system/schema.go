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

package system

import (
    "github.com/golang/protobuf/proto"
    "github.com/kataras/go-errors"
    "github.com/readystock/arctonyx"
    "github.com/readystock/golinq"
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
    return ctx.db.Set(getTablePath(table.TableName), b)
}

func (ctx *SSchema) GetTable(tableName string) (*NTable, error) {
    bytes, err := ctx.db.Get(getTablePath(tableName))
    if err != nil {
        return nil, err
    }
    if len(bytes) == 0 {
        return nil, nil
    }
    return ctx.convertBytesToTable(bytes)
}

func (ctx *SSchema) GetTables() ([]NTable, error) {
    tableBytes, err := ctx.db.GetPrefix(getTablesPath())
    if err != nil {
        return nil, err
    }
    tables, err := ctx.convertBytesToTables(tableBytes...)
    if err != nil {
        return nil, err
    }
    return tables, nil
}

func (ctx *SSchema) GetAccountsTable() (*NTable, error) {
    tables, err := ctx.GetTables()
    if err != nil {
        return nil, err
    }
    t := linq.From(tables).FirstWithT(func(t NTable) bool {
        return t.TableType == NTableType_ACCOUNT
    })
    if t != nil {
        table := t.(NTable)
        return &table, nil
    }
    return nil, nil
}

func (ctx *SSchema) DropTable(tableName string) (*NTable, error) {
    return nil, nil
}

func (ctx *SSchema) convertBytesToTable(bytes []byte) (*NTable, error) {
    table := NTable{}
    err := proto.Unmarshal(bytes, &table)
    if err != nil {
        return nil, err
    }
    return &table, err
}

func (ctx *SSchema) convertBytesToTables(tableBytes ...arctonyx.KeyValue) ([]NTable, error) {
    tables := make([]NTable, len(tableBytes))
    for i, kv := range tableBytes {
        table, err := ctx.convertBytesToTable(kv.Value)
        if err != nil {
            return nil, err
        }
        tables[i] = *table
    }
    return tables, nil
}