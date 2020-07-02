package main

import (
	"bytes"
	"encoding/json"
	"io"
)

type Payload struct {
	Account    *Account `json:"account"`
	Violations []string `json:"violations"`
}

func Parse(reader io.Reader) *bytes.Buffer {
	var err error

	var payload Payload
	_ = json.NewDecoder(reader).Decode(&payload)

	if payload.Account != nil {
		err = Initialize(*payload.Account)
		payload.Account = CurrentAccount // todo revise this attribuition
	}

	payload.Violations = []string{}
	if err != nil {
		payload.Violations = append(payload.Violations, err.Error())
	}

	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	_ = encoder.Encode(payload)

	return buffer
}
