package dto

type (
	AccountState struct {
		Balance           int64
		Code              string
		Data              string
		FrozenHash        string
		LastTransactionId *TransactionID
		SyncUtime         int64
	}
)
