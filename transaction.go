package main

import (
	"crypto/sha256"
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

func (tx *Transaction) IsValid() bool { // TODO: Стоит улучшить, добавить че та или декомпозировать проверки
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
