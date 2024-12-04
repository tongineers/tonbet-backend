//go:generate mockgen -source=./interfaces.go -destination=./interfaces_mock.go -package=listener DiceContract,Repository

package listener

import (
	"github.com/tongineers/dice-ton-api/internal/models"
)

type (
	DiceContract interface {
		GetActiveBets() ([]*models.Bet, error)
	}

	Repository interface {
		ReadByIDs(ids ...int) ([]*models.Bet, error)
		Update(bets ...*models.Bet) error
	}
)
