package sql

import (
	"net"
)

type SessionArgs struct {
	Database        string
	User            string
	ApplicationName string
	// RemoteAddr is the client's address. This is nil iff this is an internal
	// client.
	RemoteAddr net.Addr
}
