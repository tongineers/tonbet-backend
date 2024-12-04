package listener_test

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"github.com/tongineers/dice-ton-api/config"
	"github.com/tongineers/dice-ton-api/internal/app/dependencies"
	"github.com/tongineers/dice-ton-api/internal/interfaces"
	"github.com/tongineers/dice-ton-api/internal/models"
	"github.com/tongineers/dice-ton-api/internal/services/listener"
)

func TestListener(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Listener Suite")
}

var (
	bets = []*models.Bet{
		{
			ID: 1,
		},
		{
			ID: 2,
		},
		{
			ID: 3,
		},
		{
			ID: 4,
		},
	}
)

var _ = Describe("Test Listener", func() {
	var (
		service *listener.Service
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

		service = listener.New(container)
	})

	Describe("Test Listening Active Bets", func() {
		JustBeforeEach(func() {
			err = service.Do()
		})

		When("have new active bets", func() {
			BeforeEach(func() {
				mockDice.EXPECT().GetActiveBets().Return(bets, nil)
				mockRepo.EXPECT().ReadByIDs([]int{1, 2, 3, 4}).Return(bets[:2], nil)
				mockRepo.EXPECT().Update(bets[2:]).Return(nil)
			})

			It("should not have error occured", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("active bets found but don't have new bets", func() {
			BeforeEach(func() {
				mockDice.EXPECT().GetActiveBets().Return(bets, nil)
				mockRepo.EXPECT().ReadByIDs([]int{1, 2, 3, 4}).Return(bets, nil)
			})

			It("should not have error occured", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("active bets not found", func() {
			BeforeEach(func() {
				mockDice.EXPECT().GetActiveBets().Return([]*models.Bet{}, nil)
			})

			It("should not have error occured", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
