package client

import (
	"block-balance/internal/types"
	"block-balance/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client interface {
	GetLastBlockNumber() (int64, error)
	GetBlockTransactionsByNumber(number string) ([]types.Transaction, error)
}

type client struct {
	address   string
	requestId string
	apiKey    string
}

func NewClient(address, RequestId, apiKey string) (Client, error) {
	addr, err := url.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("failed to parse address: %w", err)
	}
	addr.Path = fmt.Sprintf("%s/%s", addr.Path, apiKey)

	return &client{
		address:   addr.String(),
		requestId: RequestId,
	}, nil
}

func (c *client) GetLastBlockNumber() (int64, error) {
	res, err := c.rpcRequest("eth_blockNumber", "")
	if err != nil {
		return 0, fmt.Errorf("failed to get request: %w", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read body: %w", err)
	}

	var LastBlock types.BlockNumber
	if err := json.Unmarshal(body, &LastBlock); err != nil {
		return 0, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return utils.HexToInt(LastBlock.Result)
}

func (c *client) GetBlockTransactionsByNumber(number string) ([]types.Transaction, error) {
	res, err := c.rpcRequest("eth_getBlockByNumber", fmt.Sprintf(`"%s", true`, number))
	if err != nil {
		return nil, fmt.Errorf("failed to get request: %w", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	var blockData types.BlockData
	if err := json.Unmarshal(body, &blockData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return blockData.Result.Transactions, nil
}

func (c *client) rpcRequest(rpcMethod, params string) (*http.Response, error) {
	payload := strings.NewReader(fmt.Sprintf(`{
    "jsonrpc": "2.0",
    "method": "%s",
    "params": [%s],
    "id": "%s"
}`, rpcMethod, params, c.requestId))

	req, err := http.NewRequest(http.MethodGet, c.address, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do http request: %w", err)
	}

	return res, nil
}
