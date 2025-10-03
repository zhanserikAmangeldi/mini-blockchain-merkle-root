package main

import "fmt"

type Blockchain struct {
	Blocks            []Block
	CurrentDifficulty int
	TargetBlockTime   int
}

func (blockchain *Blockchain) isValid() bool {
	if len(blockchain.Blocks) == 0 {
		return true
	}

	if !blockchain.Blocks[0].IsValid() {
		return false
	}

	for i := 1; i < len(blockchain.Blocks); i++ {
		currentBlock := blockchain.Blocks[i]
		prevBlock := blockchain.Blocks[i-1]

		if !currentBlock.IsValid() {
			return false
		}

		if currentBlock.PrevHash != prevBlock.Hash {
			return false
		}

		if currentBlock.BlockHeight != prevBlock.BlockHeight+1 {
			return false
		}
	}

	return true
}

func (blockchain *Blockchain) AddBlock(newBlock Block) bool {
	if !newBlock.IsValid() {
		fmt.Println("Invalid block, cannot add to blockchain")
		return false
	}

	if len(blockchain.Blocks) > 0 {
		lastBlock := blockchain.Blocks[len(blockchain.Blocks)-1]
		if newBlock.PrevHash != lastBlock.Hash || newBlock.BlockHeight != lastBlock.BlockHeight+1 {
			fmt.Println("Block does not link properly to the last block, cannot add to blockchain")
			return false
		}
	}

	blockchain.Blocks = append(blockchain.Blocks, newBlock)
	return true
}

func (blockchain *Blockchain) GetLastBlock() *Block {
	if len(blockchain.Blocks) == 0 {
		return nil
	}
	return &blockchain.Blocks[len(blockchain.Blocks)-1]
}

func (blockchain *Blockchain) adjustDifficulty() {
	if len(blockchain.Blocks) < 2 {
		return
	}

	lastBlock := blockchain.Blocks[len(blockchain.Blocks)-1]
	prevBlock := blockchain.Blocks[len(blockchain.Blocks)-2]

	actualTime := lastBlock.Timestamp - prevBlock.Timestamp

	if actualTime < int64(blockchain.TargetBlockTime/2) {
		blockchain.CurrentDifficulty++
		fmt.Printf("\nToo fast, increasing difficulty: from %d to %d\n", blockchain.CurrentDifficulty-1, blockchain.CurrentDifficulty)
	} else if actualTime > int64(blockchain.TargetBlockTime*2) {
		if blockchain.CurrentDifficulty > 1 {
			blockchain.CurrentDifficulty--
			fmt.Printf("\nToo slow, decreasing difficulty: from %d to %d\n", blockchain.CurrentDifficulty+1, blockchain.CurrentDifficulty)
		}
	}

	fmt.Println(blockchain.CurrentDifficulty)
}
