package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSimilarTransaction(t *testing.T) {
	tests := map[string]func(*testing.T){
		"Should detect a similar transaction": func(t *testing.T) {
			// given
			tr := &Transaction{
				Merchant: "Hello World",
				Amount:   100,
			}

			// when
			similar := tr.isSimilar(Transaction{
				Merchant: "Hello World",
				Amount:   100,
			})

			// then
			assert.True(t, similar)
		},
		"Should not detect a similar transaction due to different merchant": func(t *testing.T) {
			// given
			tr := &Transaction{
				Merchant: "Hello",
				Amount:   100,
			}

			// when
			similar := tr.isSimilar(Transaction{
				Merchant: "Hello World",
				Amount:   100,
			})

			// then
			assert.False(t, similar)
		},
		"Should not detect a similar transaction due to different amount": func(t *testing.T) {
			// given
			tr := &Transaction{
				Merchant: "Hello World",
				Amount:   100,
			}

			// when
			similar := tr.isSimilar(Transaction{
				Merchant: "Hello World",
				Amount:   50,
			})

			// then
			assert.False(t, similar)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}
