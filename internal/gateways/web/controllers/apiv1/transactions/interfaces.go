package transactions

import "github.com/tongineers/tonlib-go-api/internal/dto"

type (
	TONClient interface {
		GetTransactions(in *dto.GetTransactions) ([]*dto.Transaction, error)
		GetAccountState(in *dto.GetAccountState) (*dto.AccountState, error)
		GetActiveBets(in *dto.GetActiveBets) ([]*dto.Bet, error)
	}
)
