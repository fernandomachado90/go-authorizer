package main

import "github.com/stretchr/testify/mock"

type dbMock struct {
	mock.Mock
}

func NewDatabaseMock() *dbMock {
	return &dbMock{}
}

func (db *dbMock) CreateAccount(acc Account) (Account, error) {
	args := db.Called(acc)
	if args == nil {
		return acc, nil
	}
	res := args.Get(0).(Account)
	err := args.Error(1)
	return res, err
}

func (db *dbMock) UpdateAccount(acc Account) Account {
	_ = db.Called(acc)
	return acc
}

func (db *dbMock) CurrentAccount() (Account, error) {
	args := db.Called()
	res := args.Get(0).(Account)
	err := args.Error(1)
	return res, err
}
