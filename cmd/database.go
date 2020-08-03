package main

import (
	"errors"
)

type db interface {
	CreateAccount(Account) (Account, error)
	UpdateAccount(Account) (Account, error)
}

type dbMemory struct {
	account map[int]Account
}

func NewMemoryDB() *dbMemory {
	return &dbMemory{
		account: map[int]Account{},
	}
}

func (db *dbMemory) CreateAccount(acc Account) (Account, error) {
	if len(db.account) > 0 {
		return db.account[0], errors.New("an account already exists")
	}
	db.account[0] = acc
	return acc, nil
}

func (db *dbMemory) UpdateAccount(acc Account) (Account, error) {
	db.account[0] = acc
	return acc, nil
}
