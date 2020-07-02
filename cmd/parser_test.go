package main

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := map[string]func(*testing.T){
		"Should initialize account with no violations": func(t *testing.T) {
			// given
			CurrentAccount = nil

			var stdin bytes.Buffer
			stdin.Write([]byte(`{ "account": { "activeCard": true, "availableLimit": 100 } }`))

			// when
			stdout := Parse(&stdin)

			// then
			assert.JSONEq(t, `{"account":{"activeCard":true,"availableLimit":100},"violations":[]}`, stdout.String())
		},
		"Should not initialize account due to account already initialized violation": func(t *testing.T) {
			// given
			CurrentAccount = &Account{
				ActiveCard:     true,
				AvailableLimit: 500,
			}

			var stdin bytes.Buffer
			stdin.Write([]byte(`{ "account": { "activeCard": false, "availableLimit": 100 } }`))

			// when
			stdout := Parse(&stdin)

			// then
			expected := fmt.Sprintf(`{"account":{"activeCard":true,"availableLimit":500},"violations":["%s"]}`, AccountAlreadyInitialized)
			assert.JSONEq(t, expected, stdout.String())
		},
		"Should authorize transaction with no violations": func(t *testing.T) {
			// given
			CurrentAccount = &Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			}

			var stdin bytes.Buffer
			stdin.Write([]byte(` { "transaction": { "merchant": "Acme Corporation", "amount": 20, "time": "2020-07-12T10:00:00.000Z" } }`))

			// when
			stdout := Parse(&stdin)

			// then
			assert.JSONEq(t, `{"account":{"activeCard":true,"availableLimit":80},"violations":[]}`, stdout.String())
		},
		"Should not authorize transaction due to insufficient limit violation": func(t *testing.T) {
			// given
			CurrentAccount = &Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			}

			var stdin bytes.Buffer
			stdin.Write([]byte(` { "transaction": { "merchant": "Acme Corporation", "amount": 999, "time": "2020-07-12T10:00:00.000Z" } }`))

			// when
			stdout := Parse(&stdin)

			// then
			expected := fmt.Sprintf(`{"account":{"activeCard":true,"availableLimit":100},"violations":["%s"]}`, InsufficientLimit)
			assert.JSONEq(t, expected, stdout.String())
		},
		"Should not authorize transaction due to card not active violation": func(t *testing.T) {
			// given
			CurrentAccount = &Account{
				ActiveCard:     false,
				AvailableLimit: 100,
			}

			var stdin bytes.Buffer
			stdin.Write([]byte(` { "transaction": { "merchant": "Acme Corporation", "amount": 20, "time": "2020-07-12T10:00:00.000Z" } }`))

			// when
			stdout := Parse(&stdin)

			// then
			expected := fmt.Sprintf(`{"account":{"activeCard":false,"availableLimit":100},"violations":["%s"]}`, CardNotActive)
			assert.JSONEq(t, expected, stdout.String())
		},
		"Should not authorize transaction due to high frequency on a small interval violation": func(t *testing.T) {
			// given
			CurrentAccount = &Account{
				ActiveCard:     true,
				AvailableLimit: 100,
				transactions: []Transaction{
					{Time: time.Date(2020, 7, 12, 10, 30, 0, 0, time.UTC)},
					{Time: time.Date(2020, 7, 12, 10, 31, 0, 0, time.UTC)},
					{Time: time.Date(2020, 7, 12, 10, 31, 30, 0, time.UTC)},
				},
			}

			var stdin bytes.Buffer
			stdin.Write([]byte(` { "transaction": { "merchant": "Acme Corporation", "amount": 50, "time": "2020-7-12T10:30:00.000Z" } }`))

			// when
			stdout := Parse(&stdin)

			// then
			expected := fmt.Sprintf(`{"account":{"activeCard":true,"availableLimit":100},"violations":["%s"]}`, HighFrequencySmallInterval)
			assert.JSONEq(t, expected, stdout.String())
		},
		"Should not authorize transaction due to doubled transaction violation": func(t *testing.T) {
			// given
			CurrentAccount = &Account{
				ActiveCard:     true,
				AvailableLimit: 100,
				transactions: []Transaction{
					{
						Merchant: "Acme Corporation",
						Amount:   50,
					},
				},
			}

			var stdin bytes.Buffer
			stdin.Write([]byte(` { "transaction": { "merchant": "Acme Corporation", "amount": 50, "time": "2020-7-12T10:30:00.000Z" } }`))

			// when
			stdout := Parse(&stdin)

			// then
			expected := fmt.Sprintf(`{"account":{"activeCard":true,"availableLimit":100},"violations":["%s"]}`, DoubledTransaction)
			assert.JSONEq(t, expected, stdout.String())
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}
