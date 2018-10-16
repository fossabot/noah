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
 */

package system

import (
	"flag"
	"github.com/Ready-Stock/Noah/db/util/snowflake"
	"github.com/Ready-Stock/raft-badger"
)

const (
	CoordinatorsPath           = "/coordinators/"
	SettingsPath               = "/settings/"
	TablesPath                 = "/tables/"
	NodeIDSequencePath         = "/sequences/internal/nodes"
	AccountIDSequencePath      = "/sequences/internal/accounts"
	CoordinatorIDSequencePath  = "/sequences/internal/coordinators"
	PreloadPoolConnectionCount = 5
)

type baseContext struct {
	db        *raft_badger.Store
	snowflake *snowflake.Snowflake
}

type SContext struct {
	baseContext
	Settings  *SSettings
	Accounts  *SAccounts
	Schema    *SSchema
	Pool      *SPool
	Nodes     *SNode
	Sequences *SSequence
	Flags     SFlags
}

type SFlags struct {
	HTTPPort      int
	PostgresPort  int
	DataDirectory string
	WalDirectory  string
	LogLevel      string
}

func NewSystemContext() (*SContext, error) {
	flag.Parse()
	db, err := raft_badger.CreateStore("data", ":5431", ":5430", "")
	if err != nil {
		return nil, err
	}
	base := baseContext{
		snowflake: snowflake.NewSnowflake(1),
		db:        db,
	}
	sctx := SContext{
		Flags: SFlags{
			HTTPPort:      HttpPort,
			PostgresPort:  PostgresPort,
			DataDirectory: DataDirectory,
			WalDirectory:  WalDirectory,
			LogLevel:      LogLevel,
		},
	}

	settings := SSettings(base)
	accounts := SAccounts(base)
	schema := SSchema(base)
	sequences := SSequence(base)
	nodes := SNode(base)
	pool := SPool{baseContext: &base}
	sctx.Settings = &settings
	sctx.Accounts = &accounts
	sctx.Schema = &schema
	sctx.Pool = &pool
	sctx.Nodes = &nodes
	sctx.Sequences = &sequences
	return &sctx, nil
}

func (ctx *SContext) Close() {
	ctx.db.Close()
}
