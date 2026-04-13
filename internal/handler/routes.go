package handler

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterAllRoutes(h *handlers, api huma.API) {
	accountsTag := []string{"Accounts"}
	transactionsTag := []string{"Transactions"}

	huma.Register(api, huma.Operation{ //nolint: exhaustruct
		Method:        http.MethodPost,
		Path:          "/accounts",
		Tags:          accountsTag,
		Description:   "Create a new account",
		DefaultStatus: http.StatusCreated,
	}, h.accountHandler.HandleCreateAccount)

	huma.Register(api, huma.Operation{ //nolint: exhaustruct
		Method:      http.MethodGet,
		Path:        "/accounts/{account_id}",
		Tags:        accountsTag,
		Description: "Retrieve account by ID",
	}, h.accountHandler.HandleGetAccountByID)

	huma.Register(api, huma.Operation{ //nolint: exhaustruct
		Method:        http.MethodPost,
		Path:          "/transactions",
		Tags:          transactionsTag,
		Description:   "Create transaction",
		DefaultStatus: http.StatusCreated,
	}, h.transactionHandler.HandleCreateTransaction)
}
