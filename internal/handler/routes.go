package handler

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterAllRoutes(h *handlers, api huma.API) {
	accountsTag := []string{"Accounts"}
	transactionsTag := []string{"Transactions"}

	huma.Register(api, huma.Operation{ //nolint: exhaustruct
		Method: http.MethodPost,
		Path:   "/accounts",
		Tags:   accountsTag,
		Description: `## Create a new account

Creates a new account with the provided document number.

The document number cannot be a empty string.

**Error Codes:**
- 422 Unprocessable Entity: Invalid input data
- 409 Conflict: Document number already exists
`,
		Errors:        []int{http.StatusUnprocessableEntity, http.StatusConflict},
		DefaultStatus: http.StatusCreated,
	}, h.accountHandler.HandleCreateAccount)

	huma.Register(api, huma.Operation{ //nolint: exhaustruct
		Method: http.MethodGet,
		Path:   "/accounts/{account_id}",
		Tags:   accountsTag,
		Description: `## Retrieve account by ID

Retrieves the account details for the specified account ID.

**Error Codes:**
- 404 Not Found: No account exists with the provided ID
`,
		Errors: []int{http.StatusNotFound},
	}, h.accountHandler.HandleGetAccountByID)

	huma.Register(api, huma.Operation{ //nolint: exhaustruct
		Method: http.MethodPost,
		Path:   "/transactions",
		Tags:   transactionsTag,
		Description: `## Create a transaction

Creates a new transaction for a specified account. The transaction includes the account ID, operation type ID, and the amount. The amount must be always a positive number.

**Error Codes:**
- 422 Unprocessable Entity: Invalid input data (e.g., invalid account ID, invalid operation type ID, or invalid amount)
- 400 Bad Request: Invalid account ID or operation type ID
`,
		DefaultStatus: http.StatusCreated,
		Errors:        []int{http.StatusUnprocessableEntity, http.StatusBadRequest},
	}, h.transactionHandler.HandleCreateTransaction)
}
