package transactions

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	SaveTransaction(ctx context.Context, transaction *Transaction) (string, error)
	GetAllUserTransactions(ctx context.Context, userId int) ([]*Transaction, error)
	GetAllUserTransactionsByFilter(ctx context.Context, userId int, input TransactionFilter) ([]*Transaction, error)
	GetTransactionById(ctx context.Context, userId int, transactionId string) (*Transaction, error)
	DeleteTransaction(ctx context.Context, transactionId string) error
	SaveBulkTransactions(ctx context.Context, transactions []*Transaction) (int, error)
}

type repo struct {
	conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) Repository {
	return &repo{
		conn,
	}
}

func (r *repo) SaveTransaction(ctx context.Context, transaction *Transaction) (transactionId string, err error) {
	err = r.conn.QueryRow(ctx,
		"INSERT INTO transactions (id, user_id, amount, transaction_type, description, date) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		transaction.ID, transaction.UserID, transaction.Amount, string(transaction.TransactionType), transaction.Description, transaction.Date).Scan(&transactionId)
	if err != nil {
		return "", err
	}
	return transactionId, nil
}

func prepareData(transactions []*Transaction) [][]interface{} {
	rows := make([][]interface{}, len(transactions))
	for i, item := range transactions {
		rows[i] = []interface{}{
			item.ID,
			item.UserID,
			item.Description,
			item.Amount,
			item.Date.Format("2006-01-02 15:04:05"),
			string(item.TransactionType),
		}
	}
	return rows
}

func (r *repo) SaveBulkTransactions(ctx context.Context, transactions []*Transaction) (int, error) {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
	}()
	_, err = tx.Exec(ctx, "CREATE TEMPORARY TABLE IF NOT EXISTS _temp_upsert_transactions (LIKE transactions INCLUDING ALL) ON COMMIT DROP")
	if err != nil {
		return 0, err
	}
	data := prepareData(transactions)
	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"_temp_upsert_transactions"},
		[]string{
			"id",
			"user_id",
			"description",
			"amount",
			"date",
			"transaction_type",
		},
		pgx.CopyFromRows(data),
	)
	if err != nil {
		return 0, err
	}
	upsertQuery := `
		INSERT INTO transactions (id, user_id, description, amount, date, transaction_type)
		SELECT id, user_id, description, amount, date, transaction_type
		FROM _temp_upsert_transactions
		ON CONFLICT (id) DO NOTHING
		RETURNING id;
	`
	result, err := tx.Exec(ctx, upsertQuery)
	if err != nil {
		return 0, err
	}
	return int(result.RowsAffected()), nil
}

func (r *repo) DeleteTransaction(ctx context.Context, transactionId string) error {
	_, err := r.conn.Exec(ctx, "DELETE FROM transactions WHERE id = $1", transactionId)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) GetAllUserTransactions(ctx context.Context, userId int) ([]*Transaction, error) {
	rows, err := r.conn.Query(ctx, "SELECT id, user_id, amount, transaction_type, description, date FROM transactions WHERE user_id = $1 ORDER BY date DESC", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	transactions := make([]*Transaction, 0)
	for rows.Next() {
		transaction := &Transaction{}
		err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.TransactionType, &transaction.Description, &transaction.Date)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (r *repo) GetAllUserTransactionsByFilter(ctx context.Context, userId int, input TransactionFilter) ([]*Transaction, error) {
	stmt := `SELECT id, user_id, amount, transaction_type, description, date FROM transactions WHERE user_id = $1`
	if input.StartDate != "" {
		stmt += ` AND date >= '` + input.StartDate + `'`
	}
	if input.EndDate != "" {
		stmt += ` AND date <= '` + input.EndDate + `'`
	}
	if input.Month != 0 {
		stmt += ` AND EXTRACT(MONTH FROM date) = ` + fmt.Sprintf("%d", input.Month)
	}
	if input.Year != 0 {
		stmt += ` AND EXTRACT(YEAR FROM date) = ` + fmt.Sprintf("%d", input.Year)
	} else {
		stmt += ` AND EXTRACT(YEAR FROM date) = ` + fmt.Sprintf("%d", time.Now().Year())
	}
	rows, err := r.conn.Query(ctx, stmt, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	transactions := make([]*Transaction, 0)
	for rows.Next() {
		transaction := &Transaction{}
		err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.TransactionType, &transaction.Description, &transaction.Date)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (r *repo) GetTransactionById(ctx context.Context, userId int, transactionId string) (*Transaction, error) {
	transaction := &Transaction{}
	err := r.conn.QueryRow(ctx, "SELECT id, user_id, amount, transaction_type, description, date FROM transactions WHERE user_id = $1 and id = $2", userId, transactionId).
		Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.TransactionType, &transaction.Description, &transaction.Date)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}
