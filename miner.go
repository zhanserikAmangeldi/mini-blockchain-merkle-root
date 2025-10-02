package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
)

func mine(num int, previousHash string) string {
	counter := 0
	lead_zeros := num

	for {
		word := previousHash + strconv.Itoa(counter)
		hashed_word := sha256.Sum256([]byte(word))

		hex_word := fmt.Sprintf("%x", hashed_word)
		zero_or_not, err := strconv.Atoi(hex_word[:lead_zeros])
		if err != nil {
			zero_or_not = 1
		}

		if zero_or_not == 0 && hashed_word[lead_zeros] != 0 {
			return hex_word
		}

		counter++
	}
}
