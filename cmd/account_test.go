package main

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitializeAccount(t *testing.T) {
	tests := map[string]func(*testing.T){
		"Should initialize account when there are no accounts": func(t *testing.T) {
			// given
			CurrentAccount = nil

			// when
			errs := Initialize(Account{
				ActiveCard:     true,
				AvailableLimit: 123,
			})

			// then
			assert.NotEmpty(t, CurrentAccount)
			assert.Empty(t, errs)
		},
		"Should not initialize account when an account is already initialized": func(t *testing.T) {
			// given
			CurrentAccount = &Account{
				ActiveCard:     true,
				AvailableLimit: 123,
			}

			// when
			errs := Initialize(Account{
				ActiveCard:     false,
				AvailableLimit: 456,
			})

			// then
			assert.NotEmpty(t, CurrentAccount)
			assert.Len(t, errs, 1)
			assert.Contains(t, errs, errors.New(AccountAlreadyInitialized))
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}

func TestAuthorizeTransaction(t *testing.T) {
	tests := map[string]func(*testing.T){
		"Should authorize transaction with no violations": func(t *testing.T) {
			// given
			account := &Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			}

			// when
			errs := account.Authorize(Transaction{
				Merchant: "Acme Corporation",
				Amount:   20,
				Time:     time.Now(),
			})

			// then
			assert.Equal(t, 80, account.AvailableLimit)
			assert.Len(t, account.transactions, 1)
			assert.Empty(t, errs)
		},
		"Should not authorize transaction due to insufficient limit violation": func(t *testing.T) {
			// given
			account := &Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			}

			// when
			errs := account.Authorize(Transaction{
				Merchant: "Acme Corporation",
				Amount:   200,
				Time:     time.Now(),
			})

			// then
			assert.Equal(t, 100, account.AvailableLimit)
			assert.Len(t, account.transactions, 0)
			assert.Len(t, errs, 1)
			assert.Contains(t, errs, errors.New(InsufficientLimit))
		},
		"Should not authorize transaction due to card not active violation": func(t *testing.T) {
			// given
			account := &Account{
				ActiveCard:     false,
				AvailableLimit: 100,
			}

			// when
			errs := account.Authorize(Transaction{
				Merchant: "Acme Corporation",
				Amount:   20,
				Time:     time.Now(),
			})

			// then
			assert.Equal(t, 100, account.AvailableLimit)
			assert.Len(t, account.transactions, 0)
			assert.Len(t, errs, 1)
			assert.Contains(t, errs, errors.New(CardNotActive))
		},
		"Should not authorize transaction due to high frequency on small interval violation": func(t *testing.T) {
			// given
			account := &Account{
				ActiveCard:     true,
				AvailableLimit: 100,
				transactions: []Transaction{
					{Time: time.Date(2020, 7, 12, 10, 30, 0, 0, time.UTC)},
					{Time: time.Date(2020, 7, 12, 10, 31, 0, 0, time.UTC)},
					{Time: time.Date(2020, 7, 12, 10, 31, 30, 0, time.UTC)},
				},
			}

			// when
			errs := account.Authorize(Transaction{
				Merchant: "Acme Corporation",
				Amount:   20,
				Time:     time.Date(2020, 7, 12, 10, 32, 0, 0, time.UTC),
			})

			// then
			assert.Equal(t, 100, account.AvailableLimit)
			assert.Len(t, account.transactions, 3)
			assert.Len(t, errs, 1)
			assert.Contains(t, errs, errors.New(HighFrequencySmallInterval))
		},
		"Should not authorize transaction due to doubled transaction violation": func(t *testing.T) {
			// given
			account := &Account{
				ActiveCard:     true,
				AvailableLimit: 100,
				transactions: []Transaction{
					{
						Merchant: "Acme Corporation",
						Amount:   20,
						Time:     time.Date(2020, 7, 12, 10, 30, 0, 0, time.UTC),
					},
				},
			}

			// when
			errs := account.Authorize(Transaction{
				Merchant: "Acme Corporation",
				Amount:   20,
				Time:     time.Date(2020, 7, 12, 10, 31, 0, 0, time.UTC),
			})

			// then
			assert.Equal(t, 100, account.AvailableLimit)
			assert.Len(t, account.transactions, 1)
			assert.Len(t, errs, 1)
			assert.Contains(t, errs, errors.New(DoubledTransaction))
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}

func TestCountMatches(t *testing.T) {
	tests := map[string]func(*testing.T){
		"Should count matches according to the defined interval": func(t *testing.T) {
			// given
			account := &Account{
				transactions: []Transaction{
					{
						Merchant: "Alpha",
						Amount:   10,
						Time:     time.Date(2020, 7, 12, 10, 30, 0, 0, time.UTC),
					},
					{
						Merchant: "Beta",
						Amount:   20,
						Time:     time.Date(2020, 7, 12, 10, 31, 0, 0, time.UTC),
					},
					{
						Merchant: "Gamma",
						Amount:   30,
						Time:     time.Date(2020, 7, 12, 10, 32, 0, 0, time.UTC),
					},
				},
			}

			// when
			matches := account.countMatches(Transaction{
				Merchant: "Gamma",
				Amount:   30,
				Time:     time.Date(2020, 7, 12, 10, 32, 1, 0, time.UTC),
			})

			// then
			assert.Equal(t, 2, matches.frequency)
			assert.Equal(t, 1, matches.similarity)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}
