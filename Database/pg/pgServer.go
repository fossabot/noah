package pg

import (
	"net"
	"fmt"
	"time"
	"github.com/pkg/errors"
	"github.com/Ready-Stock/Noah/Database/pg/pgwire/pgerror"
	"github.com/Ready-Stock/Noah/Database/pg/sql"
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



func StartIncomingConnection(in <-chan *net.TCPConn, out chan<- *net.TCPConn) {
	for conn := range in {
		handleConnection(conn)
		out <- conn
	}
}

func handleConnection(conn *net.TCPConn) error {
	fmt.Println("Handling connection from ", conn.RemoteAddr().String())

	var buf ReadBuffer
	_, err := buf.ReadUntypedMsg(conn)
	if err != nil {
		return err
	}
	version, err := buf.GetUint32()
	if err != nil {
		return err
	}

	// sendErr := func(err error) error {
	// 	msgBuilder := newWriteBuffer(s.metrics.BytesOutCount)
	// 	_ /* err */ = writeErr(err, msgBuilder, conn)
	// 	_ = conn.Close()
	// 	return err
	// }

	if version != version30 {
		return fmt.Errorf("unknown protocol version %d", version)
	}
	var sArgs sql.SessionArgs
	if sArgs, err = parseOptions(buf.Msg); err != nil {
		return pgerror.NewError(pgerror.CodeProtocolViolationError, err.Error())
	}
	fmt.Println("Connection for user ", sArgs.User)
	return nil
}

func parseOptions(data []byte) (sql.SessionArgs, error) {
	args := sql.SessionArgs{}
	buf := ReadBuffer{Msg: data}
	for {
		key, err := buf.GetString()
		if err != nil {
			return sql.SessionArgs{}, errors.Errorf("error reading option key: %s", err)
		}
		if len(key) == 0 {
			break
		}
		value, err := buf.GetString()
		if err != nil {
			return sql.SessionArgs{}, errors.Errorf("error reading option value: %s", err)
		}
		switch key {
		case "database":
			args.Database = value
		case "user":
			args.User = value
		case "application_name":
			args.ApplicationName = value
		default:

		}
	}
	return args, nil
}


