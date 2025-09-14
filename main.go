package main

import "time"

type Transaction struct {
	Hash      string
	Amount    float64
	From      string
	To        string
	Timestamp int64
}

type Block struct {
	BlockHeight  int
	Hash         string
	PrevHash     string
	Timestamp    int64
	Transactions []Transaction
	MerleRoot    string
}

type Blockchain struct {
	Blocks []Block
}

type MemPool struct {
	Transactions []Transaction
	TimeToLive   int64
}

func (mempool *MemPool) AddTransaction(tx Transaction) {
	mempool.Transactions = append(mempool.Transactions, tx)
}

func (mempool *MemPool) CleanUp() {
	now := time.Now().Unix()
	newTxs := []Transaction{}
	for _, tx := range mempool.Transactions {
		if now-tx.Timestamp <= mempool.TimeToLive {
			newTxs = append(newTxs, tx)
		}
	}

	mempool.Transactions = newTxs
}

func (mempool *MemPool) GetTransactions(limit int) []Transaction {
	mempool.CleanUp()
	if len(mempool.Transactions) < limit {
		limit = len(mempool.Transactions)
	}
	txs := mempool.Transactions[:limit]
	mempool.Transactions = mempool.Transactions[limit:]

	return txs
}

func (blockchain *Blockchain) AddBlock(newBlock Block) {
	blockchain.Blocks = append(blockchain.Blocks, newBlock)
}

func main() {
}
