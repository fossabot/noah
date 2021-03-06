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

package coordinator

import (
	"github.com/readystock/golog"
	"github.com/readystock/noah/db/base"
	"github.com/readystock/noah/db/sql/pgwire"
	"github.com/readystock/noah/db/system"
	"github.com/readystock/noah/db/util"
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

func Start(sctx *system.SContext) (err error) {
	defer util.CatchPanic(&err)
	for {
		conn, err := sctx.PGListen.AcceptTCP()
		if err != nil {
			golog.Error(err.Error())
		}

		go func() {
			if err := handleConnection(sctx, conn); err != nil {
				golog.Errorf("failed handling connection %v", err)
			}
		}()
	}
}

func handleConnection(sctx *system.SContext, conn *net.TCPConn) error {
	golog.Infof("Handling connection from %s", conn.RemoteAddr().String())
	serv := pgwire.MakeServer(&base.Config{
		Insecure: true,
	})
	return serv.ServeConn(sctx, conn)
}
