package main

import (
	"fmt"
	"time"
)

func main() {
	ethClient := NewEthClient("https://cloudflare-eth.com")
	ethParser := NewEthParser(ethClient)

	// Uniswap v2 router 2 address. Should be able to trigger transactions very fast.
	address := "0x7a250d5630b4cf539739df2c5dacb4c659f2488d"
	ethParser.Subscribe(address)
	go func() {
		<-time.After(time.Minute)
		// Print first 5 transactions
		for i, transaction := range ethParser.GetTransactions(address)[:5] {
			fmt.Printf("%v: %+v\n", i, transaction)
		}
		ethParser.Stop()
	}()
	ethParser.Start()
}
