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

package main

import (
	"flag"
	"fmt"
	"github.com/Ready-Stock/Noah/db/coordinator"
	"github.com/Ready-Stock/Noah/db/system"
	"github.com/Ready-Stock/Noah/db/util/snowflake"
	"github.com/Ready-Stock/badger"
	"github.com/kataras/golog"
)

type ServiceType string

const (
	Coordinator ServiceType = "coordinator"
	Tablet      ServiceType = "tablet"
)

var (
	NodeIDSequencePath = []byte("/sequences/node_id_sequence")
)

var (
	RunType       = *flag.String("type", "coordinator", "Type of handler to run, defaults to `coordinator`. Valid values are: `tablet` and `coordinator`.")
	Role          = *flag.String("role", "master", "Type of role for the coordinator, not valid for tablets. Value values are `master` and `follower`. Defaults to `master`")
	Join          = *flag.String("join", "", "The IP of another coordinator in the cluster.")
	HttpPort      = *flag.Int("http-port", 8080, "Listen port for Noah's HTTP REST interface.")
	PostgresPort  = *flag.Int("psql-port", 5433, "Listen port for Noah's PostgreSQL client connectivity.")
	DataDirectory = *flag.String("data-dir", "data", "Directory for Noah's embedded database.")
	WalDirectory  = *flag.String("wal-dir", "wal", "Directory for Noah's write ahead log.")
	LogLevel      = *flag.String("log-level", "debug", "Log verbosity for message written to the console window.")
)

func main() {
	flag.Parse()
	golog.SetLevel(LogLevel)

	switch ServiceType(RunType) {
	case Coordinator:
		SystemContext := system.SContext{
			Flags: system.SFlags{
				HTTPPort:      HttpPort,
				PostgresPort:  PostgresPort,
				DataDirectory: DataDirectory,
				WalDirectory:  WalDirectory,
			},
			Snowflake: snowflake.NewSnowflake(1),
		}

		opts := badger.DefaultOptions

		opts.Dir = SystemContext.Flags.DataDirectory
		opts.ValueDir = SystemContext.Flags.DataDirectory
		badger_data, err := badger.Open(opts)
		if err != nil {
			panic(err)
		}
		node_seq, err := badger_data.GetSequence(NodeIDSequencePath, 10)
		if err != nil {
			panic(err)
		}

		SystemContext.Badger = badger_data

		defer badger_data.Close()
		defer badger_wal.Close()

		SystemContext.NodeIDs = node_seq
		fmt.Println("Starting admin application with port:", SystemContext.Flags.HTTPPort)
		fmt.Println("Listening for connections on:", SystemContext.Flags.PostgresPort)

		// go api.StartApp(&SystemContext)
		coordinator.Start(&SystemContext)

	case Tablet:
		golog.Info("Starting tablet...")
	}
}
