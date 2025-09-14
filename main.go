package main

import (
	"fmt"
	"strconv"
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
	hash := tx.From + tx.To + strconv.FormatFloat(tx.Amount, 'f', 6, 64) + strconv.FormatInt(tx.Timestamp, 10)
	return hash
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
	MerleRoot    string
}

func (block *Block) IsValid() bool {
	return block.BlockHeight >= 0 && block.Hash != "" && block.PrevHash != "" && len(block.Transactions) > 0
}

func (block *Block) ComputeHash() string {
	hash := strconv.Itoa(int(block.Timestamp))
	return hash
}

func (block *Block) AddTransaction(tx Transaction) {
	if !tx.IsValid() {
		return
	}

	block.Transactions = append(block.Transactions, tx)
}

func NewBlock(prevBlock *Block, transactions []Transaction) *Block {
	height := 0
	prevHash := "0000000000000000000000000000000000000000000000000000000000000000000"
	if prevBlock != nil {
		height = prevBlock.BlockHeight + 1
		prevHash = prevBlock.Hash
	}

	newBlock := &Block{
		BlockHeight:  height,
		PrevHash:     prevHash,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
	}
	newBlock.Hash = newBlock.ComputeHash()
	return newBlock
}

type Blockchain struct {
	Blocks []Block
}

type MemPool struct {
	Transactions []Transaction
	TimeToLive   int64
}

func (mempool *MemPool) AddTransaction(tx Transaction) {
	if !tx.IsValid() {
		return
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
	txs := mempool.Transactions[:limit]
	mempool.Transactions = mempool.Transactions[limit:]

	return txs
}

func (blockchain *Blockchain) AddBlock(newBlock Block) {
	if !newBlock.IsValid() {
		return
	}
	blockchain.Blocks = append(blockchain.Blocks, newBlock)
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
	fmt.Println("Mempool before creating block:", mempool.Transactions)

	genesisBlock := NewBlock(nil, mempool.GetTransactions(3))

	fmt.Println(genesisBlock)
	blockchain.AddBlock(*genesisBlock)

	fmt.Println("Blockchain:", blockchain)
	tx4 := NewTransaction("Aral", "Ilyas", 70)
	mempool.AddTransaction(*tx4)

	blockSecond := NewBlock(genesisBlock, mempool.GetTransactions(3))
	blockchain.AddBlock(*blockSecond)
	fmt.Println(blockchain)
}
