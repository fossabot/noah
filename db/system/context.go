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

package system

import (
	"github.com/pkg/errors"
	"github.com/readystock/arctonyx"
	"github.com/readystock/noah/db/sql/driver/npgx"
	"github.com/readystock/noah/db/util/snowflake"
	"net"
	"time"
)

type baseContext struct {
	db        *arctonyx.Store
	snowflake *snowflake.Snowflake
}

type SContext struct {
	baseContext

	PGListen  *net.TCPListener
	isRunning bool

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

	addr, err := net.ResolveTCPAddr("tcp", pgWireAddr)
	if err != nil {
		panic(errors.Errorf("unable to resolve address %q: %v", pgWireAddr, err))
	}

	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(errors.Errorf("unable to listen on address %q: %v", addr, err))
	}

	sctx := SContext{
		baseContext: base,
		PGListen:    listen,
	}

	settings := SSettings(base)
	accounts := SAccounts(base)
	schema := SSchema(base)
	sequences := SSequence(base)
	nodes := SNode(base)
	query := SQuery(base)
	pool := SPool{baseContext: &base, nodePools: map[uint64]*npgx.ConnPool{}}
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

func (ctx *SContext) IsLeader() bool {
	return ctx.db.IsLeader()
}

func (ctx *SContext) GrpcListenAddr() string {
	return ctx.db.ListenAddr()
}

func (ctx *SContext) PgListenAddr() string {
	return ctx.PGListen.Addr().String()
}

func (ctx *SContext) FinishedStartup() {
	ctx.isRunning = true
}

func (ctx *SContext) IsRunning() bool {
	return ctx.isRunning
}

func (ctx *SContext) Close() {
	ctx.db.Close()
}
