package sql

import (
	"io"
	"fmt"
)

type Server struct {

}

type connExecutor struct {
	server *Server
	stmtBuf *StmtBuf
	clientComm ClientComm
}

func NewServer() *Server {
	return &Server{

	}
}

func (s *Server) ServeConn(stmtBuf *StmtBuf, clientComm ClientComm) error {
	ex := s.newConnExecutor(stmtBuf, clientComm)
	return ex.run()
}

func (s *Server) newConnExecutor(stmtBuf *StmtBuf, clientComm ClientComm) *connExecutor {
	ex := &connExecutor{
		server:      s,
		stmtBuf:     stmtBuf,
		clientComm:  clientComm,
	}
	return ex
}

func (ex *connExecutor) run() error {
	for {
		cmd, pos, err := ex.stmtBuf.curCmd()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		fmt.Printf("[pos:%d] executing %s\n", pos, cmd)
		switch cmd.(type) {
		case ExecStmt:
			ex.clientComm.CreateEmptyQueryResult(pos)

		case ExecPortal:
			ex.clientComm.CreateEmptyQueryResult(pos)

		}
		ex.stmtBuf.advanceOne()
	}
}