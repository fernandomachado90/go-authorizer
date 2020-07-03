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

func Parse(reader io.Reader) *bytes.Buffer {
	var input Payload
	_ = json.NewDecoder(reader).Decode(&input)

	var errs []error
	if input.Account != nil {
		errs = Initialize(*input.Account)
	} else if input.Transaction != nil {
		errs = CurrentAccount.Authorize(*input.Transaction)
	} else {
		return &bytes.Buffer{} // undefined operation
	}

	var output = Payload{
		Account:    CurrentAccount,
		Violations: []string{},
	}
	for _, err := range errs {
		output.Violations = append(output.Violations, err.Error())
	}

	buffer := &bytes.Buffer{}
	_ = json.NewEncoder(buffer).Encode(&output)
	return buffer
}
