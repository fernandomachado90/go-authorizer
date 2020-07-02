package main

type Account struct {
	ActiveCard     bool
	AvailableLimit int
}

func Create(account Account) (Account, error) {
	return account, nil
}
