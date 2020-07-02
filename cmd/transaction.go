package main

import (
	"time"
)

type Transaction struct {
	Merchant string    `json:"merchant"`
	Amount   int       `json:"amount"`
	Time     time.Time `json:"time"`
}

func (tr *Transaction) isSimilar(other Transaction) bool {
	return tr.Amount == other.Amount && tr.Merchant == other.Merchant
}
