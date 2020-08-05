package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMemoryDB(t *testing.T) {
	tests := map[string]func(*testing.T){
		"Should create a new in-memory database": func(t *testing.T) {
			// when
			db := NewMemoryDB()

			// then
			assert.NotEmpty(t, db)
			assert.Empty(t, db.account)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}

func TestCreateAccount(t *testing.T) {
	tests := map[string]func(*testing.T){
		"Should create an account": func(t *testing.T) {
			// given
			db := NewMemoryDB()
			acc := Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			}

			// when
			res, err := db.CreateAccount(acc)

			// then
			assert.Equal(t, acc, res)
			assert.NoError(t, err)
		},
		"Should not create an account because one already exists": func(t *testing.T) {
			// given
			db := NewMemoryDB()
			existing := Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			}
			db.CreateAccount(existing)

			acc := Account{
				ActiveCard:     true,
				AvailableLimit: 200,
			}

			// when
			res, err := db.CreateAccount(acc)

			// then
			assert.Equal(t, existing, res)
			assert.Error(t, err, "account already exists")
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}

func TestUpdateAccount(t *testing.T) {
	tests := map[string]func(*testing.T){
		"Should update an account": func(t *testing.T) {
			// given
			db := NewMemoryDB()
			existing := Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			}
			db.CreateAccount(existing)

			acc := Account{
				ActiveCard:     false,
				AvailableLimit: 200,
			}

			// when
			res := db.UpdateAccount(acc)

			// then
			assert.Equal(t, acc, res)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}

func TestCurrentAccount(t *testing.T) {
	tests := map[string]func(*testing.T){
		"Should get current account": func(t *testing.T) {
			// given
			db := NewMemoryDB()
			existing := Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			}
			db.CreateAccount(existing)

			// when
			res, err := db.CurrentAccount()

			// then
			assert.Equal(t, existing, res)
			assert.NoError(t, err)
		},
		"Should not get current account because it does not exists": func(t *testing.T) {
			// given
			db := NewMemoryDB()

			// when
			res, err := db.CurrentAccount()

			// then
			assert.Empty(t, res)
			assert.Error(t, err, "no account set")
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}
