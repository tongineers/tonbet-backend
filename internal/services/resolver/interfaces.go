//go:generate mockgen -source=./interfaces.go -destination=./interfaces_mock.go -package=resolver DiceContract,Repository

package resolver

import (
	"github.com/tongineers/dice-ton-api/internal/models"
)

type (
	DiceContract interface {
		ResolveBet(betID int, seed string) error
	}

	Repository interface {
		ReadByStatus(status models.BetStatus) ([]*models.Bet, error)
		Update(bets ...*models.Bet) error
	}
)
