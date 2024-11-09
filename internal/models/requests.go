package models

type (
	GetTransactions struct {
		Addr string
		Hash string
		Lt   uint64
	}

	GetAccountState struct {
		Addr string
	}

	GetBetSeed struct {
		BetID int
	}

	GetActiveBets struct{}

	GetSeqno struct{}

	SendMessage struct {
		Body []byte
	}
)
