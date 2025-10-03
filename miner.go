package main

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/exp/rand"
)

func mine(lastBlock *Block, txs []Transaction, difficulty int) Block {
	nonce := 0
	var hash string

	newBlock := Block{
		BlockHeight:  lastBlock.BlockHeight + 1,
		PrevHash:     lastBlock.Hash,
		Transactions: txs,
		Timestamp:    time.Now().Unix(),
		Difficulty:   difficulty,
	}
	newBlock.ComputeMerkleRoot()

	for {
		newBlock.Nonce = nonce
		hash = newBlock.computeHashInternal()

		if strings.HasPrefix(hash, strings.Repeat("0", difficulty)) {
			newBlock.Hash = hash
			break
		}
		nonce++
	}

	return newBlock
}

func createBlock(lastBlock *Block, mempool *MemPool, miner string, difficulty int) Block {
	rewardTx := Transaction{
		From:      "COINBASE",
		To:        miner,
		Amount:    1,
		Timestamp: time.Now().Unix(),
	}

	rewardTx.Hash = rewardTx.ComputeHash()

	txs := mempool.GetTransactions(len(mempool.Transactions))

	txs = append([]Transaction{rewardTx}, txs...)

	return mine(lastBlock, txs, difficulty)
}

func proofOfMine(mempool *MemPool, blockChan chan<- Block, blockchain *Blockchain) {
	miner := "Zhanserik"

	for {
		lastBlock := blockchain.GetLastBlock()
		if lastBlock == nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		block := createBlock(lastBlock, mempool, miner, blockchain.CurrentDifficulty)

		if ok := blockchain.AddBlock(block); !ok {
			fmt.Println("Block rejected: ", block)
			continue
		}

		fmt.Printf("	Mined block %d | txs: %d | hash: %s...\n",
			block.BlockHeight, len(block.Transactions), hashShort(block.Hash))
		for _, tx := range block.Transactions {
			fmt.Printf("		%s -> %s (%.2f)\n", tx.From, tx.To, tx.Amount)
		}
		blockChan <- block

		blockchain.adjustDifficulty()
	}
}

func TransactionsGenerator(memory *MemPool) {
	users := []string{
		"Anuza",
		"Ayazhan",
		"Dana",
		"Ainel",
		"Zhanserik",
		"Bisen",
	}
	for {
		time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)
		tx := NewTransaction(
			users[rand.Intn(len(users))],
			users[rand.Intn(len(users))],
			float64(rand.Intn(100)),
		)
		memory.AddTransaction(*tx)
		fmt.Printf("\nNew Transaction: %s -> %s (%f)\n", hashShort(tx.From), hashShort(tx.To), tx.Amount)
	}
}

func demonstrateMiningProcesses() {
	blockChan := make(chan Block)

	genesisBlock := NewBlock(nil, []Transaction{})

	AAA_blockchain := &Blockchain{CurrentDifficulty: 6, TargetBlockTime: 20}
	memory := &MemPool{TimeToLive: 30}

	AAA_blockchain.AddBlock(*genesisBlock)

	go TransactionsGenerator(memory)
	go proofOfMine(memory, blockChan, AAA_blockchain)

	for block := range blockChan {
		fmt.Printf("📦 Block %d confirmed (txs: %d)\n",
			block.BlockHeight, len(block.Transactions))
	}
}
