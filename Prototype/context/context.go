package context

type TxnState string
const (
	StateNoTxn TxnState = "StateNoTxn"
	StateInTxn TxnState = "StateInTxn"
)

type SessionContext struct {
	TransactionState TxnState
}


