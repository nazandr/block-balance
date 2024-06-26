package service

import (
	"block-balance/internal/client"
	"block-balance/internal/types"
	"block-balance/internal/utils"
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"

	"golang.org/x/time/rate"
)

type Service struct {
	client client.Client

	addresses   *types.Addresses
	blockAmount int

	numWorkers int

	wg *sync.WaitGroup
	// Используем пакет rate из репозитория golang.org/x, рекомендация с go.dev. Можно заменить time.Ticker.
	// Считаю это уместным, учитывая требования по использованию stdlib, так как golang.org/x
	// это репозитории являющиеся частью проекта Go, но не входящий в основное дерево Go.
	// Они разрабатываются в соответствии с более слабыми требованиями к совместимости, чем ядро Go.
	rateLimiter rate.Limiter
}

func NewService(client client.Client, blockAmount, rps, numWorkers int) *Service {
	return &Service{
		client:      client,
		addresses:   types.NewAddresses(),
		blockAmount: blockAmount,
		numWorkers:  numWorkers,
		wg:          &sync.WaitGroup{},
		rateLimiter: *rate.NewLimiter(rate.Limit(rps), 60),
	}
}

func (s *Service) GetMostChangedAccount(ctx context.Context) (string, *big.Int, error) {
	lbNumber, err := s.client.GetLastBlockNumber()
	if err != nil {
		return "", nil, fmt.Errorf("failed to get last block number: %w", err)
	}

	blockNumCh := make(chan int, s.blockAmount)
	txCh := make(chan []types.Transaction, s.blockAmount)

	nw := s.numWorkers
	if s.blockAmount < s.numWorkers {
		nw = s.blockAmount
	}

	for w := 0; w < nw; w++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			for blockNum := range blockNumCh {
				select {
				case <-ctx.Done():
					return
				default:
					if err := s.rateLimiter.Wait(ctx); err != nil {
						log.Printf("[ERROR] rate limiter wait failed: %s", err)
						return
					}

					transactions, err := s.client.GetBlockTransactionsByNumber(utils.IntToHex(blockNum))
					if err != nil {
						log.Printf("[WARN] failed to get transactions of block %d: %s", blockNum, err)
						continue
					}
					txCh <- transactions
				}
			}
		}()
	}

	go func() {
		for blockNum := int(lbNumber); blockNum > int(lbNumber)-s.blockAmount; blockNum-- {
			blockNumCh <- blockNum
		}
		close(blockNumCh)
	}()

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

			gas, err := utils.HexToBigInt(tx.Gas)
			if err != nil {
				log.Printf("[WARN] %s", err)
				continue
			}

			gasPrice, err := utils.HexToBigInt(tx.GasPrice)
			if err != nil {
				log.Printf("[WARN] %s", err)
				continue
			}

			val = val.Add(val, gas.Mul(gas, gasPrice))
			s.addresses.Store(tx.From, tx.To, val)
		}
	}

	a, v := s.addresses.FindLargest()
	return a, v, nil
}
