package tonapi

import (
	"errors"
	"fmt"

	tonlib "github.com/mercuryoio/tonlib-go/v2"
	"github.com/tongineers/tonlib-go-api/config"
	"github.com/tongineers/tonlib-go-api/internal/dto"
	"github.com/tongineers/tonlib-go-api/internal/gateways/web/controllers/apiv1/transactions"
	"go.uber.org/zap"
)

var (
	_ transactions.TONClient = (*Service)(nil)
)

type (
	Service struct {
		client *tonlib.Client
		key    *tonlib.InputKey
		conf   *config.Config
		logger *zap.Logger
	}

	Opt func(s *Service)
)

func WithClient(c *tonlib.Client) Opt {
	return func(s *Service) {
		s.client = c
	}
}

func WithKey(k *tonlib.InputKey) Opt {
	return func(s *Service) {
		s.key = k
	}
}

func WithConfig(c *config.Config) Opt {
	return func(s *Service) {
		s.conf = c
	}
}

func WithLogger(l *zap.Logger) Opt {
	return func(s *Service) {
		s.logger = l
	}
}

func (s *Service) apply(opts ...Opt) {
	for _, opt := range opts {
		opt(s)
	}
}

func New(opts ...Opt) (*Service, error) {
	s := &Service{}
	s.apply(opts...)

	return s, nil
}

func (s *Service) GetTransactions(in *dto.GetTransactions) ([]*dto.Transaction, error) {
	resp, err := s.client.RawGetTransactions(*tonlib.NewAccountAddress(in.Addr), *tonlib.NewInternalTransactionId(in.Hash, tonlib.JSONInt64(in.Lt)), *s.key)
	if err != nil {
		// need to restart container
		//panic(err)
		//s.api.UpdateTonConnection()
		return nil, err
	}

	txns := make([]*dto.Transaction, 0)
	for _, trx := range resp.Transactions {
		msgData := trx.InMsg.MsgData.(map[string]interface{})
		msgDataText := ""
		if msgData["@type"] == "msg.dataText" {
			msgDataText = msgData["text"].(string)
		}
		inMsg := &dto.Message{
			BodyHash:    trx.InMsg.BodyHash,
			CreatedLt:   int64(trx.InMsg.CreatedLt),
			Destination: trx.InMsg.Destination.AccountAddress,
			FwdFee:      int64(trx.InMsg.FwdFee),
			IhrFee:      int64(trx.InMsg.IhrFee),
			Message:     msgDataText,
			Source:      trx.InMsg.Source.AccountAddress,
			Value:       int64(trx.InMsg.Value),
		}

		outMsgs := make([]*dto.Message, 0)
		for _, msg := range trx.OutMsgs {
			msgData := msg.MsgData.(map[string]interface{})
			msgDataText := ""
			if msgData["@type"] == "msg.dataText" {
				msgDataText = msgData["text"].(string)
			}

			outMsg := &dto.Message{
				BodyHash:    msg.BodyHash,
				CreatedLt:   int64(msg.CreatedLt),
				Destination: msg.Destination.AccountAddress,
				FwdFee:      int64(msg.FwdFee),
				IhrFee:      int64(msg.IhrFee),
				Message:     msgDataText,
				Source:      msg.Source.AccountAddress,
				Value:       int64(msg.Value),
			}
			outMsgs = append(outMsgs, outMsg)
		}

		txnID := &dto.TransactionID{
			Hash: trx.TransactionId.Hash,
			Lt:   int64(trx.TransactionId.Lt),
		}

		txn := &dto.Transaction{
			Data:       trx.Data,
			Fee:        int64(trx.Fee),
			InMsg:      inMsg,
			OtherFee:   int64(trx.OtherFee),
			OutMsgs:    outMsgs,
			StorageFee: int64(trx.StorageFee),
			TxnID:      txnID,
		}
		txns = append(txns, txn)
	}

	return txns, nil
}

func (s *Service) GetAccountState(in *dto.GetAccountState) (*dto.AccountState, error) {
	resp, err := s.client.RawGetAccountState(*tonlib.NewAccountAddress(in.Addr))
	if err != nil {
		// need to restart container
		//panic(err)
		//s.api.UpdateTonConnection()
		return nil, err
	}

	if !isRawFullAccountState(resp) {
		return nil, errors.New("Invalid return type for GetAccountState")
	}

	txnID := &dto.TransactionID{
		Hash: resp.LastTransactionId.Hash,
		Lt:   int64(resp.LastTransactionId.Lt),
	}

	return &dto.AccountState{
		Balance:           int64(resp.Balance),
		Code:              resp.Code,
		Data:              resp.Data,
		FrozenHash:        resp.FrozenHash,
		LastTransactionId: txnID,
		SyncUtime:         resp.SyncUtime,
	}, nil
}

// no longer in use
func (s *Service) GetBetSeed(in *dto.GetBetSeed) {

}

func (s *Service) GetActiveBets(in *dto.GetActiveBets) ([]*dto.Bet, error) {
	address := tonlib.NewAccountAddress(s.conf.TONContractAddress)
	smcInfo, err := s.client.SmcLoad(*address)
	if err != nil {
		//s.api.UpdateTonConnection()
		return nil, err
	}

	methodName := "active_bets"
	methodID := struct {
		Type  string `json:"@type"`
		Extra string `json:"@extra"`
		Name  string `json:"name"`
	}{
		Type: "smc.methodIdName",
		Name: methodName,
	}

	stack := make([]tonlib.TvmStackEntry, 0)
	res, err := s.runGetMethod(smcInfo.Id, methodID, stack)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("empty response")
	}

	return nil, nil

	// var entries CustomTvmStackEntry
	// asBytes, err := json.Marshal(res[0])
	// if err != nil {
	// 	return nil, err
	// }
	// err = json.Unmarshal(asBytes, &bets)
	// if err != nil {
	// 	return nil, err
	// }

	// var activeBets []*pb.ActiveBet
	// for _, element := range bets.List.Elements {
	// 	var betIdRaw CustomTvmStackEntryNumber
	// 	asBytes, err := json.Marshal(element.Tuple.Elements[0])
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	err = json.Unmarshal(asBytes, &betIdRaw)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	betId, err := strconv.Atoi(betIdRaw.Number.Number)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	other := element.Tuple.Elements[1]

	// 	var tmp _CustomTvmStackEntryTuple
	// 	asBytes, err = json.Marshal(other)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	err = json.Unmarshal(asBytes, &tmp)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	params := tmp.Tuple.Elements

	// 	rollUnder, err := strconv.Atoi(params[0].Number.Number)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	amount, err := strconv.Atoi(params[1].Number.Number)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	wc1, err := strconv.Atoi(params[2].Number.Number)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	address1 := params[3].Number.Number
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	wc2, err := strconv.Atoi(params[4].Number.Number)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	address2 := params[5].Number.Number
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	refBonus, err := strconv.Atoi(params[6].Number.Number)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	seed := params[7].Number.Number

	// 	bet := &pb.ActiveBet{
	// 		Id:            int32(betId),
	// 		RollUnder:     int32(rollUnder),
	// 		Amount:        int64(amount),
	// 		PlayerAddress: &pb.TonAddress{Workchain: int32(wc1), Address: address1},
	// 		RefAddress:    &pb.TonAddress{Workchain: int32(wc2), Address: address2},
	// 		RefBonus:      int64(refBonus),
	// 		Seed:          seed,
	// 	}

	// 	activeBets = append(activeBets, bet)
	// }

	// return &pb.GetActiveBetsResponse{
	// 	Bets: activeBets,
	// }, nil
}

// no longer in use
func (s *Service) GetSeqno(in *dto.GetSeqno) {

}

func (s *Service) SendMessage(in *dto.SendMessage) {

}

func (s *Service) runGetMethod(id int64, method interface{}, stack []tonlib.TvmStackEntry) ([]tonlib.TvmStackEntry, error) {
	resp, err := s.client.SmcRunGetMethod(id, method, stack)
	if err != nil {
		// need to restart container
		//panic(err)
		//s.api.UpdateTonConnection()
		return nil, err
	}

	return resp.Stack, nil
}
