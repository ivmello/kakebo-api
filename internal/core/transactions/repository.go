package transactions

import (
	"context"

	"github.com/ivmello/kakebo-go-api/internal/adapters/database"
)

type Repository interface {
	SaveTransaction(ctx context.Context, transaction *Transaction) (string, error)
	GetAllUserTransactions(ctx context.Context, userId int) ([]*Transaction, error)
	GetTransactionById(ctx context.Context, userId int, transactionId string) (*Transaction, error)
	DeleteTransaction(ctx context.Context, transactionId string) error
}

type repo struct {
	conn database.Connection
}

func NewRepository(conn database.Connection) Repository {
	return &repo{
		conn,
	}
}

func (r *repo) SaveTransaction(ctx context.Context, transaction *Transaction) (transactionId string, err error) {
	err = r.conn.QueryRow(ctx,
		"INSERT INTO transactions (id, user_id, amount, transaction_type, description, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		transaction.ID, transaction.UserID, transaction.Amount, string(transaction.TransactionType), transaction.Description, transaction.CreatedAt).Scan(&transactionId)
	if err != nil {
		return "", err
	}
	return transactionId, nil
}

func (r *repo) DeleteTransaction(ctx context.Context, transactionId string) error {
	_, err := r.conn.Exec(ctx, "DELETE FROM transactions WHERE id = $1", transactionId)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) GetAllUserTransactions(ctx context.Context, userId int) ([]*Transaction, error) {
	rows, err := r.conn.Query(ctx, "SELECT id, user_id, amount, transaction_type, description, created_at FROM transactions WHERE user_id = $1 ORDER BY created_at DESC", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	transactions := make([]*Transaction, 0)
	for rows.Next() {
		transaction := &Transaction{}
		err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.TransactionType, &transaction.Description, &transaction.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (r *repo) GetTransactionById(ctx context.Context, userId int, transactionId string) (*Transaction, error) {
	transaction := &Transaction{}
	err := r.conn.QueryRow(ctx, "SELECT id, user_id, amount, transaction_type, description, created_at FROM transactions WHERE user_id = $1 and id = $2", userId, transactionId).
		Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.TransactionType, &transaction.Description, &transaction.CreatedAt)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}
