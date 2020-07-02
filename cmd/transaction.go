package main

import (
	"errors"
	"time"
)

type Transaction struct {
	Merchant string    `json:"merchant"`
	Amount   int       `json:"amount"`
	Time     time.Time `json:"time"`
}

func (acc *Account) Authorize(tr Transaction) []error {
	var errs []error

	if acc.AvailableLimit-tr.Amount < 0 {
		errs = append(errs, errors.New(InsufficientLimit))
	}
	if !acc.ActiveCard {
		errs = append(errs, errors.New(CardNotActive))
	}

	if errs == nil {
		acc.AvailableLimit -= tr.Amount
	}
	return errs
}
