package transactions

import (
	"context"
	"time"

	"github.com/ivmello/kakebo-go-api/internal/core/transactions/dto"
	"github.com/ivmello/kakebo-go-api/internal/core/transactions/entity"
	"github.com/ivmello/kakebo-go-api/internal/utils"
)

type Service interface {
	CreateTransaction(ctx context.Context, input dto.CreateTransactionInput) (dto.CreateTransactionOutput, error)
	GetAllUserTransactions(ctx context.Context, userId int) ([]dto.TransactionOutput, error)
	GetTransaction(ctx context.Context, userId int, transactionId string) (dto.TransactionOutput, error)
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

func (s *service) CreateTransaction(ctx context.Context, input dto.CreateTransactionInput) (dto.CreateTransactionOutput, error) {
	userId := ctx.Value(utils.USER_ID_KEY).(int)
	transaction := entity.NewTransaction("", userId, input.Amount, entity.TransactionType(input.TransactionType), input.Description)
	transactionId, err := s.repo.SaveTransaction(ctx, transaction)
	if err != nil {
		return dto.CreateTransactionOutput{}, err
	}
	return dto.CreateTransactionOutput{
		ID:     transactionId,
		Status: "created",
	}, nil
}

func (s *service) GetAllUserTransactions(ctx context.Context, userId int) ([]dto.TransactionOutput, error) {
	transactions, err := s.repo.GetAllUserTransactions(ctx, userId)
	if err != nil {
		return nil, err
	}
	output := make([]dto.TransactionOutput, 0)
	for _, transaction := range transactions {
		output = append(output, dto.TransactionOutput{
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

func (s *service) GetTransaction(ctx context.Context, userId int, transactionId string) (dto.TransactionOutput, error) {
	transaction, _ := s.repo.GetTransactionById(ctx, userId, transactionId)
	if transaction == nil {
		return dto.TransactionOutput{}, ErrTransactionNotFound
	}
	output := dto.TransactionOutput{
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
