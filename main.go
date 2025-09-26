package main

import (
	"fmt"
	"strings"
)

func main() {
	demonstrateBlockchain()
}

func (b *Block) demonstrateMerkleTree() {
	fmt.Printf("\n=== Merkle Tree Structure For Block %s ===\n", hashShort(b.Hash))
	fmt.Printf("Level %d:%s %s\n", len(b.MerkleTree), strings.Repeat(" ", len(b.MerkleTree)*10), hashShort(b.MerkleRoot))

	for i := len(b.MerkleTree) - 1; i >= 0; i-- {
		level := b.MerkleTree[i]
		represent_level := make([]string, len(level))

		for j, hash := range level {
			represent_level[j] = hashShort(hash)
		}

		fmt.Printf("Level %d:%s %s\n", i, strings.Repeat(" ", i*10), strings.Join(represent_level, ", "))
	}
}

func hashShort(hash string) string {
	if len(hash) <= 16 {
		return hash
	}
	return hash[:16] + "...."
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

	blockchain.Blocks[0].demonstrateMerkleTree()
	blockchain.Blocks[1].demonstrateMerkleTree()
}
