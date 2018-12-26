package sql

import (
	"github.com/readystock/noah/db/system"
)

func CreateConnExecutor(sctx *system.SContext) *connExecutor {
	ex := &connExecutor{
		prepStmtsNamespace: prepStmtNamespace{
			prepStmts: make(map[string]prepStmtEntry),
			portals:   make(map[string]portalEntry),
		},
		SystemContext:    sctx,
		TransactionState: TransactionState_None,
		TransactionMode:  TransactionMode_AutoCommit,
	}
	return ex
}
