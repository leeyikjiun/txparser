package main

type DataStore interface {
	GetCurrentBlock() Block
	AddBlock(block Block)
	GetTransactions(address string) []Transaction
}
