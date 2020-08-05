package main

import (
	"bytes"
	"encoding/json"
	"io"
)

type Handler struct {
	db             DB
	accountHandler AccountHandler
}

type AccountHandler interface {
	Initialize(Account) (Account, []error)
	Authorize(Account, Transaction) (Account, []error)
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

func (h *Handler) Dispatch(request interface{}) (Account, []error) {
	switch req := request.(type) {
	case Account:
		return h.accountHandler.Initialize(req)
	case Transaction:
		acc, _ := h.db.CurrentAccount()
		return h.accountHandler.Authorize(acc, req)
	default:
		acc, _ := h.db.CurrentAccount()
		return acc, nil
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
