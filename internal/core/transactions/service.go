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
	GetAllUserTransactions(ctx context.Context) (dto.GetAllUserTransactionsOutput, error)
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
	transaction := entity.NewTransaction("", userId, input.Amount, input.Description)
	transactionId, err := s.repo.SaveTransaction(ctx, transaction)
	if err != nil {
		return dto.CreateTransactionOutput{}, err
	}
	return dto.CreateTransactionOutput{
		ID:     transactionId,
		Status: "created",
	}, nil
}

func (s *service) GetAllUserTransactions(ctx context.Context) (dto.GetAllUserTransactionsOutput, error) {
	userId := ctx.Value(utils.USER_ID_KEY).(int)
	transactions, err := s.repo.GetAllUserTransactions(ctx, userId)
	if err != nil {
		return dto.GetAllUserTransactionsOutput{}, err
	}
	output := dto.GetAllUserTransactionsOutput{
		Transactions: make([]dto.TransactionOutput, 0),
	}
	for _, transaction := range transactions {
		output.Transactions = append(output.Transactions, dto.TransactionOutput{
			ID:          transaction.ID,
			UserID:      transaction.UserID,
			Amount:      transaction.Amount,
			Description: transaction.Description,
			CreatedAt:   transaction.CreatedAt.Local().Format(time.RFC3339),
		})
	}
	return output, nil
}
