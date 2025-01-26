package resolver

import (
	"time"

	"go.uber.org/zap"

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
		logger *zap.Logger
	}
)

const (
	DoDelay          = 1 * time.Second
	DBBackOffTimeout = 3 * time.Second
)

// New ...
func New(
	dice DiceContract,
	repo Repository,
	logger *zap.Logger,
) *Service {
	return &Service{
		dice:   dice,
		repo:   repo,
		logger: logger,
	}
}

// Do ...
func (s *Service) Do() error {
	bets, err := s.repo.ReadByStatus(models.BetStatusNew)
	if err != nil {
		return err
	}

	if len(bets) == 0 {
		return nil
	}

	toUpdate := make([]*models.Bet, 0)
	for _, bet := range bets {
		s.logger.Info("start resolving bet with ID",
			zap.Int("id", bet.ID),
		)

		err := s.dice.ResolveBet(bet.ID, bet.Seed)
		if err != nil {
			s.logger.Error("error resolving bet with ID",
				zap.Int("id", bet.ID),
				zap.Error(err),
			)
			continue
		}

		copy := *bet
		copy.Status = models.BetStatusSent
		toUpdate = append(toUpdate, &copy)
	}

	if len(toUpdate) == 0 {
		return nil
	}

	for {
		err = s.repo.Update(toUpdate...)
		if err != nil {
			s.logger.Error("failed to save resolved bets to DB",
				zap.Error(err),
			)
			s.logger.Info("trying to save bets again after timeout...",
				zap.Duration("timeout", DBBackOffTimeout),
			)

			time.Sleep(DBBackOffTimeout)
			continue
		}

		break
	}

	return nil
}

func (s *Service) Run() {
	for {
		select {
		case <-time.After(DoDelay):
			s.Do()
		}
	}
}
