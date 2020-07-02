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

func Authorize(tr Transaction) []error {
	var errs []error

	if CurrentAccount.AvailableLimit-tr.Amount < 0 {
		errs = append(errs, errors.New(InsufficientLimit))
	}
	if !CurrentAccount.ActiveCard {
		errs = append(errs, errors.New(CardNotActive))
	}

	if errs == nil {
		CurrentAccount.AvailableLimit -= tr.Amount
	}
	return errs
}
