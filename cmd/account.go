package main

import "errors"

type Account struct {
	ActiveCard     bool `json:"activeCard"`
	AvailableLimit int  `json:"availableLimit"`
}

func Create(a Account) error {
	if CurrentAccount != nil {
		return errors.New(AccountAlreadyInitialized)
	}
	CurrentAccount = &a
	return nil
}
