package service

import (
	"mini-ledger/internal/domain"
	"time"
)

type (
	CreateTransactionCommand struct {
		AccountID     int64
		OperationType domain.OperationType
		Amount        int64
		Currency      string
		Scale         int16
		EventDate     time.Time
	}
)
