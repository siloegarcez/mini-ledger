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
	accountRepository := repository.NewAccountsRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)
	operationTypeRepository := repository.NewOperationTypeRepository(db)
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
