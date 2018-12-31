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

func (ctx *SSchema) UpdateTable(table NTable) error {
    if bytes, err := ctx.convertTableToBytes(table); err != nil {
        return err
    } else {
        return ctx.db.Set(getTablePath(table.TableName), bytes)
    }
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

func (ctx *SSchema) DropTable(tableName string) (error) {
    return ctx.db.Delete(getTablePath(tableName))
}

func (ctx *SSchema) convertTableToBytes(table NTable) ([]byte, error) {
    return proto.Marshal(&table)
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