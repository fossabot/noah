package pg

import (
	"net"
	"fmt"
	"bytes"
	"github.com/Ready-Stock/Noah/Configuration"
)

func StartIncomingConnection(in <-chan *net.TCPConn, out chan<- *net.TCPConn) {
	for conn := range in {
		handleConnection(conn)
		out <- conn
	}
}

func handleConnection(conn *net.TCPConn) {
	fmt.Println("Handling connection from ", conn.RemoteAddr().String())
	buf := &bytes.Buffer{}
	for {
		data := make([]byte, Conf.Configuration.Database.ReadBuffer)
		n, err := conn.Read(data)
		if err != nil {
			fmt.Println(err.Error())
		}
		buf.Write(data[:n])
		if data[0] == 13 && data[1] == 10 {
			break
		}
	}
}