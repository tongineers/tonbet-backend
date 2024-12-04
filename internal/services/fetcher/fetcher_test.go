package fetcher_test

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"github.com/tongineers/dice-ton-api/config"
	"github.com/tongineers/dice-ton-api/internal/app/dependencies"
	"github.com/tongineers/dice-ton-api/internal/interfaces"
	"github.com/tongineers/dice-ton-api/internal/models"
	"github.com/tongineers/dice-ton-api/internal/services/fetcher"
	"github.com/tongineers/dice-ton-api/internal/utils"
)

func TestFetcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Fetcher Suite")
}

var (
	bets = []*models.Bet{
		{
			ID:     1,
			Status: models.BetStatusSent,
		},
		{
			ID:     2,
			Status: models.BetStatusSent,
		},
		{
			ID:     3,
			Status: models.BetStatusSent,
		},
		{
			ID:     4,
			Status: models.BetStatusSent,
		},
	}
)

var _ = Describe("Test Fetcher", func() {
	var (
		service *fetcher.Service
		err     error

		ctrl     *gomock.Controller
		mockDice *interfaces.MockDiceContract
		mockRepo *interfaces.MockRepository

		container = &dependencies.Container{
			Config: config.LoadConfig(),
			Logger: zap.NewExample(),
		}
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockDice = interfaces.NewMockDiceContract(ctrl)
		mockRepo = interfaces.NewMockRepository(ctrl)

		container.DiceContract = mockDice
		container.Repository = mockRepo

		service = fetcher.New(container)
	})

	Describe("Test Fetching Game Results", func() {
		JustBeforeEach(func() {
			err = service.Do()
		})

		var (
			chBets   chan *models.Bet
			copyBets []*models.Bet
		)

		BeforeEach(func() {
			chBets = make(chan *models.Bet)
			copyBets = utils.DupSliceOfPointers(bets)
		})

		When("", func() {
			BeforeEach(func() {
				mockRepo.EXPECT().GetLastResolvedBetLT().Return(uint64(0), nil)
				mockDice.EXPECT().SubscribeOnFinishedBets(context.Background(), uint64(0)).DoAndReturn(func(context.Context, uint64) (<-chan *models.Bet, error) {
					go func() {
						defer close(chBets)

						chBets <- bets[0]
						chBets <- bets[1]
						chBets <- bets[2]
						chBets <- bets[3]
					}()

					return chBets, nil
				})

				mockRepo.EXPECT().Update(copyBets[0]).Return(nil)
				mockRepo.EXPECT().Update(copyBets[1]).Return(nil)
				mockRepo.EXPECT().Update(copyBets[2]).Return(nil)
				mockRepo.EXPECT().Update(copyBets[3]).Return(nil)

				copyBets[0].Status = models.BetStatusResolved
				copyBets[1].Status = models.BetStatusResolved
				copyBets[2].Status = models.BetStatusResolved
				copyBets[3].Status = models.BetStatusResolved
			})

			It("should not have error occured", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
