package transactions

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

type CreateBulkTransactionOutput struct {
	Qnt    int    `json:"qnt"`
	Status string `json:"status"`
}

type TransactionOutput struct {
	ID              string `json:"id"`
	UserID          int    `json:"user_id"`
	Amount          int    `json:"amount"`
	TransactionType string `json:"transaction_type"`
	Description     string `json:"description"`
	Date            string `json:"date"`
}
