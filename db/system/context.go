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
    "flag"
    "github.com/readystock/arctonyx"
    "github.com/readystock/noah/db/util/snowflake"
    "time"
)

type baseContext struct {
    db        *arctonyx.Store
    snowflake *snowflake.Snowflake
}

type SContext struct {
    baseContext

    PGWireAddress string

    Settings  *SSettings
    Accounts  *SAccounts
    Schema    *SSchema
    Pool      *SPool
    Nodes     *SNode
    Sequences *SSequence
    Setup     *SSetup
    Query     *SQuery
}

func NewSystemContext(dataDirectory, listenAddr, joinAddr, pgWireAddr string) (*SContext, error) {
    flag.Parse()
    db, err := arctonyx.CreateStore(dataDirectory, listenAddr, joinAddr)
    if err != nil {
        return nil, err
    }
    // Wait 5 seconds, this should be enough for the store to elect itself as leader if needed.
    time.Sleep(5 * time.Second)

    base := baseContext{
        snowflake: snowflake.NewSnowflake(uint16(db.NodeID())),
        db:        db,
    }

    sctx := SContext{
        baseContext:   base,
        PGWireAddress: pgWireAddr,
    }

    settings := SSettings(base)
    accounts := SAccounts(base)
    schema := SSchema(base)
    sequences := SSequence(base)
    nodes := SNode(base)
    query := SQuery(base)
    pool := SPool{baseContext: &base}
    setup := SSetup(base)

    sctx.Settings = &settings
    sctx.Accounts = &accounts
    sctx.Schema = &schema
    sctx.Pool = &pool
    sctx.Nodes = &nodes
    sctx.Sequences = &sequences
    sctx.Query = &query
    sctx.Setup = &setup
    return &sctx, nil
}

// func (ctx *SContext) InitStore() error {
//      val, err := ctx.Settings.GetSettingUint64(InitialSetupTimestamp)
//      if err != nil {
//          return err
//      }
//
//
// }

func (ctx *SContext) NewSnowflake() (uint64, error) {
    return ctx.snowflake.NextID()
}

func (ctx *SContext) CoordinatorID() uint64 {
    return ctx.db.NodeID()
}

func (ctx *SContext) Close() {
    ctx.db.Close()
}
