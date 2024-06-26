package transactions

import (
	"context"

	"github.com/ivmello/kakebo-go-api/internal/core/transactions/entity"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) SaveTransaction(ctx context.Context, transaction *entity.Transaction) (string, error) {
	args := m.Called(ctx, transaction)
	return args.String(0), args.Error(1)
}

func (m *RepositoryMock) GetAllUserTransactions(ctx context.Context, userId int) ([]*entity.Transaction, error) {
	args := m.Called(ctx, userId)
	transaction := args.Get(0)
	if transaction == nil {
		return nil, args.Error(1)
	}
	return transaction.([]*entity.Transaction), args.Error(1)
}
