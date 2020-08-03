package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMemoryDB(t *testing.T) {
	tests := map[string]func(*testing.T){
		"Should create a new memory database": func(t *testing.T) {
			// given

			// when
			db := NewMemoryDB()

			// then
			assert.NotEmpty(t, db)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}
