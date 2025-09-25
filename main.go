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
	demonstrateBlockchain()
}

func demonstrateBlockchain() {
	fmt.Println(strings.Repeat("=", 20) + " Blockchain Demonstration " + strings.Repeat("=", 20))

	blockchain := &Blockchain{}
	mempool := &MemPool{TimeToLive: 300}

	fmt.Println("1. Creating transactions...")
	tx1 := NewTransaction("Anuza", "Ayazhan", 100.50)
	tx2 := NewTransaction("Dana", "Ainel", 50.25)
	tx3 := NewTransaction("Zhanserik", "Bisen", 30.75)

	fmt.Printf("   TX1: %s -> %s (%.2f) Hash: %s\n", tx1.From, tx1.To, tx1.Amount, tx1.Hash[:16]+"...")
	fmt.Printf("   TX2: %s -> %s (%.2f) Hash: %s\n", tx2.From, tx2.To, tx2.Amount, tx2.Hash[:16]+"...")
	fmt.Printf("   TX3: %s -> %s (%.2f) Hash: %s\n", tx3.From, tx3.To, tx3.Amount, tx3.Hash[:16]+"...")

	mempool.AddTransaction(*tx1)
	mempool.AddTransaction(*tx2)
	mempool.AddTransaction(*tx3)

	fmt.Println("\n2. Creating genesis block...")
	genesisBlock := NewBlock(nil, mempool.GetTransactions(3))
	fmt.Printf("   Genesis Block Height: %d\n", genesisBlock.BlockHeight)
	fmt.Printf("   Genesis Block Hash: %s\n", genesisBlock.Hash[:32]+"...")
	fmt.Printf("   Merkle Root: %s\n", genesisBlock.MerkleRoot[:32]+"...")
	fmt.Printf("   Transactions: %d\n", len(genesisBlock.Transactions))

	fmt.Println("\n3. Validating and adding genesis block...")
	if genesisBlock.IsValid() {
		fmt.Println("   ✓ Genesis block is valid")
		if blockchain.AddBlock(*genesisBlock) {
			fmt.Println("   ✓ Genesis block added to blockchain")
		}
	} else {
		fmt.Println("   ✗ Genesis block is invalid")
	}

	fmt.Println("\n4. Creating second block...")
	tx4 := NewTransaction("Bisen", "Zhanserik", 70.80)
	tx5 := NewTransaction("Ayazhan", "Dana", 25.30)

	mempool.AddTransaction(*tx4)
	mempool.AddTransaction(*tx5)

	secondBlock := NewBlock(blockchain.GetLastBlock(), mempool.GetTransactions(2))
	fmt.Printf("   Second Block Height: %d\n", secondBlock.BlockHeight)
	fmt.Printf("   Second Block Hash: %s\n", secondBlock.Hash[:32]+"...")
	fmt.Printf("   Previous Hash: %s\n", secondBlock.PrevHash[:32]+"...")
	fmt.Printf("   Merkle Root: %s\n", secondBlock.MerkleRoot[:32]+"...")

	fmt.Println("\n5. Validating and adding second block...")
	if secondBlock.IsValid() {
		fmt.Println("   ✓ Second block is valid")
		if blockchain.AddBlock(*secondBlock) {
			fmt.Println("   ✓ Second block added to blockchain")
		}
	}

	fmt.Println("\n6. Validating entire blockchain...")
	if blockchain.isValid() {
		fmt.Println("   ✓ Blockchain is valid")
	} else {
		fmt.Println("   ✗ Blockchain is invalid")
	}

	fmt.Printf("\n=== Final Blockchain State ===\n")
	fmt.Printf("Total blocks: %d\n", len(blockchain.Blocks))
	for i, block := range blockchain.Blocks {
		fmt.Printf("Block %d: Hash=%s, Transactions=%d\n",
			i, block.Hash[:16]+"...", len(block.Transactions))
	}
}
