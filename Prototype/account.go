package Prototype

var (
	Accounts = map[int64]Account{
		1: Account {
			AccountID: 1,
			AccountName: "Test Account 1",
		},
		2: Account {
			AccountID: 2,
			AccountName: "Test Account 2",
		},
	}
)

type Account struct {
	AccountID int64
	AccountName string
}
