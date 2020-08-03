package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCountMatches(t *testing.T) {
	tests := map[string]func(*testing.T){
		"Should count and group transaction matches according to interval rules": func(t *testing.T) {
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
