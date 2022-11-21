package service

import (
	"block-balance/internal/client"
	"block-balance/internal/types"
	"block-balance/internal/utils"
	"log"
	"math/big"
	"sync"
)

type Service struct {
	client client.Client

	addresses   *types.Addresses
	blockAmount int

	wg *sync.WaitGroup
}

func NewService(client client.Client, blockAmount int) *Service {
	return &Service{
		client:      client,
		addresses:   types.NewAddresses(),
		blockAmount: blockAmount,
		wg:          &sync.WaitGroup{},
	}
}

func (s *Service) GetMostChangedAccount() (string, *big.Int, error) {
	lbNumber, err := s.client.GetLastBlockNumber()
	if err != nil {
		return "", nil, err
	}

	txCh := make(chan []types.Transaction, s.blockAmount)

	s.wg.Add(s.blockAmount)
	for blockNum := int(lbNumber); blockNum > int(lbNumber)-s.blockAmount; blockNum-- {
		go func(blockNum int) {
			defer s.wg.Done()
			transactions, err := s.client.GetBlockTransactionsByNumber(utils.IntToHex(blockNum))
			if err != nil {
				log.Printf("[WARN] failed to get transactions of block %d: %s", blockNum, err)
				return
			}
			txCh <- transactions
		}(blockNum)
	}

	go func() {
		s.wg.Wait()
		close(txCh)
	}()

	for txs := range txCh {
		for _, tx := range txs {
			if tx.Value == "0x0" || tx.Value == "0" {
				continue
			}

			val, err := utils.HexToBigInt(tx.Value)
			if err != nil {
				log.Printf("[WARN] %s", err)
				continue
			}
			s.addresses.Store(tx.From, tx.To, val)
		}
	}

	a, v := s.addresses.FindLargest()
	return a, v, nil
}
