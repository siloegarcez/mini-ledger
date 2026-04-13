package service

type (
	CreateTransactionCommand struct {
		AccountID       int64
		OperationTypeID int64
		Amount          string
	}
)
