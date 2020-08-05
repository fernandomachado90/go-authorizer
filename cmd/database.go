package main

import (
	"errors"
)

type db interface {
	CreateAccount(Account) (Account, error)
	UpdateAccount(Account) Account
	CurrentAccount() (Account, error)
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
		return db.account[0], errors.New("account already exists")
	}
	db.account[0] = acc
	return db.account[0], nil
}

func (db *dbMemory) UpdateAccount(acc Account) Account {
	db.account[0] = acc
	return db.account[0]
}

func (db *dbMemory) CurrentAccount() (Account, error) {
	if len(db.account) == 0 {
		return Account{}, errors.New("no account set")
	}
	return db.account[0], nil
}
