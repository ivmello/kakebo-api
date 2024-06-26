package users

import (
	"context"

	"github.com/ivmello/kakebo-go-api/internal/core/users/entity"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) SaveUser(ctx context.Context, user *entity.User) (int, error) {
	args := m.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func (m *RepositoryMock) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *RepositoryMock) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *RepositoryMock) UpdateUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *RepositoryMock) DeleteUser(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *RepositoryMock) ListUsers(ctx context.Context) ([]*entity.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *RepositoryMock) GetUserByEmailAndPassword(ctx context.Context, email, password string) (*entity.User, error) {
	args := m.Called(ctx, email, password)
	return args.Get(0).(*entity.User), args.Error(1)
}
