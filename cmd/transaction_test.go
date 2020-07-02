package main

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuthorizeTransaction(t *testing.T) {
	tests := map[string]func(*testing.T){
		"Should authorize transaction with no violations": func(t *testing.T) {
			// given
			CurrentAccount = &Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			}

			// when
			errs := Authorize(Transaction{
				Merchant: "Acme Corporation",
				Amount:   20,
				Time:     time.Now(),
			})

			// then
			assert.Equal(t, 80, CurrentAccount.AvailableLimit)
			assert.Empty(t, errs)
		},
		"Should not authorize transaction due to insufficient limit violation": func(t *testing.T) {
			// given
			CurrentAccount = &Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			}

			// when
			errs := Authorize(Transaction{
				Merchant: "Acme Corporation",
				Amount:   200,
				Time:     time.Now(),
			})

			// then
			assert.Equal(t, 100, CurrentAccount.AvailableLimit)
			assert.Contains(t, errs, errors.New(InsufficientLimit))
		},
		"Should not authorize transaction due to card not active violation": func(t *testing.T) {
			// given
			CurrentAccount = &Account{
				ActiveCard:     false,
				AvailableLimit: 100,
			}

			// when
			errs := Authorize(Transaction{
				Merchant: "Acme Corporation",
				Amount:   10,
				Time:     time.Now(),
			})

			// then
			assert.Equal(t, 100, CurrentAccount.AvailableLimit)
			assert.Contains(t, errs, errors.New(CardNotActive))
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}
