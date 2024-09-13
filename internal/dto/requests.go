package dto

type (
	GetTransactions struct {
		Addr string
		Hash string
		Lt   int
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
