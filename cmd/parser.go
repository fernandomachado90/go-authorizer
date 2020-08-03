package main

import (
	"bytes"
	"encoding/json"
	"io"
)

type Payload struct {
	Account     *Account     `json:"account"`
	Transaction *Transaction `json:"transaction,omitempty"`
	Violations  []string     `json:"violations"`
}

type Parser struct {
	accountManager *AccountManager
	accountActive  Account
}

func (p *Parser) Parse(reader io.Reader) *bytes.Buffer {
	var input Payload
	_ = json.NewDecoder(reader).Decode(&input)

	var errs []error
	var acc Account
	if input.Account != nil {
		acc, errs = p.accountManager.Initialize(*input.Account)
	} else if input.Transaction != nil {
		acc, errs = p.accountManager.Authorize(p.accountActive, *input.Transaction)
	} else {
		return &bytes.Buffer{} // undefined operation
	}

	p.accountActive = acc
	var output = Payload{
		Account:    &acc,
		Violations: []string{},
	}
	for _, err := range errs {
		output.Violations = append(output.Violations, err.Error())
	}

	buffer := &bytes.Buffer{}
	_ = json.NewEncoder(buffer).Encode(&output)
	return buffer
}
