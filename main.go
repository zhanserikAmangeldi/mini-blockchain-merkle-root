package main

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

func (blockchain *Blockchain) AddBlock(newBlock Block) {
	blockchain.Blocks = append(blockchain.Blocks, newBlock)
}

func main() {
}
