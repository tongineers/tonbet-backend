package fetcher

import (
	"context"

	"go.uber.org/zap"

	"github.com/tongineers/tonbet-backend/config"
	"github.com/tongineers/tonbet-backend/internal/models"
	"github.com/tongineers/tonbet-backend/pkg/workerpool"
)

var (
	_ workerpool.Task = (*Service)(nil)
)

type (
	// Service ...
	Service struct {
		dice   DiceContract
		repo   Repository
		conf   *config.Config
		logger *zap.Logger
	}
)

// New ...
func New(
	dice DiceContract,
	repo Repository,
	conf *config.Config,
	logger *zap.Logger,
) *Service {
	return &Service{
		dice:   dice,
		repo:   repo,
		conf:   conf,
		logger: logger,
	}
}

// Do ...
func (s *Service) Do() error {
	lastLT, err := s.repo.GetLastResolvedBetLT()
	if err != nil {
		return err
	}

	if lastLT == 0 {
		lastLT = s.conf.TONLastTransactionLT
	}

	bets, err := s.dice.SubscribeOnFinishedBets(context.Background(), lastLT)
	if err != nil {
		return err
	}

	for bet := range bets {
		s.logger.Info("new game result detected",
			zap.Int("betID", bet.ID),
			zap.Int("rollUnder", bet.RollUnder),
			zap.Int("randomRoll", bet.RandomRoll),
			zap.Uint64("payout", bet.Payout),
			zap.Uint64("lastLT", bet.LastLT),
			zap.String("lastHash", bet.LastHash),
		)

		copy := *bet
		copy.Status = models.BetStatusResolved

		err = s.repo.Update(&copy)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Run() {
	for {
		err := s.Do()
		if err != nil {
			continue
		}
	}
}
