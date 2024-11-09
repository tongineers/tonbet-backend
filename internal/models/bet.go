package models

const (
	BetStatusNew BetStatus = iota
	BetStatusResolved
	BetStatusFinished
)

type (
	Bet struct {
		ID            int32     `json:"id"`
		RollUnder     int32     `json:"rollUnder"`
		Amount        int64     `json:"amount"`
		PlayerAddress string    `json:"playerAddress"`
		RefAddress    string    `json:"refAddress"`
		RefBonus      int64     `json:"refBonus"`
		Seed          string    `json:"seed"`
		Status        BetStatus `json:"status"`
	}

	BetStatus int8
)
