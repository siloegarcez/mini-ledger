package handler_test

import (
	"mini-ledger/internal/handler"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateTransaction(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Integration Test: Test_CreateTransaction")
	}
	api, cleanup := allocateResources(t)
	defer cleanup()

	documentNumber := "12345678900"

	// Create an account
	createResp := api.Post("/accounts", handler.AccountCreateRequest{
		DocumentNumber: documentNumber,
	})
	assert.Equal(t, http.StatusCreated, createResp.Code)

	assert.True(t, strings.Contains(createResp.Body.String(), documentNumber))
}

func Test_CreateTransactionNonExistentAccount(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Integration Test: Test_CreateTransactionNonExistentAccount")
	}
	api, cleanup := allocateResources(t)
	defer cleanup()

	createResp := api.Post("/transactions", handler.TransactionCreateRequest{
		AccountID:       999, // Non-existent account ID
		OperationTypeID: 1,
		Amount:          "100.00",
	})
	assert.Equal(t, http.StatusBadRequest, createResp.Code)
}

func Test_CreateTransactionInvalidOperationType(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Integration Test: Test_CreateTransactionInvalidOperationType")
	}
	api, cleanup := allocateResources(t)
	defer cleanup()

	documentNumber := "12345678900"

	// Create an account
	createResp := api.Post("/accounts", handler.AccountCreateRequest{
		DocumentNumber: documentNumber,
	})
	assert.Equal(t, http.StatusCreated, createResp.Code)

	assert.True(t, strings.Contains(createResp.Body.String(), documentNumber))

	createTransacResp := api.Post("/transactions", handler.TransactionCreateRequest{
		AccountID:       1,
		OperationTypeID: 999, // Non-existent operation type ID
		Amount:          "100.00",
	})
	assert.Equal(t, http.StatusBadRequest, createTransacResp.Code)
}

func Test_CreateTransactionNegativeAmount(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Integration Test: Test_CreateTransactionNegativeAmount")
	}
	api, cleanup := allocateResources(t)
	defer cleanup()

	documentNumber := "12345678900"

	// Create an account
	createResp := api.Post("/accounts", handler.AccountCreateRequest{
		DocumentNumber: documentNumber,
	})
	assert.Equal(t, http.StatusCreated, createResp.Code)

	assert.True(t, strings.Contains(createResp.Body.String(), documentNumber))

	createTransacResp := api.Post("/transactions", handler.TransactionCreateRequest{
		AccountID:       1,
		OperationTypeID: 1,
		Amount:          "-100.00", // Negative amount
	})
	assert.Equal(t, http.StatusUnprocessableEntity, createTransacResp.Code)
}

func Test_CreateDebitTransaction(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Integration Test: Test_CreateDebitTransaction")
	}
	api, cleanup := allocateResources(t)
	defer cleanup()

	documentNumber := "12345678900"

	// Create an account
	createResp := api.Post("/accounts", handler.AccountCreateRequest{
		DocumentNumber: documentNumber,
	})
	assert.Equal(t, http.StatusCreated, createResp.Code)

	assert.True(t, strings.Contains(createResp.Body.String(), documentNumber))

	createTransacResp := api.Post("/transactions", handler.TransactionCreateRequest{
		AccountID:       1,
		OperationTypeID: 1, // Assuming this is a debit operation type ID
		Amount:          "100.00",
	})
	assert.Equal(t, http.StatusCreated, createTransacResp.Code)

	assert.True(t, strings.Contains(createTransacResp.Body.String(), "-100.00"))
}

func Test_CreateCreditTransaction(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Integration Test: Test_CreateDebitTransaction")
	}
	api, cleanup := allocateResources(t)
	defer cleanup()

	documentNumber := "12345678900"

	// Create an account
	createResp := api.Post("/accounts", handler.AccountCreateRequest{
		DocumentNumber: documentNumber,
	})
	assert.Equal(t, http.StatusCreated, createResp.Code)

	assert.True(t, strings.Contains(createResp.Body.String(), documentNumber))

	createTransacResp := api.Post("/transactions", handler.TransactionCreateRequest{
		AccountID:       1,
		OperationTypeID: 4,
		Amount:          "123.45",
	})
	assert.Equal(t, http.StatusCreated, createTransacResp.Code)

	assert.True(t, strings.Contains(createTransacResp.Body.String(), "123.45"))
	assert.False(t, strings.Contains(createTransacResp.Body.String(), "-123.45"))
}
