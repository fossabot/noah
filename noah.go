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

package main

import (
	"github.com/kataras/golog"
	"github.com/readystock/noah/api"
	"github.com/readystock/noah/db/coordinator"
	"github.com/readystock/noah/db/system"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

type ServiceType string

const (
	Coordinator ServiceType = "coordinator"
	Tablet      ServiceType = "tablet"
)

var (
	systemContext *system.SContext
	ch            chan os.Signal
)

func main() {
	StartCoordinator()
}

func StartCoordinator() {
	sctx, err := system.NewSystemContext()
	if err != nil {
		panic(err)
	}
	golog.SetLevel(sctx.Flags.LogLevel)
	golog.Infof("Coordinator ID [%d] starting...", sctx.CoordinatorID())
	golog.Infof("Starting admin application with port [%d]", sctx.Flags.HTTPPort)
	golog.Infof("Listening for PostgreSQL connection on port [%d]", sctx.Flags.PostgresPort)
	systemContext = sctx
	ch = make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	signal.Notify(ch, os.Interrupt, syscall.SIGSEGV)
	go func() {
		<-ch
		StopCoordinator()
	}()

	logrus.SetLevel(logrus.DebugLevel)
	go api.StartApp(sctx)
	coordinator.Start(sctx)
}

func StartCoordinatorExt() {

}

func StopCoordinator() {
	systemContext.Close()
	os.Exit(0)
}
