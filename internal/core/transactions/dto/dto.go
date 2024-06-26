package dto

type CreateTransactionInput struct {
	UserID          int    `json:"user_id"`
	Amount          int    `json:"amount"`
	TransactionType string `json:"transaction_type"`
	Description     string `json:"description"`
}

type CreateTransactionOutput struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type TransactionOutput struct {
	ID              string `json:"id"`
	UserID          int    `json:"user_id"`
	Amount          int    `json:"amount"`
	TransactionType string `json:"transaction_type"`
	Description     string `json:"description"`
	CreatedAt       string `json:"created_at"`
}

type GetAllUserTransactionsOutput struct {
	Transactions []TransactionOutput `json:"transactions"`
}
