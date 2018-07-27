package sqlbase

import (
	"github.com/Ready-Stock/Noah/Database/sql/pgwire/pgerror"
)

const (
	txnAbortedMsg = "current transaction is aborted, commands ignored " +
		"until end of transaction block"
	txnCommittedMsg = "current transaction is committed, commands ignored " +
		"until end of transaction block"
	txnRetryMsgPrefix = "restart transaction"
)

func NewRetryError(cause error) error {
	return pgerror.NewErrorf(
		pgerror.CodeSerializationFailureError, "%s: %s", txnRetryMsgPrefix, cause)
}