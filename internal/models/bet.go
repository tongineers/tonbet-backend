package models

import (
	"time"
)

type (
	Bet struct {
		ID            int       `json:"id"            gorm:"primarykey"`
		RollUnder     int       `json:"rollUnder"     gorm:"default:(-);"`
		RandomRoll    int       `json:"randomRoll"`
		Amount        uint64    `json:"amount"        gorm:"default:(-);"`
		Payout        uint64    `json:"payout"`
		PlayerAddress string    `json:"playerAddress" gorm:"default:(-);"`
		RefAddress    string    `json:"refAddress"    gorm:"default:(-);"`
		RefBonus      uint64    `json:"refBonus"`
		Seed          string    `json:"seed"          gorm:"default:(-);"`
		Status        BetStatus `json:"status"`
		LastLT        uint64    `json:"lastLT"`
		LastHash      string    `json:"lastHash"`
		CreatedAt     time.Time `json:"createdAt"     gorm:"default:(-);"`
		UpdatedAt     time.Time `json:"updatedAt"`
	}

	BetStatus int8
)

const (
	BetStatusNew BetStatus = iota
	BetStatusSent
	BetStatusResolved
)
