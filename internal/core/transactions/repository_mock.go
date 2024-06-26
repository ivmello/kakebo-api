package transactions

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) SaveTransaction(ctx context.Context, transaction *Transaction) (string, error) {
	args := m.Called(ctx, transaction)
	return args.String(0), args.Error(1)
}

func (m *RepositoryMock) GetAllUserTransactions(ctx context.Context, userId int) ([]*Transaction, error) {
	args := m.Called(ctx, userId)
	transaction := args.Get(0)
	if transaction == nil {
		return nil, args.Error(1)
	}
	return transaction.([]*Transaction), args.Error(1)
}

func (m *RepositoryMock) GetTransactionById(ctx context.Context, userId int, transactionId string) (*Transaction, error) {
	args := m.Called(ctx, userId, transactionId)
	transaction := args.Get(0)
	if transaction == nil {
		return nil, args.Error(1)
	}
	return transaction.(*Transaction), args.Error(1)
}

func (m *RepositoryMock) DeleteTransaction(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *RepositoryMock) GetAllUserTransactionsByFilter(ctx context.Context, userId int, input TransactionFilter) ([]*Transaction, error) {
	args := m.Called(ctx, userId, input)
	transaction := args.Get(0)
	if transaction == nil {
		return nil, args.Error(1)
	}
	return transaction.([]*Transaction), args.Error(1)
}
