package smartcont

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/tvm/cell"

	"github.com/tongineers/dice-ton-api/config"
	"github.com/tongineers/dice-ton-api/internal/models"
)

type (
	Service struct {
		api  *ton.APIClient
		conf *config.Config
	}
)

const (
	ResolveQueryFileName = "resolve-query.fif"
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

func New(
	api *ton.APIClient,
	conf *config.Config,
) *Service {
	return &Service{
		api:  api,
		conf: conf,
	}
}

// GetAccountState returns account state by the specific address.
func (s *Service) GetAccountState(accountAddr string) (*models.AccountState, error) {
	master, err := s.api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		return nil, err
	}

	addr, err := address.ParseAddr(accountAddr)
	if err != nil {
		return nil, err
	}

	res, err := s.api.WaitForBlock(master.SeqNo).GetAccount(context.Background(), master, addr)
	if err != nil {
		return nil, err
	}

	return &models.AccountState{
		Balance:  res.State.Balance.Nano().Int64(),
		Data:     res.Data.Dump(),
		LastHash: string(res.LastTxHash),
		LastLt:   res.LastTxLT,
	}, nil
}

// SubscribeOnFinishedBets returns (asynchronously) game results starts from the specific LT.
func (s *Service) SubscribeOnFinishedBets(ctx context.Context, fromLT uint64) (<-chan *models.Bet, error) {
	addr, err := address.ParseAddr(s.conf.TONContractAddr)
	if err != nil {
		return nil, err
	}

	transactions := make(chan *tlb.Transaction)
	go s.api.SubscribeOnTransactions(ctx, addr, 28136901000003, transactions)

	bets := make(chan *models.Bet)
	go func() {
		for tx := range transactions {
			// only external messages can contains the game result
			if tx.IO.In != nil && tx.IO.In.MsgType == tlb.MsgTypeExternalIn {
				if tx.OutMsgCount != 1 {
					continue
				}

				out, err := tx.IO.Out.ToSlice()
				if err != nil {
					continue
				}

				external := out[0].AsInternal()
				betID, randomRoll, err := parseOutMessage(external.Comment())
				if err != nil {
					continue
				}

				encodedHash := base64.StdEncoding.EncodeToString(tx.Hash)

				bets <- &models.Bet{
					ID:         betID,
					RandomRoll: randomRoll,
					Payout:     external.Amount.Nano().Uint64(),
					LastHash:   encodedHash,
					LastLT:     tx.LT,
				}
			}
		}
	}()

	return bets, nil
}

// GetActiveBets returns active bets via running smartcontract GET method.
func (s *Service) GetActiveBets() ([]*models.Bet, error) {
	master, err := s.api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		return nil, err
	}

	addr, err := address.ParseAddr(s.conf.TONContractAddr)
	if err != nil {
		return nil, err
	}

	res, err := s.api.WaitForBlock(master.SeqNo).RunGetMethod(context.Background(), master, addr, "active_bets")
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`[^0-9 ]`)
	bets := make([]*models.Bet, 0)
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

		amount, err := strconv.ParseUint(data[Amount], 10, 64)
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

		refBonus, err := strconv.ParseUint(data[RefBonus], 10, 64)
		if err != nil {
			return nil, err
		}

		seed := data[Seed]

		bets = append(bets, &models.Bet{
			ID:            betID,
			RollUnder:     rollUnder,
			Amount:        amount,
			PlayerAddress: playerAddr,
			RefAddress:    refAddr,
			RefBonus:      refBonus,
			Seed:          seed,
			Status:        models.BetStatusNew,
		})
	}

	return bets, nil
}

// ResolveBet resolving active bet with specific betID and seed.
func (s *Service) ResolveBet(betID int, seed string) error {
	fileNameWithPath := ResolveQueryFileName
	fileNameStart := strings.LastIndex(fileNameWithPath, "/")
	fileName := fileNameWithPath[fileNameStart+1:]

	bocFile := strings.Replace(fileName, ".fif", ".boc", 1)
	_ = os.Remove(bocFile)

	var out bytes.Buffer
	cmd := exec.Command("fift", "-s", fileNameWithPath, s.conf.TONOwnerKeyFile, s.conf.TONContractAddr, strconv.Itoa(betID), seed)
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return err
	}

	f, err := os.Open(bocFile)
	if err != nil {
		return err
	}
	defer f.Close()

	b := &bytes.Buffer{}
	io.Copy(b, f)

	return s.sendMessage(b.Bytes())
}

func (s *Service) sendMessage(body []byte) error {
	master, err := s.api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		return err
	}

	addr, err := address.ParseAddr(s.conf.TONContractAddr)
	if err != nil {
		return err
	}

	cell, err := cell.FromBOC(body)
	if err != nil {
		return err
	}

	msg := &tlb.ExternalMessage{
		DstAddr: addr,
		Body:    cell,
	}

	err = s.api.WaitForBlock(master.SeqNo).SendExternalMessage(context.Background(), msg)
	if err != nil {
		return err
	}

	return nil
}

func parseOutMessage(msg string) (int, int, error) {
	r, _ := regexp.Compile(`TONBET.IO - \[#(\d+)] Your number is (\d+), all numbers greater than (\d+) have won.`)
	matches := r.FindStringSubmatch(string(msg))

	if len(matches) > 0 {
		betID, _ := strconv.Atoi(matches[1])
		randomRoll, _ := strconv.Atoi(matches[3])

		return betID, randomRoll, nil
	}

	return 0, 0, fmt.Errorf("message does not match expected pattern")
}
