package Database

import (
	"net"
	"github.com/Ready-Stock/Noah/Configuration"
	"github.com/pkg/errors"
	"github.com/Ready-Stock/Noah/Database/pg"
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
			go pg.StartIncomingConnection(pending, complete)
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
