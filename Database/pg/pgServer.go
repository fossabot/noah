package pg

import (
	"net"
	"fmt"
)

func StartIncomingConnection(in <-chan *net.TCPConn, out chan<- *net.TCPConn) {
	for conn := range in {
		handleConnection(conn)
		out <- conn
	}
}

func handleConnection(conn *net.TCPConn) {
	fmt.Println("Handling connection from ", conn.RemoteAddr().String())

	var buf ReadBuffer
	n, err := buf.ReadUntypedMsg(conn)
	if err != nil {
		return err
	}
	version, err := buf.GetUint32()
	if err != nil {
		return err
	}
	errSSLRequired := false
}