package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	Hash      string
	Amount    float64
	From      string
	To        string
	Timestamp int64
}

func (tx *Transaction) IsValid() bool {
	return tx.Amount > 0 && tx.From != "" && tx.To != ""
}

func (tx *Transaction) ComputeHash() string {
	transactionData := tx.From + tx.To + strconv.FormatFloat(tx.Amount, 'f', 6, 64) + strconv.FormatInt(tx.Timestamp, 10)
	hash := sha256.Sum256([]byte(transactionData))
	return fmt.Sprintf("%x", hash)
}

func NewTransaction(from, to string, amount float64) *Transaction {
	tx := &Transaction{
		From:      from,
		To:        to,
		Amount:    amount,
		Timestamp: time.Now().Unix(),
	}
	tx.Hash = tx.ComputeHash()
	return tx
}

type Block struct {
	BlockHeight  int
	Hash         string
	PrevHash     string
	Timestamp    int64
	Transactions []Transaction
	MerkleRoot   string
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
	if len(block.Transactions) == 0 {
		return ""
	}

	layer := []string{}
	for _, tx := range block.Transactions {
		layer = append(layer, tx.From)
	}

	for len(layer) > 1 {
		var nextLayer []string

		for i := 0; i < len(layer)-1; i += 2 {
			left := layer[i]
			right := left

			if i+1 < len(layer) {
				right = layer[i+1]
			}

			combined := left + right
			hash := sha256.Sum256([]byte(combined))
			nextLayer = append(nextLayer, fmt.Sprintf("%x", hash))
		}

		layer = nextLayer
		fmt.Println("Merkle layer:", layer)
	}

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
	}

	newBlock.ComputeMerkleRoot()
	newBlock.ComputeHash()

	return newBlock
}

type Blockchain struct {
	Blocks []Block
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

func main() {
	blockchain := &Blockchain{}
	mempool := &MemPool{TimeToLive: 300}

	tx1 := NewTransaction("Aral", "Zhanserik", 100)
	mempool.AddTransaction(*tx1)
	tx2 := NewTransaction("Danyal", "Ilyas", 50)
	mempool.AddTransaction(*tx2)
	tx3 := NewTransaction("Zhanserik", "Danyal", 30)
	mempool.AddTransaction(*tx3)
	// fmt.Println("Mempool before creating block:", mempool.Transactions)

	genesisBlock := NewBlock(nil, mempool.GetTransactions(3))

	// fmt.Println(genesisBlock)
	blockchain.AddBlock(*genesisBlock)

	// fmt.Println("Blockchain:", blockchain)
	tx4 := NewTransaction("Aral", "Ilyas", 70)
	mempool.AddTransaction(*tx4)

	blockSecond := NewBlock(genesisBlock, mempool.GetTransactions(3))
	blockchain.AddBlock(*blockSecond)
	// fmt.Println(blockchain)
	blockchain.Blocks[0].ComputeMerkleRoot()
	fmt.Println(blockchain.Blocks[0])
	fmt.Println(blockchain.Blocks[0].MerkleRoot)
}
