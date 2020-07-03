package main

import (
	"errors"
)

type Account struct {
	ActiveCard     bool `json:"activeCard"`
	AvailableLimit int  `json:"availableLimit"`
	transactions   []Transaction
}

var CurrentAccount *Account

const (
	IntervalMinutes          = 2
	MaxFrequencyPerInterval  = 3
	MaxSimilarityPerInterval = 1
)

const (
	AccountAlreadyInitialized  = "account-already-initialized"
	InsufficientLimit          = "insufficient-limit"
	CardNotActive              = "card-not-active"
	HighFrequencySmallInterval = "high-frequency-small-interval"
	DoubledTransaction         = "doubled-transaction"
)

func Initialize(acc Account) []error {
	if CurrentAccount != nil {
		return []error{errors.New(AccountAlreadyInitialized)}
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
	matches := acc.countMatches(tr)
	if matches.frequency == MaxFrequencyPerInterval {
		errs = append(errs, errors.New(HighFrequencySmallInterval))
	}
	if matches.similarity == MaxSimilarityPerInterval {
		errs = append(errs, errors.New(DoubledTransaction))
	}

	if errs == nil {
		*acc = Account{
			ActiveCard:     acc.ActiveCard,
			AvailableLimit: acc.AvailableLimit - tr.Amount,
			transactions:   append(acc.transactions, tr),
		}
	}

	return errs
}

func (acc *Account) countMatches(newTransaction Transaction) matches {
	matches := matches{}
	for _, t := range acc.transactions {
		minutesSinceLastTransaction := newTransaction.Time.Sub(t.Time).Minutes()
		if minutesSinceLastTransaction > IntervalMinutes {
			continue
		}
		if newTransaction.isSimilar(t) {
			matches.similarity++
		}
		matches.frequency++
	}
	return matches
}

type matches struct {
	frequency  int
	similarity int
}
