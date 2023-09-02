package main

type InMemoryDataStore struct {
	blocks []Block

	// Index transactions by address
	transactionsByAddress map[string][]Transaction
}

func NewInMemoryDataStore() *InMemoryDataStore {
	return &InMemoryDataStore{
		transactionsByAddress: make(map[string][]Transaction),
	}
}

func (imds *InMemoryDataStore) GetCurrentBlock() Block {
	if len(imds.blocks) == 0 {
		return Block{}
	}
	return imds.blocks[len(imds.blocks)-1]
}

func (imds *InMemoryDataStore) AddBlock(block Block) {
	imds.blocks = append(imds.blocks, block)
	for _, transaction := range block.Transactions {
		imds.transactionsByAddress[transaction.From] = append(imds.transactionsByAddress[transaction.From], transaction)
		imds.transactionsByAddress[transaction.To] = append(imds.transactionsByAddress[transaction.To], transaction)
	}
}

func (imds *InMemoryDataStore) GetTransactions(address string) []Transaction {
	return imds.transactionsByAddress[address]
}
