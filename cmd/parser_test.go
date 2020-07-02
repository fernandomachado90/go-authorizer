package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := map[string]func(*testing.T){
		"Should create account when there are no account": func(t *testing.T) {
			// given
			CurrentAccount = nil

			var stdin bytes.Buffer
			stdin.Write([]byte(`{ "account": { "activeCard": true, "availableLimit": 100 } }`))

			// when
			stdout := Parse(&stdin)

			// then
			assert.JSONEq(t, `{"account":{"activeCard":true,"availableLimit":100},"violations":[]}`, stdout.String())
		},
		"Should not create account when an account is already created": func(t *testing.T) {
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
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}
