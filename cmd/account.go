package main

import "errors"

type Account struct {
	ActiveCard     bool `json:"activeCard"`
	AvailableLimit int  `json:"availableLimit"`
}

func Create(acc Account) error {
	if CurrentAccount != nil {
		return errors.New(AccountAlreadyInitialized)
	}
	CurrentAccount = &acc
	return nil
}
