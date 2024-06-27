package transactions

import (
	"context"
	"encoding/csv"
	"io"
	"strconv"
	"time"
)

type Service interface {
	CreateTransaction(ctx context.Context, userId int, input CreateTransactionInput) (CreateTransactionOutput, error)
	GetAllUserTransactions(ctx context.Context, userId int) ([]TransactionOutput, error)
	GetTransaction(ctx context.Context, userId int, transactionId string) (TransactionOutput, error)
	DeleteTransaction(ctx context.Context, userId int, transactionId string) error
	ImportTransactionsFromCSV(ctx context.Context, userId int, file io.Reader) (CreateBulkTransactionOutput, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo,
	}
}

func (s *service) CreateTransaction(ctx context.Context, userId int, input CreateTransactionInput) (CreateTransactionOutput, error) {
	transaction := NewTransaction("", userId, input.Amount, TransactionType(input.TransactionType), input.Description)
	transactionId, err := s.repo.SaveTransaction(ctx, transaction)
	if err != nil {
		return CreateTransactionOutput{}, err
	}
	return CreateTransactionOutput{
		ID:     transactionId,
		Status: "created",
	}, nil
}

func (s *service) parseCSV(file io.Reader, userId int) ([]*Transaction, error) {
	reader := csv.NewReader(file)
	var transactions []*Transaction
	_, err := reader.Read()
	if err != nil {
		return nil, err
	}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		dateStr := record[0]
		date, err := time.Parse("02/01/2006", dateStr)
		if err != nil {
			return nil, err
		}
		amountStr := record[1]
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return nil, err
		}
		amountInt := int(amount * 100)
		transactionType := CREDIT
		if amountInt < 0 {
			transactionType = DEBIT
		}
		identifier := record[2]
		description := record[3]
		transaction := NewTransaction(identifier, userId, amountInt, transactionType, description)
		transaction.Date = date
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (s *service) ImportTransactionsFromCSV(ctx context.Context, userId int, file io.Reader) (CreateBulkTransactionOutput, error) {
	output := CreateBulkTransactionOutput{
		Qnt:    0,
		Status: "no action",
	}
	transactions, err := s.parseCSV(file, userId)
	if err != nil {
		return output, err
	}
	result, err := s.repo.SaveBulkTransactions(ctx, transactions)
	if err != nil {
		return output, err
	}
	output.Qnt = result
	if result > 0 {
		output.Status = "created"
	}
	return output, nil
}

func (s *service) GetAllUserTransactions(ctx context.Context, userId int) ([]TransactionOutput, error) {
	transactions, err := s.repo.GetAllUserTransactions(ctx, userId)
	if err != nil {
		return nil, err
	}
	output := make([]TransactionOutput, 0)
	for _, transaction := range transactions {
		output = append(output, TransactionOutput{
			ID:              transaction.ID,
			UserID:          transaction.UserID,
			Amount:          transaction.Amount,
			TransactionType: string(transaction.TransactionType),
			Description:     transaction.Description,
			Date:            transaction.Date.Local().Format(time.RFC3339),
		})
	}
	return output, nil
}

func (s *service) GetTransaction(ctx context.Context, userId int, transactionId string) (TransactionOutput, error) {
	transaction, _ := s.repo.GetTransactionById(ctx, userId, transactionId)
	if transaction == nil {
		return TransactionOutput{}, ErrTransactionNotFound
	}
	output := TransactionOutput{
		ID:              transaction.ID,
		UserID:          transaction.UserID,
		Amount:          transaction.Amount,
		TransactionType: string(transaction.TransactionType),
		Description:     transaction.Description,
		Date:            transaction.Date.Local().Format(time.RFC3339),
	}
	return output, nil
}

func (s *service) DeleteTransaction(ctx context.Context, userId int, transactionId string) error {
	transaction, _ := s.repo.GetTransactionById(ctx, userId, transactionId)
	if transaction == nil {
		return ErrTransactionNotFound
	}
	return s.repo.DeleteTransaction(ctx, transaction.ID)
}
