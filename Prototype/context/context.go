package context

type TxnState string

const (
	StateNoTxn        TxnState = "StateNoTxn"
	StateInTxn        TxnState = "StateInTxn"
	StateCommittedTxn TxnState = "StateCommittedTxn"
)

type SessionContext struct {
	TransactionState TxnState
	Nodes            map[int]NodeContext
}

type NodeContext struct {
	NodeID           int
	TransactionState TxnState
}
