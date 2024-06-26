package main

import (
	"block-balance/internal/client"
	"block-balance/internal/service"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

var opts struct {
	BlockAmount int
	APIKey      string
	RPS         int
	NumWorkers  int
	Timeout     time.Duration
}

func main() {
	// Read environment variables
	blockAmountStr := os.Getenv("BLOCK_AMOUNT")
	apiKey := os.Getenv("API_KEY")
	rpsStr := os.Getenv("RPS")
	numWorkersStr := os.Getenv("NUM_WORKERS")
	timeoutStr := os.Getenv("TIMEOUT")

	opts.BlockAmount = convertToInt(blockAmountStr, "100")
	opts.APIKey = apiKey
	opts.RPS = convertToInt(rpsStr, "10")
	opts.NumWorkers = convertToInt(numWorkersStr, "10")

	ctx := context.Background()
	if timeoutStr != "" {
		var err error
		opts.Timeout, err = time.ParseDuration(timeoutStr)
		if err != nil {
			log.Fatalf("[ERROR] Error parsing TIMEOUT: %v\n", err)
		}
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	client, err := client.NewClient("https://go.getblock.io", "getblock.io", opts.APIKey)
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}

	service := service.NewService(client, opts.BlockAmount, opts.RPS, opts.NumWorkers)

	start := time.Now()
	address, amount, err := service.GetMostChangedAccount(ctx)
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}

	fmt.Printf("Address, which balance has changed (in any direction) more than the others over the last %d blocks:\n%s : %v Wei\n",
		opts.BlockAmount, address, amount)
	fmt.Printf("Time elapsed: %v\n", time.Since(start))
}

func convertToInt(s, defaultValue string) int {
	if s == "" {
		s = defaultValue
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("[ERROR] Error converting %s to int: %v\n", s, err)
	}
	return i
}
