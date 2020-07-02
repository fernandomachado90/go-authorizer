package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	tests := map[string]func(*testing.T){
		"Should initialize entity": func(t *testing.T) {
			//given

			// when
			account := Account{}

			// then
			assert.Empty(t, account)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}
