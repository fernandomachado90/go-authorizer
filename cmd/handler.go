package main

import (
	"bytes"
	"encoding/json"
	"io"
)

type Handler struct {
	accountManager *AccountManager
}

func (h *Handler) Decode(reader io.Reader) interface{} {
	type payload struct {
		Account     *Account     `json:"account"`
		Transaction *Transaction `json:"transaction"`
	}

	var input payload
	_ = json.NewDecoder(reader).Decode(&input)

	if input.Account != nil {
		return *input.Account
	}
	if input.Transaction != nil {
		return *input.Transaction
	}
	return nil
}

func (h *Handler) Process(request interface{}) (Account, []error) {
	switch r := request.(type) {
	case Account:
		return h.accountManager.Initialize(r)
	case Transaction:
		acc, _ := h.accountManager.db.CurrentAccount()
		return h.accountManager.Authorize(acc, r)
	default:
		return Account{}, nil
	}
}

func (h *Handler) Encode(acc Account, errs []error) *bytes.Buffer {
	type payload struct {
		Account    *Account `json:"account"`
		Violations []string `json:"violations"`
	}

	var output = payload{
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
