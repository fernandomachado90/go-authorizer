package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	tests := map[string]func(*testing.T){
		"Should create account when there are no accounts": func(t *testing.T) {
			// given
			CurrentAccount = nil

			// when
			err := Create(Account{
				ActiveCard:     true,
				AvailableLimit: 123,
			})

			// then
			assert.NotEmpty(t, CurrentAccount)
			assert.NoError(t, err)
		},
		"Should not create account when an account is already created": func(t *testing.T) {
			// given
			CurrentAccount = &Account{
				ActiveCard:     true,
				AvailableLimit: 123,
			}

			// when
			err := Create(Account{
				ActiveCard:     false,
				AvailableLimit: 456,
			})

			// then
			assert.NotEmpty(t, CurrentAccount)
			assert.Error(t, err, AccountAlreadyInitialized)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}
