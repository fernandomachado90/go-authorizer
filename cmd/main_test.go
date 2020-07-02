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
			`{ "account": { "activeCard": true, "availableLimit": 100 } }`,
			`{ "account": { "activeCard": true, "availableLimit": 100 }, "violations": [] }`,
		},
		{
			`{ "transaction": { "merchant": "Alpha", "amount": 20, "time": "2020-07-12T10:00:00.000Z" } }`,
			`{ "account": { "activeCard": true, "availableLimit": 80 }, "violations": [] }`,
		},
		{
			`{ "transaction": { "merchant": "Beta", "amount": 90, "time": "2020-07-12T11:00:00.000Z" } }`,
			`{ "account": { "activeCard": true, "availableLimit": 80 }, "violations": [ "insufficient-limit" ] }`,
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
