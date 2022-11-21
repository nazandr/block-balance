package main

import (
	"block-balance/internal/client"
	"block-balance/internal/service"
	"fmt"
	"os"

	log "github.com/go-pkgz/lgr"
	"github.com/umputun/go-flags"
)

var opts struct {
	BlockAmount int    `short:"b" long:"block-amount" env:"BLOCK_AMOUNT" default:"100" description:"Amount of blocks to scan"`
	APIKye      string `short:"a" long:"api-key" env:"API_KEY" description:"Getblock.io API key"`
	Debug       bool   `long:"dbg" env:"DEBUG" description:"Enable debug mode"`
}

func main() {
	p := flags.NewParser(&opts, flags.PassDoubleDash|flags.HelpFlag)
	if _, err := p.Parse(); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		p.WriteHelp(os.Stderr)
		os.Exit(2)
	}

	setupLog(opts.Debug)

	client := client.NewClient("https://eth.getblock.io/mainnet/", "getblock.io", opts.APIKye)

	service := service.NewService(client, opts.BlockAmount)

	address, amount, err := service.GetMostChangedAccount()
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}

	fmt.Printf("Address, which balance has changed (in any direction) more than the others over the last %d blocks:\n %s : %v Wei\n",
		opts.BlockAmount, address, amount)
}

func setupLog(dbg bool) {
	logOpts := []log.Option{log.Msec, log.LevelBraces, log.StackTraceOnError}
	if dbg {
		logOpts = append(logOpts, log.Debug, log.CallerFunc)
	}
	log.SetupStdLogger(logOpts...)
	log.Setup(logOpts...)
}
