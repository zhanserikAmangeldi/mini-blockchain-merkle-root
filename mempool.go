package main

import (
	"fmt"
	"time"
)

type MemPool struct {
	Transactions []Transaction
	TimeToLive   int64
}

func (mempool *MemPool) AddTransaction(tx Transaction) {
	if !tx.IsValid() {
		fmt.Println("Invalid transaction, not adding to mempool")
		return
	}

	for _, existingTx := range mempool.Transactions {
		if existingTx.Hash == tx.Hash {
			fmt.Println("Transaction already exists in mempool")
			return
		}
	}

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

	if limit == 0 {
		return []Transaction{}
	}

	txs := make([]Transaction, limit)
	copy(txs, mempool.Transactions[:limit])
	mempool.Transactions = mempool.Transactions[limit:]

	return txs
}
