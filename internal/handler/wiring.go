package handler

import (
	"database/sql"
	"mini-ledger/internal/repository"
	"mini-ledger/internal/service"
)

type handlers struct {
	accountHandler     AccountHandler
	transactionHandler TransactionHandler
}

func WireLayers(db *sql.DB) *handlers {
	accountRepository := repository.NewAccountsRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)
	operationTypeRepository := repository.NewOperationTypeRepository(db)

	accountService := service.NewAccountService(accountRepository)
	transactionService := service.NewTransactionService(
		transactionRepository,
		operationTypeRepository,
	)

	accountHandler := NewAccountHandler(accountService)
	transactionHandler := NewTransactionHandler(transactionService)

	return &handlers{
		accountHandler:     accountHandler,
		transactionHandler: transactionHandler,
	}
}
