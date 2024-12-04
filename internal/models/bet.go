package models

const (
	BetStatusNew BetStatus = iota
	BetStatusSent
	BetStatusResolved
)

type (
	Bet struct {
		ID            int       `json:"id"`
		RollUnder     int       `json:"rollUnder" gorm:"default:(-);"`
		RandomRoll    int       `json:"randomRoll"`
		Amount        uint64    `json:"amount" gorm:"default:(-);"`
		Payout        uint64    `json:"payout"`
		PlayerAddress string    `json:"playerAddress" gorm:"default:(-);"`
		RefAddress    string    `json:"refAddress" gorm:"default:(-);"`
		RefBonus      uint64    `json:"refBonus"`
		Seed          string    `json:"seed" gorm:"default:(-);"`
		Status        BetStatus `json:"status"`
		LastLT        uint64    `json:"lastLT"`
		LastHash      string    `json:"lastHash"`
	}

	BetStatus int8
)
