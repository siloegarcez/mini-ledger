package handler

import (
	"database/sql"
	repository "mini-ledger/internal/repository/postgres"
	"mini-ledger/internal/service"
)

type handlers struct {
	accountHandler     AccountHandler
	transactionHandler TransactionHandler
}

func WireLayers(db *sql.DB) *handlers {
	executor := repository.NewExecutor(db)
	accountRepository := repository.NewAccountsRepository(executor)
	transactionRepository := repository.NewTransactionRepository(executor)
	operationTypeRepository := repository.NewOperationTypeRepository(executor)
	transactor := repository.NewTransactor(db)

	accountService := service.NewAccountService(accountRepository, transactor)
	transactionService := service.NewTransactionService(
		transactionRepository,
		operationTypeRepository,
		transactor,
	)

	accountHandler := NewAccountHandler(accountService)
	transactionHandler := NewTransactionHandler(transactionService)

	return &handlers{
		accountHandler:     accountHandler,
		transactionHandler: transactionHandler,
	}
}
