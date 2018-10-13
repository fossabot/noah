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
 */

package system

import (
	"flag"
	"github.com/Ready-Stock/Noah/db/util/snowflake"
	"github.com/Ready-Stock/badger"
	"github.com/Ready-Stock/raft-badger"
)

const (
	CoordinatorsPath           = "/coordinators/"
	NodesPath                  = "/nodes/"
	SettingsPath               = "/settings/"
	TablesPath                 = "/tables/"
	NodeIDSequencePath         = "/sequences/internal/nodes"
	AccountIDSequencePath      = "/sequences/internal/accounts"
	CoordinatorIDSequencePath  = "/sequences/internal/coordinators"
	PreloadPoolConnectionCount = 5
)

type BaseContext struct {
	db   *raft_badger.Store
}

type SContext struct {
	BaseContext
	NodeIDSequence *badger.Sequence
	Snowflake      *snowflake.Snowflake
	Flags          SFlags
	Pool           *NodePool
	Wal            NWal
	Schema         *NSchema
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
	sctx := SContext{
		Flags: SFlags{
			HTTPPort:      HttpPort,
			PostgresPort:  PostgresPort,
			DataDirectory: DataDirectory,
			WalDirectory:  WalDirectory,
			LogLevel:      LogLevel,
		},
		Snowflake: snowflake.NewSnowflake(1),
	}

	return &sctx, nil
}

func (ctx *SContext) Close() {
	ctx.db.Close()
}

func (ctx *BaseContext) GetSettings() (*map[string]string, error) {
	m := map[string]string{}
	if values, err := ctx.db.GetPrefix([]byte(SettingsPath)); err != nil {
		return nil, err
	} else {
		for _, value := range values {
			m[string(value.Key[len([]byte(SettingsPath))-1:])] = string(value.Value)
		}
	}
}
