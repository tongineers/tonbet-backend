package tonapi

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	_ "github.com/xssnick/tonutils-go/tvm/cell"
	"go.uber.org/zap"

	"github.com/tongineers/dice-ton-api/config"
	pb "github.com/tongineers/dice-ton-api/gen/go/tonapi/v1"
	appgo "github.com/tongineers/dice-ton-api/pkg/app-go"
)

var (
	_ appgo.GRPCService = (*Service)(nil)
)

type (
	Service struct {
		client *ton.APIClient
		conf   *config.Config
		logger *zap.Logger
		pb.UnimplementedTonApiServiceServer
	}

	Opt func(s *Service)
)

const (
	BetID = iota
	RollUnder
	Amount
	PlayerWorkchain
	PlayerAddress
	RefWorkchain
	RefAddress
	RefBonus
	Seed
)

func WithClient(c *ton.APIClient) Opt {
	return func(s *Service) {
		s.client = c
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

func New(opts ...Opt) (*Service, error) {
	s := &Service{}
	s.apply(opts...)

	return s, nil
}

func (s *Service) apply(opts ...Opt) {
	for _, opt := range opts {
		opt(s)
	}
}

func (s *Service) ServiceDef() *appgo.GRPCOptions {
	return &appgo.GRPCOptions{
		Handler:     pb.RegisterTonApiServiceHandler,
		ServiceDesc: &pb.TonApiService_ServiceDesc,
		ServiceImpl: s,
	}
}

func (s *Service) FetchTransactions(ctx context.Context, in *pb.FetchTransactionsRequest) (*pb.FetchTransactionsResponse, error) {
	freshBlock, err := s.client.CurrentMasterchainInfo(context.Background())
	if err != nil {
		return nil, err
	}

	addr, err := address.ParseAddr(in.GetAddress())
	if err != nil {
		return nil, err
	}

	hash, err := base64.StdEncoding.DecodeString(in.GetHash())
	if err != nil {
		return nil, err
	}

	res, err := s.client.WaitForBlock(freshBlock.SeqNo).ListTransactions(context.Background(), addr, 100, uint64(in.GetLt()), hash)
	if err != nil {
		return nil, err
	}

	txns := make([]*pb.Transaction, 0)
	for _, txn := range res {
		if txn.IO.In.MsgType != tlb.MsgTypeInternal {
			continue
		}

		msgData := txn.IO.In.AsInternal()
		inMsg := &pb.RawMessage{
			Source:      msgData.SrcAddr.String(),
			Destination: msgData.DstAddr.String(),
			Value:       msgData.Amount.Nano().Int64(),
			FwdFee:      msgData.FwdFee.Nano().Int64(),
			IhrFee:      msgData.IHRFee.Nano().Int64(),
			Message:     msgData.Comment(),
			//BodyHash:    string(msgData.Body.Hash()),
			CreatedLt: int64(msgData.CreatedLT),
		}

		outMsgs := make([]*pb.RawMessage, 0)
		if txn.OutMsgCount > 0 {
			msgs, err := txn.IO.Out.ToSlice()
			if err != nil {
				continue
			}

			for _, msg := range msgs {
				msgData := msg.AsInternal()
				outMsgs = append(outMsgs, &pb.RawMessage{
					Source:      msgData.SrcAddr.String(),
					Destination: msgData.DstAddr.String(),
					Value:       msgData.Amount.Nano().Int64(),
					FwdFee:      msgData.FwdFee.Nano().Int64(),
					IhrFee:      msgData.IHRFee.Nano().Int64(),
					Message:     msgData.Comment(),
					//BodyHash:    string(msgData.Body.Hash()),
					CreatedLt: int64(msgData.CreatedLT),
				})

				// bs := msg.Msg.Payload().BitsSize()
				// s := msg.Msg.Payload().BeginParse()
				// b, err := s.LoadSlice(bs)
				// if err != nil {
				// 	panic(err)
				// }

				// comment += string(b)

				// s, err = s.LoadRef()
				// if err != nil {
				// 	panic(err)
				// }

				// bs = s.BitsLeft()
				// b, err = s.LoadSlice(bs)
				// if err != nil {
				// 	panic(err)
				// }

				// comment += string(b)
			}
		}

		txns = append(txns, &pb.Transaction{
			TransactionId: &pb.InternalTransactionId{
				//Hash: string(txn.Hash),
				Lt: int64(txn.LT),
			},
			//Data:    txn.Dump(),
			InMsg:   inMsg,
			OutMsgs: outMsgs,
			Fee:     txn.TotalFees.Coins.Nano().Int64(),
		})
	}

	fmt.Println(txns)

	return &pb.FetchTransactionsResponse{
		Items: txns,
	}, nil
}

func (s *Service) GetAccountState(ctx context.Context, in *pb.GetAccountStateRequest) (*pb.GetAccountStateResponse, error) {
	freshBlock, err := s.client.CurrentMasterchainInfo(context.Background())
	if err != nil {
		return nil, err
	}

	addr, err := address.ParseAddr(in.GetAccountAddress())
	if err != nil {
		return nil, err
	}

	res, err := s.client.WaitForBlock(freshBlock.SeqNo).GetAccount(context.Background(), freshBlock, addr)
	if err != nil {
		return nil, err
	}

	transactionId := &pb.InternalTransactionId{
		Hash: string(res.LastTxHash),
		Lt:   int64(res.LastTxLT),
	}

	return &pb.GetAccountStateResponse{
		Balance:           res.State.Balance.Nano().Int64(),
		Code:              res.Code.Dump(),
		Data:              res.Data.Dump(),
		LastTransactionId: transactionId,
	}, nil
}

func (s *Service) GetActiveBets(ctx context.Context, in *pb.GetActiveBetsRequest) (*pb.GetActiveBetsResponse, error) {
	freshBlock, err := s.client.CurrentMasterchainInfo(context.Background())
	if err != nil {
		return nil, err
	}

	addr, err := address.ParseAddr(s.conf.TONContractAddress)
	if err != nil {
		return nil, err
	}

	res, err := s.client.RunGetMethod(context.Background(), freshBlock, addr, "active_bets")
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`[^0-9 ]`)
	bets := make([]*pb.ActiveBet, 0)
	for _, item := range res.MustTuple(0) {
		item = re.ReplaceAllString(fmt.Sprintf("%s", item), "")
		data := strings.Split(item.(string), " ")

		betID, err := strconv.Atoi(data[BetID])
		if err != nil {
			return nil, err
		}

		rollUnder, err := strconv.Atoi(data[RollUnder])
		if err != nil {
			return nil, err
		}

		amount, err := strconv.Atoi(data[Amount])
		if err != nil {
			return nil, err
		}

		wc1, err := strconv.Atoi(data[PlayerWorkchain])
		if err != nil {
			return nil, err
		}
		playerAddr, err := toHumanRepresentationAddr(int8(wc1), data[PlayerAddress])
		if err != nil {
			return nil, err
		}

		wc2, err := strconv.Atoi(data[RefWorkchain])
		if err != nil {
			return nil, err
		}
		refAddr, err := toHumanRepresentationAddr(int8(wc2), data[RefAddress])
		if err != nil {
			return nil, err
		}

		refBonus, err := strconv.Atoi(data[RefBonus])
		if err != nil {
			return nil, err
		}

		seed := data[Seed]

		bets = append(bets, &pb.ActiveBet{
			Id:            int32(betID),
			RollUnder:     int32(rollUnder),
			Amount:        int64(amount),
			PlayerAddress: &pb.TonAddress{Workchain: int32(wc1), Address: playerAddr},
			RefAddress:    &pb.TonAddress{Workchain: int32(wc2), Address: refAddr},
			RefBonus:      int64(refBonus),
			Seed:          seed,
		})
	}

	return &pb.GetActiveBetsResponse{
		Bets: bets,
	}, nil
}

func (s *Service) SendMessage(ctx context.Context, in *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	return &pb.SendMessageResponse{}, nil
}
