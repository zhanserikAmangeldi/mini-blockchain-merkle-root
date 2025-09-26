package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	BlockHeight  int
	Hash         string
	PrevHash     string
	Timestamp    int64
	Transactions []Transaction
	MerkleRoot   string
	MerkleTree   [][]string
	Nonce        int
}

func (block *Block) IsValid() bool {
	if block.BlockHeight < 0 || len(block.Transactions) == 0 {
		return false
	}

	for _, tx := range block.Transactions {
		if !tx.IsValid() {
			return false
		}

		expectedHash := (&Transaction{
			From:      tx.From,
			To:        tx.To,
			Amount:    tx.Amount,
			Timestamp: tx.Timestamp,
		}).ComputeHash()

		if tx.Hash != expectedHash {
			return false
		}
	}

	expectedMerkleRoot := block.computeMerkleRootInternal()
	if block.MerkleRoot != expectedMerkleRoot {
		return false
	}

	expectedBlockHash := block.computeHashInternal()
	return block.Hash == expectedBlockHash
}

func (block *Block) computeHashInternal() string {
	blockData := strconv.Itoa(block.BlockHeight) +
		block.PrevHash +
		strconv.FormatInt(block.Timestamp, 10) +
		block.MerkleRoot +
		strconv.Itoa(block.Nonce)

	hash := sha256.Sum256([]byte(blockData))

	return fmt.Sprintf("%x", hash)
}

func (block *Block) ComputeHash() string {
	block.Hash = block.computeHashInternal()
	return block.Hash
}

func (block *Block) AddTransaction(tx Transaction) {
	if !tx.IsValid() {
		fmt.Println("Invalid transaction, not adding to block.")
		return
	}

	block.Transactions = append(block.Transactions, tx)

	block.ComputeMerkleRoot()
}

func (block *Block) ComputeMerkleRoot() {
	block.MerkleRoot = block.computeMerkleRootInternal()
}

func (block *Block) computeMerkleRootInternal() string {
	fmt.Printf("\nCOMPUTING MERKLE ROOT FOR BLOCK %d\n", block.BlockHeight)
	fmt.Println(strings.Repeat("=", 40))
	if len(block.Transactions) == 0 {
		return ""
	}
	block.MerkleTree = [][]string{}

	layer := []string{}

	fmt.Printf("Starting with %d transactions:\n", len(block.Transactions))
	for i, tx := range block.Transactions {
		layer = append(layer, tx.Hash)
		fmt.Printf("   TX%d: %s -> %s (%.2f) | Hash: %s\n",
			i+1, tx.From, tx.To, tx.Amount, tx.Hash)
	}

	levelNum := 0

	for len(layer) > 1 {
		var nextLayer []string
		fmt.Println(1)
		block.MerkleTree = append(block.MerkleTree, layer)
		fmt.Println(2)

		levelNum++
		fmt.Printf("\nLevel %d Processing (%d nodes):\n", levelNum, len(layer))
		fmt.Println("Layer:", layer)
		for i := 0; i < len(layer); i += 2 {
			left := layer[i]
			right := left

			if i+1 < len(layer) {
				right = layer[i+1]
				fmt.Printf("   Pair %d: %s + %s\n", (i/2)+1, left, right)
			} else {
				fmt.Printf("   Pair %d: %s + %s (duplicated - odd number)\n", (i/2)+1, left, right)
			}

			combined := left + right
			hash := sha256.Sum256([]byte(combined))
			nextLayer = append(nextLayer, fmt.Sprintf("%x", hash))

			fmt.Printf("      Result: %x\n", hash)
		}

		layer = nextLayer
		fmt.Printf("   → Next level will have %d nodes\n", len(nextLayer))
	}

	fmt.Printf("\nMERKLE ROOT: %s\n", layer[0])
	fmt.Println(strings.Repeat("=", 80))
	return layer[0]
}

func NewBlock(prevBlock *Block, transactions []Transaction) *Block {
	height := 0
	prevHash := strings.Repeat("0", 64)

	if prevBlock != nil {
		height = prevBlock.BlockHeight + 1
		prevHash = prevBlock.Hash
	}

	newBlock := &Block{
		BlockHeight:  height,
		PrevHash:     prevHash,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		Nonce:        0,
		MerkleRoot:   "",
		MerkleTree:   [][]string{},
	}

	newBlock.ComputeMerkleRoot()
	newBlock.ComputeHash()

	return newBlock
}
