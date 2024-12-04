//go:generate mockgen -source=./interfaces.go -destination=./interfaces_mock.go -package=fetcher DiceContract,Repository

package fetcher

import (
	"context"

	"github.com/tongineers/dice-ton-api/internal/models"
)

type (
	DiceContract interface {
		SubscribeOnFinishedBets(ctx context.Context, fromLT uint64) (<-chan *models.Bet, error)
	}

	Repository interface {
		GetLastResolvedBetLT() (uint64, error)
		Update(bets ...*models.Bet) error
	}
)
