package main

import (
	"errors"
	"fmt"
)

var CurrentAccount *Account

type Account struct {
	ActiveCard     bool `json:"activeCard"`
	AvailableLimit int  `json:"availableLimit"`
	transactions   []Transaction
}

type bufferMatches struct {
	frequency int
	similar   int
}

const BufferIntervalMinutes = 2
const MaxFrequencyPerInterval = 3
const MaxSimilarityPerInterval = 1

func Initialize(acc Account) error {
	if CurrentAccount != nil {
		return errors.New(AccountAlreadyInitialized)
	}
	CurrentAccount = &acc
	return nil
}

func (acc *Account) Authorize(tr Transaction) []error {
	var errs []error

	if acc.AvailableLimit-tr.Amount < 0 {
		errs = append(errs, errors.New(InsufficientLimit))
	}
	if !acc.ActiveCard {
		errs = append(errs, errors.New(CardNotActive))
	}
	matches := acc.countBufferMatches(tr)
	if matches.frequency == MaxFrequencyPerInterval {
		errs = append(errs, errors.New(HighFrequencySmallInterval))
	}
	if matches.similar == MaxSimilarityPerInterval {
		errs = append(errs, errors.New(DoubledTransaction))
	}

	if errs == nil {
		shouldClearBuffer := matches.frequency == 0
		acc.commit(tr, shouldClearBuffer)
	}
	return errs
}

func (acc *Account) commit(tr Transaction, clearBuffer bool) {
	acc.AvailableLimit -= tr.Amount
	if clearBuffer {
		acc.transactions = []Transaction{}
	}
	acc.transactions = append(acc.transactions, tr)
}

func (acc *Account) countBufferMatches(newTransaction Transaction) bufferMatches {
	matches := bufferMatches{}
	for _, t := range acc.transactions {
		timeSinceLastTransaction := newTransaction.Time.Sub(t.Time)
		fmt.Println(newTransaction.Time)
		fmt.Println(t.Time)
		fmt.Println(timeSinceLastTransaction.Minutes())
		if timeSinceLastTransaction.Minutes() > BufferIntervalMinutes {
			continue
		}
		if newTransaction.isSimilar(t) {
			matches.similar++
		}
		matches.frequency++
	}
	return matches
}
