package listener

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
	DoDelay = 1 * time.Second
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
	bets, err := s.dice.GetActiveBets()
	if err != nil {
		return err
	}

	if len(bets) == 0 {
		return nil
	}

	ids := make([]int, 0)
	for _, bet := range bets {
		ids = append(ids, bet.ID)
	}

	existedBets, err := s.repo.ReadByIDs(ids...)
	if err != nil {
		return err
	}

	toUpdate := make([]*models.Bet, 0)
	for _, bet := range bets {
		alreadyExist := false
		for _, existed := range existedBets {
			if bet.ID == existed.ID {
				alreadyExist = true
				break
			}
		}

		if alreadyExist {
			continue
		}

		copy := *bet
		copy.Status = models.BetStatusNew
		copy.CreatedAt = time.Now()

		s.logger.Info("new active bet found",
			zap.Int("id", bet.ID),
			zap.Int("rollUnder", bet.RollUnder),
			zap.Uint64("amount", bet.Amount),
			zap.String("playerAddress", bet.PlayerAddress),
			zap.String("refAddress", bet.RefAddress),
			zap.String("seed", bet.Seed),
		)

		toUpdate = append(toUpdate, &copy)
	}

	if len(toUpdate) == 0 {
		return nil
	}

	return s.repo.Update(toUpdate...)
}

func (s *Service) Run() {
	for {
		select {
		case <-time.After(DoDelay):
			s.Do()
		}
	}
}
