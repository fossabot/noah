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
		var res ResultBase
		switch cmd.(type) {
		case ExecStmt:

		case ExecPortal:
			// ExecPortal is handled like ExecStmt, except that the placeholder info
			// is taken from the portal.

		case PrepareStmt:
			res = ex.clientComm.CreatePrepareResult(pos)
		case DescribeStmt:

		case BindStmt:
			res = ex.clientComm.CreateBindResult(pos)
		case DeletePreparedStmt:
			res = ex.clientComm.CreateDeleteResult(pos)
		case SendError:
			res = ex.clientComm.CreateErrorResult(pos)
		case Sync:
			// Note that the Sync result will flush results to the network connection.
			res = ex.clientComm.CreateSyncResult(pos)
		case CopyIn:
			res = ex.clientComm.CreateCopyInResult(pos)
		case DrainRequest:
			// We received a drain request. We terminate immediately if we're not in a
			// transaction. If we are in a transaction, we'll finish as soon as a Sync
			// command (i.e. the end of a batch) is processed outside of a
			// transaction.

		case Flush:
			// Closing the res will flush the connection's buffer.
			res = ex.clientComm.CreateFlushResult(pos)
		default:
			panic(fmt.Sprintf("unsupported command type: %T", cmd))
		}

		ex.stmtBuf.advanceOne()

		res.Close()
	}
}