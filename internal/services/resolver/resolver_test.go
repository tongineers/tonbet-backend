package resolver_test

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"github.com/tongineers/tonbet-backend/config"
	"github.com/tongineers/tonbet-backend/internal/app/dependencies"
	"github.com/tongineers/tonbet-backend/internal/interfaces"
	"github.com/tongineers/tonbet-backend/internal/models"
	"github.com/tongineers/tonbet-backend/internal/services/resolver"
	"github.com/tongineers/tonbet-backend/internal/utils"
)

func TestResolver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Resolver Suite")
}

var (
	bets = []*models.Bet{
		{
			ID:     1,
			Seed:   "F27991EB4857143790D0FBDA632F163A0B28FF5F3E00CC76567D054622094CB0",
			Status: models.BetStatusNew,
		},
		{
			ID:     2,
			Seed:   "A5D167C647F1FB1ECCE093F442DAFC6C942805288083BA5F27E1AAB90D23222B",
			Status: models.BetStatusNew,
		},
		{
			ID:     3,
			Seed:   "09BCF49F3DEA91774898A50BB4EF6A1BE7BD4A3C4247ED7B10EA6EE29AAF64A6",
			Status: models.BetStatusNew,
		},
		{
			ID:     4,
			Seed:   "6A04E81073FCBE2926BC08085B49256F79C62CFC2A54BC08BCFE037FE130B113",
			Status: models.BetStatusNew,
		},
	}

	errSomethingWentWrong = errors.New("something went wrong")
)

var _ = Describe("Test Resolver", func() {
	var (
		service *resolver.Service
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

		service = resolver.New(container)
	})

	Describe("Test Resolving New Bets", func() {
		JustBeforeEach(func() {
			err = service.Do()
		})

		var (
			copyBets = make([]*models.Bet, 0)
		)

		BeforeEach(func() {
			copyBets = utils.DupSliceOfPointers(bets)
		})

		When("all bets are resolved", func() {
			BeforeEach(func() {
				mockRepo.EXPECT().ReadByStatus(models.BetStatusNew).Return(bets, nil)
				mockDice.EXPECT().ResolveBet(bets[0].ID, bets[0].Seed).Return(nil)
				mockDice.EXPECT().ResolveBet(bets[1].ID, bets[1].Seed).Return(nil)
				mockDice.EXPECT().ResolveBet(bets[2].ID, bets[2].Seed).Return(nil)
				mockDice.EXPECT().ResolveBet(bets[3].ID, bets[3].Seed).Return(nil)
				mockRepo.EXPECT().Update(copyBets).Return(nil)

				copyBets[0].Status = models.BetStatusSent
				copyBets[1].Status = models.BetStatusSent
				copyBets[2].Status = models.BetStatusSent
				copyBets[3].Status = models.BetStatusSent
			})

			It("should not have error occured", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("resolving failed for some bets", func() {
			BeforeEach(func() {
				mockRepo.EXPECT().ReadByStatus(models.BetStatusNew).Return(bets, nil)
				mockDice.EXPECT().ResolveBet(bets[0].ID, bets[0].Seed).Return(nil)
				mockDice.EXPECT().ResolveBet(bets[1].ID, bets[1].Seed).Return(errSomethingWentWrong)
				mockDice.EXPECT().ResolveBet(bets[2].ID, bets[2].Seed).Return(nil)
				mockDice.EXPECT().ResolveBet(bets[3].ID, bets[3].Seed).Return(errSomethingWentWrong)
				mockRepo.EXPECT().Update([]*models.Bet{copyBets[0], copyBets[2]}).Return(nil)

				copyBets[0].Status = models.BetStatusSent
				copyBets[2].Status = models.BetStatusSent
			})

			It("should not have error occured", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
