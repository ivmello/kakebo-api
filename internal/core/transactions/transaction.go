package transactions

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	DEBIT  TransactionType = "debit"
	CREDIT TransactionType = "credit"
)

type TransactionFilter struct {
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
	Month     int    `json:"month,omitempty"`
	Year      int    `json:"year,omitempty"`
}

type Transaction struct {
	ID              string          `json:"id"`
	UserID          int             `json:"user_id"`
	Amount          int             `json:"amount"`
	TransactionType TransactionType `json:"transaction_type"`
	Description     string          `json:"description"`
	Date            time.Time       `json:"date"`
}

func NewTransaction(id string, userId, amount int, transactionType TransactionType, description string) *Transaction {
	sid := id
	if sid == "" {
		sid = uuid.New().String()
	}
	return &Transaction{
		ID:              sid,
		UserID:          userId,
		Amount:          amount,
		TransactionType: transactionType,
		Description:     description,
		Date:            time.Now(),
	}
}
