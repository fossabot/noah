package driver

import (
	"github.com/kataras/go-errors"
)

type TransactionIsolationLevel string

// Transaction isolation levels
const (
	Serializable    = TransactionIsolationLevel("serializable")
	RepeatableRead  = TransactionIsolationLevel("repeatable read")
	ReadCommitted   = TransactionIsolationLevel("read committed")
	ReadUncommitted = TransactionIsolationLevel("read uncommitted")
)

type TransactionAccessMode string

// Transaction access modes
const (
	ReadWrite = TransactionAccessMode("read write")
	ReadOnly  = TransactionAccessMode("read only")
)

type TransactionDeferrableMode string

// Transaction deferrable modes
const (
	Deferrable    = TransactionDeferrableMode("deferrable")
	NotDeferrable = TransactionDeferrableMode("not deferrable")
)

const (
	TransactionStatusInProgress      = 0
	TransactionStatusCommitFailure   = -1
	TransactionStatusRollbackFailure = -2
	TransactionStatusInFailure       = -3
	TransactionStatusCommitSuccess   = 1
	TransactionStatusRollbackSuccess = 2
)

type TransactionOptions struct {
	IsoLevel       TransactionIsolationLevel
	AccessMode     TransactionAccessMode
	DeferrableMode TransactionDeferrableMode
}

var (
	ErrorTransactionClosed = errors.New("transaction is closed")
	ErrorTransactionFailure = errors.New("transaction has failed")

	// ErrorTransactionCommitRollback occurs when an error has occurred in a transaction and
	// Commit() is called. PostgreSQL accepts COMMIT on aborted transactions, but
	// it is treated as ROLLBACK.
	ErrorTransactionCommitRollback = errors.New("commit unexpectedly resulted in rollback")
)

