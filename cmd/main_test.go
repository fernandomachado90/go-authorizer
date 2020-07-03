package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	// given
	CurrentAccount = nil

	type contract struct {
		input  string
		output string
	}
	tests := []contract{
		{
			`{ "account": { "activeCard": true, "availableLimit": 200 } }`,
			`{ "account": { "activeCard": true, "availableLimit": 200 }, "violations": [] }`,
		},
		{
			`{ "account": { "activeCard": false, "availableLimit": 100 } }`,
			`{ "account": { "activeCard": true, "availableLimit": 200 }, "violations": ["account-already-initialized"] }`,
		},
		{
			`{ "transaction": { "merchant": "Alpha", "amount": 20, "time": "2020-07-12T10:29:59.000Z" } }`,
			`{ "account": { "activeCard": true, "availableLimit": 180 }, "violations": [] }`,
		},
		{
			`{ "transaction": { "merchant": "Alpha", "amount": 40, "time": "2020-07-12T10:30:00.000Z" } }`,
			`{ "account": { "activeCard": true, "availableLimit": 140 }, "violations": [] }`,
		},
		{
			`{ "transaction": { "merchant": "Beta", "amount": 40, "time": "2020-07-12T10:31:00.000Z" } }`,
			`{ "account": { "activeCard": true, "availableLimit": 100 }, "violations": [] }`,
		},
		{
			`{ "transaction": { "merchant": "Omega", "amount": 100, "time": "2020-07-12T10:32:00.000Z" } }`,
			`{ "account": { "activeCard": true, "availableLimit": 0 }, "violations": [] }`,
		},
		{
			`{ "transaction": { "merchant": "Omega", "amount": 100, "time": "2020-07-12T10:32:00.000Z" } }`,
			`{ "account": { "activeCard": true, "availableLimit": 0 }, "violations": ["insufficient-limit", "high-frequency-small-interval", "doubled-transaction"] }`,
		},
	}

	for _, contract := range tests {
		// when
		var stdin bytes.Buffer
		stdin.Write([]byte(contract.input))
		stdout := Parse(&stdin)

		//then
		assert.JSONEq(t, contract.output, stdout.String())
	}
}
