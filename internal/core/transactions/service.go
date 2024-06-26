package transactions

import (
	"context"
	"time"
)

type Service interface {
	CreateTransaction(ctx context.Context, userId int, input CreateTransactionInput) (CreateTransactionOutput, error)
	GetAllUserTransactions(ctx context.Context, userId int) ([]TransactionOutput, error)
	GetTransaction(ctx context.Context, userId int, transactionId string) (TransactionOutput, error)
	DeleteTransaction(ctx context.Context, userId int, transactionId string) error
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
			CreatedAt:       transaction.CreatedAt.Local().Format(time.RFC3339),
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
		CreatedAt:       transaction.CreatedAt.Local().Format(time.RFC3339),
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
