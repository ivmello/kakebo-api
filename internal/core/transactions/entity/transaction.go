package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID          string
	UserID      int
	Amount      int
	Description string
	CreatedAt   time.Time
}

func NewTransaction(id string, userId, amount int, description string) *Transaction {
	sid := id
	if sid == "" {
		sid = uuid.New().String()
	}
	return &Transaction{
		ID:          sid,
		UserID:      userId,
		Amount:      amount,
		Description: description,
		CreatedAt:   time.Now(),
	}
}
