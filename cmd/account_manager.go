package main

import (
	"errors"
)

type AccountManager struct {
	db DB
}

func NewAccountManager(db DB) *AccountManager {
	return &AccountManager{db}
}

func (m *AccountManager) Initialize(acc Account) (Account, []error) {
	var errs []error

	acc, err := m.db.CreateAccount(acc)
	if err != nil {
		errs = append(errs, errors.New(AccountAlreadyInitialized))
	}

	return acc, errs
}

func (m *AccountManager) Authorize(acc Account, tr Transaction) (Account, []error) {
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
		acc = m.db.UpdateAccount(Account{
			ActiveCard:     acc.ActiveCard,
			AvailableLimit: acc.AvailableLimit - tr.Amount,
			transactions:   append(acc.transactions, tr),
		})
	}

	return acc, errs
}

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
