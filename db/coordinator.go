package db

import (
	"fmt"
	"github.com/Ready-Stock/Noah/db/base"
	"github.com/Ready-Stock/Noah/db/sql/pgwire"
	"github.com/Ready-Stock/Noah/db/system"
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

func Start(sctx *system.SContext) error  {
	addvertise_addr := fmt.Sprintf("127.0.0.1:%d", sctx.Flags.PostgresPort)
	if addr, err := net.ResolveTCPAddr("tcp", addvertise_addr); err != nil {
		return errors.Errorf("unable to resolve RPC address %q: %v", addvertise_addr, err)
	} else {
		listener, err := net.ListenTCP("tcp", addr)
		if err != nil {
			return errors.Errorf("unable to listen on address %q: %v", addvertise_addr, err)
		}

		pending, complete := make(chan *net.TCPConn), make(chan *net.TCPConn)

		for i := 0; i < 5; i++ {
			go StartIncomingConnection(sctx, pending, complete)
		}

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				golog.Error(err.Error())
			}
			pending <- conn
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


