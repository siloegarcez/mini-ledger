package handler_test

import (
	"mini-ledger/internal/domain"
	"mini-ledger/internal/handler"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateAccountGetAccount(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Integration Test: Test_CreateAccountGetAccount")
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

	getResp := api.Get("/accounts/1")
	assert.Equal(t, http.StatusOK, getResp.Code)
	assert.True(t, strings.Contains(getResp.Body.String(), documentNumber))
}

func Test_CreateAccountConflict(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Integration Test: Test_CreateAccountConflict")
	}
	api, cleanup := allocateResources(t)
	defer cleanup()

	documentNumber := "12345678900"

	// Create an account
	createResp := api.Post("/accounts", handler.AccountCreateRequest{
		DocumentNumber: documentNumber,
	})
	assert.Equal(t, http.StatusCreated, createResp.Code)

	// Try to create another account with the same document number
	conflictResp := api.Post("/accounts", handler.AccountCreateRequest{
		DocumentNumber: documentNumber,
	})
	assert.Equal(t, http.StatusConflict, conflictResp.Code)
}

func Test_GetNonExistentAccount(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Integration Test: Test_GetNonExistentAccount")
	}
	api, cleanup := allocateResources(t)
	defer cleanup()

	getResp := api.Get("/accounts/9999")
	assert.Equal(t, http.StatusNotFound, getResp.Code)
}

func Test_CreateEmptyDocumentAccount(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Integration Test: Test_CreateEmptyDocumentAccount")
	}
	api, cleanup := allocateResources(t)
	defer cleanup()

	createResp := api.Post("/accounts", handler.AccountCreateRequest{
		DocumentNumber: "",
	})
	assert.Equal(t, http.StatusUnprocessableEntity, createResp.Code)
}

func Test_CreateAccountInvalidDocument(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Integration Test: Test_CreateAccountInvalidDocument")
	}
	api, cleanup := allocateResources(t)
	defer cleanup()

	createResp := api.Post("/accounts", handler.AccountCreateRequest{
		DocumentNumber: strings.Repeat("9", domain.MaxDocumentNumberLength+1),
	})
	assert.Equal(t, http.StatusUnprocessableEntity, createResp.Code)
}
