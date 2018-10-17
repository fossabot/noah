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

package coordinator

import (
	"fmt"
	"github.com/Ready-Stock/Noah/db/base"
	"github.com/Ready-Stock/Noah/db/sql/pgwire"
	"github.com/Ready-Stock/Noah/db/system"
	"github.com/Ready-Stock/Noah/db/util"
	"github.com/kataras/golog"
	"github.com/pkg/errors"
	"net"
	"time"
)

const (
	// ErrSSLRequired is returned when a client attempts to connect to a
	// secure server in cleartext.
	ErrSSLRequired = "node is running secure mode, SSL connection required"

	// ErrDraining is returned when a client attempts to connect to a server
	// which is not accepting client connections.
	ErrDraining = "server is not accepting clients"
)

// Fully-qualified names for metrics.

const (
	version30  = 196608
	versionSSL = 80877103
)

// cancelMaxWait is the amount of time a draining server gives to sessions to
// react to cancellation and return before a forceful shutdown.
const cancelMaxWait = 1 * time.Second

// connReservationBatchSize determines for how many connections memory
// is pre-reserved at once.
var connReservationBatchSize = 5

var (
	sslSupported   = []byte{'S'}
	sslUnsupported = []byte{'N'}
)

type Server struct {

}

func Start(sctx *system.SContext) (err error)  {
	defer util.CatchPanic(&err)
	addvertise_addr := fmt.Sprintf("127.0.0.1:%d", sctx.Flags.PostgresPort)
	if addr, err := net.ResolveTCPAddr("tcp", addvertise_addr); err != nil {
		return errors.Errorf("unable to resolve RPC address %q: %v", addvertise_addr, err)
	} else {
		listener, err := net.ListenTCP("tcp", addr)
		if err != nil {
			return errors.Errorf("unable to listen on address %q: %v", addvertise_addr, err)
		}

		// pending, complete := make(chan *net.TCPConn), make(chan *net.TCPConn)
		//
		// for i := 0; i < 5; i++ {
		// 	go StartIncomingConnection(sctx, pending, complete)
		// }

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				golog.Error(err.Error())
			}
			go handleConnection(sctx, conn)
		}
	}
}

func StartIncomingConnection(sctx *system.SContext, in <-chan *net.TCPConn, out chan<- *net.TCPConn) {
	for conn := range in {
		handleConnection(sctx, conn)
		out <- conn
	}
}

func handleConnection(sctx *system.SContext, conn *net.TCPConn) error {
	golog.Infof("Handling connection from %s", conn.RemoteAddr().String())
	serv := pgwire.MakeServer(&base.Config{
		Insecure:true,
	})
	return serv.ServeConn(sctx, conn)
}


