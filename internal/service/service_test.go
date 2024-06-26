package service

import (
	"block-balance/internal/client"
	"block-balance/internal/types"
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

type testClient struct {
	err bool
}

func (tc *testClient) GetLastBlockNumber() (int64, error) {
	return 10, nil
}

func (tc *testClient) GetBlockTransactionsByNumber(number string) ([]types.Transaction, error) {
	if !tc.err {
		tc.err = true
		return nil, fmt.Errorf("error")
	}
	txs := []types.Transaction{
		{
			From:     "a",
			To:       "b",
			Value:    "0",
			Gas:      "0",
			GasPrice: "0",
		},
		{
			From:     "a",
			To:       "b",
			Value:    "64",
			Gas:      "0",
			GasPrice: "0",
		},
		{
			From:     "a",
			To:       "c",
			Value:    "64",
			Gas:      "0",
			GasPrice: "0",
		},
		{
			From:     "b",
			To:       "d",
			Value:    "64",
			Gas:      "0",
			GasPrice: "0",
		},
	}
	return txs, nil
}

func TestService_GetMostChangedAccount(t *testing.T) {
	type fields struct {
		client      client.Client
		addresses   *types.Addresses
		blockAmount int
		wg          *sync.WaitGroup
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		want1   *big.Int
		wantErr bool
	}{
		{
			name: "one block",
			fields: fields{
				client:      &testClient{err: true},
				addresses:   types.NewAddresses(),
				blockAmount: 1,
				wg:          &sync.WaitGroup{},
			},
			want:    "a",
			want1:   big.NewInt(200),
			wantErr: false,
		},
		{
			name: "three blocks",
			fields: fields{
				client:      &testClient{err: false},
				addresses:   types.NewAddresses(),
				blockAmount: 3,
				wg:          &sync.WaitGroup{},
			},
			want:    "a",
			want1:   big.NewInt(400),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				client:      tt.fields.client,
				addresses:   tt.fields.addresses,
				blockAmount: tt.fields.blockAmount,
				numWorkers:  2,
				wg:          tt.fields.wg,
				rateLimiter: *rate.NewLimiter(rate.Every(time.Second/10), 10),
			}

			got, got1, err := s.GetMostChangedAccount(context.Background())
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
