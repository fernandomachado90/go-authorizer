package main

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDecode(t *testing.T) {
	tests := map[string]func(*testing.T){
		"Should decode account": func(t *testing.T) {
			// given
			h := Handler{}
			var stdin bytes.Buffer
			stdin.Write([]byte(`{ "account": { "activeCard": true, "availableLimit": 100 } }`))

			// when
			res := h.Decode(&stdin)

			// then
			acc := res.(Account)
			assert.Equal(t, true, acc.ActiveCard)
			assert.Equal(t, 100, acc.AvailableLimit)
		},
		"Should decode transaction": func(t *testing.T) {
			// given
			h := Handler{}
			var stdin bytes.Buffer
			stdin.Write([]byte(`{ "transaction": { "merchant": "Acme Corporation", "amount": 20, "time": "2020-07-12T10:00:00.000Z" } }`))

			// when
			res := h.Decode(&stdin)

			// then
			tr := res.(Transaction)
			assert.Equal(t, "Acme Corporation", tr.Merchant)
			assert.Equal(t, 20, tr.Amount)
			assert.NotEmpty(t, tr.Time)
		},
		"Should not decode unknown payload": func(t *testing.T) {
			// given
			h := Handler{}
			var stdin bytes.Buffer
			stdin.Write([]byte(`{ "unknown": { "command": "here" } }`))

			// when
			res := h.Decode(&stdin)

			// then
			assert.Empty(t, res)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}

func TestEncode(t *testing.T) {
	tests := map[string]func(*testing.T){
		"Should encode response": func(t *testing.T) {
			// given
			h := Handler{}
			acc := Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			}
			errs := []error{
				errors.New("this-is-an-error"),
			}

			// when
			res := h.Encode(acc, errs)

			// then
			assert.JSONEq(t, `{"account":{"activeCard":true,"availableLimit":100},"violations":["this-is-an-error"]}`, res.String())
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}

func TestDispatch(t *testing.T) {
	// setup
	dbMock := NewDatabaseMock()
	accMock := &accountHandlerMock{}
	h := Handler{
		db:             dbMock,
		accountHandler: accMock,
	}

	tests := map[string]func(*testing.T){
		"Should dispatch initialize account request": func(t *testing.T) {
			// given
			acc := Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			}
			accMock.On("Initialize", acc).Return(acc, nil)

			// when
			res, errs := h.Dispatch(acc)

			// then
			accMock.AssertNumberOfCalls(t, "Initialize", 1)
			assert.Equal(t, acc, res)
			assert.Empty(t, errs)
		},
		"Should dispatch authorize transaction request": func(t *testing.T) {
			// given
			acc := Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			}
			tr := Transaction{
				Merchant: "Acme Corporation",
				Amount:   100,
				Time:     time.Now(),
			}
			dbMock.On("CurrentAccount").Return(acc, nil)
			accMock.On("Authorize", acc, tr)

			// when
			res, errs := h.Dispatch(tr)

			// then
			accMock.AssertNumberOfCalls(t, "Authorize", 1)
			assert.Equal(t, acc, res)
			assert.Empty(t, errs)
		},
		"Should reach fallback for unknown request": func(t *testing.T) {
			// given
			acc := Account{
				ActiveCard:     true,
				AvailableLimit: 100,
			}
			dbMock.On("CurrentAccount").Return(acc, nil)

			// when
			res, errs := h.Dispatch(nil)

			// then
			assert.Equal(t, acc, res)
			assert.Empty(t, errs)
		},
	}

	for name, run := range tests {
		t.Run(name, func(t *testing.T) {
			run(t)
		})
	}
}

type accountHandlerMock struct {
	mock.Mock
}

func (h *accountHandlerMock) Initialize(acc Account) (Account, []error) {
	_ = h.Called(acc)
	return acc, nil
}

func (h *accountHandlerMock) Authorize(acc Account, tr Transaction) (Account, []error) {
	_ = h.Called(acc, tr)
	return acc, nil
}
