package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

type EvmParser struct {
	ec                   EthClient
	ch                   chan any
	currentBlock         int
	nextBlock            int
	dataStore            DataStore
	subscriberStore      SubscriberStore
	notificationsService *NotificationsService
}

func NewEthParser(ec EthClient) *EvmParser {
	dataStore := NewInMemoryDataStore()
	block := dataStore.GetCurrentBlock()
	currentBlock := 0

	var nextBlock int
	if len(block.Number) == 0 {
		var err error
		nextBlock, err = ec.BlockNumber()
		if err != nil {
			log.Fatalf("failed to get latest block height: %v", err)
		}
	} else {
		currentBlockInt64, err := strconv.ParseInt(block.Number[:2], 16, 64)
		if err != nil {
			log.Fatalf("failed to get parse current block height: %v", err)
		}
		currentBlock = int(currentBlockInt64)
		nextBlock = currentBlock + 1
	}

	notificationsService := NewNotificationsService()
	go notificationsService.Start()

	return &EvmParser{
		ec:                   ec,
		ch:                   make(chan any),
		currentBlock:         currentBlock,
		nextBlock:            nextBlock,
		dataStore:            dataStore,
		subscriberStore:      NewInMemorySubscriberStore(),
		notificationsService: notificationsService,
	}
}

// Start parsing the blockchain from the last block parsed or the latest block.
// This is a blocking call.
func (ep *EvmParser) Start() {
	fmt.Printf("Parser started!\n")

	parseNextBlock := func() {
		block, err := ep.ec.GetBlockByNumber(ep.nextBlock)
		if err != nil {
			log.Fatalf("failed to get block at height %v: %v", ep.nextBlock, err)
		}

		if len(block.Number) != 0 {
			ep.dataStore.AddBlock(block)
			ep.currentBlock = ep.nextBlock
			ep.nextBlock += 1

			// Handle subscribers
			for _, transaction := range block.Transactions {
				if ep.subscriberStore.HasSubscriber(transaction.From) {
					ep.notificationsService.fromCh <- transaction
				}
				if ep.subscriberStore.HasSubscriber(transaction.To) {
					ep.notificationsService.toCh <- transaction
				}
			}
		}
		fmt.Printf("Waiting for next block...\n")
	}
	parseNextBlock()

	ticker := time.NewTicker(15 * time.Second)
	for {
		select {
		case <-ticker.C:
			parseNextBlock()
		case <-ep.ch:
			return
		}
	}
}

func (ep *EvmParser) Stop() {
	close(ep.ch)
}

func (ep *EvmParser) GetCurrentBlock() int {
	return ep.currentBlock
}

func (ep *EvmParser) Subscribe(address string) bool {
	if ep.subscriberStore.HasSubscriber(address) {
		return false
	}

	ep.subscriberStore.AddSubscriber(address)
	return true
}

func (ep *EvmParser) GetTransactions(address string) []Transaction {
	return ep.dataStore.GetTransactions(address)
}
