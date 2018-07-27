package Database

import (
	"net"
	"github.com/Ready-Stock/Noah/Configuration"
	"github.com/pkg/errors"
	"fmt"
	"github.com/Ready-Stock/Noah/Database/sql/pgwire"
	"github.com/Ready-Stock/Noah/Database/base"
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

func Start() error  {
	if addr, err := net.ResolveTCPAddr("tcp", Conf.Configuration.Database.AdvertiseAddress); err != nil {
		return errors.Errorf("unable to resolve RPC address %q: %v", Conf.Configuration.Database.AdvertiseAddress, err)
	} else {
		listener, err := net.ListenTCP("tcp", addr)
		if err != nil {
			return errors.Errorf("unable to listen on address %q: %v", Conf.Configuration.Database.AdvertiseAddress, err)
		}

		pending, complete := make(chan *net.TCPConn), make(chan *net.TCPConn)

		for i := 0; i < 5; i++ {
			go StartIncomingConnection(pending, complete)
		}

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				panic(err)
			}
			pending <- conn
		}


	}
}

func StartIncomingConnection(in <-chan *net.TCPConn, out chan<- *net.TCPConn) {
	for conn := range in {
		handleConnection(conn)
		out <- conn
	}
}

func handleConnection(conn *net.TCPConn) error {
	fmt.Println("Handling connection from ", conn.RemoteAddr().String())

	serv := pgwire.MakeServer(&base.Config{
		Insecure:true,
	})
	serv.ServeConn(conn)

	return nil
}


