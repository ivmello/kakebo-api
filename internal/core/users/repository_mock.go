package users

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) SaveUser(ctx context.Context, user *User) (int, error) {
	args := m.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func (m *RepositoryMock) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*User), args.Error(1)
}

func (m *RepositoryMock) GetUserByID(ctx context.Context, id int) (*User, error) {
	args := m.Called(ctx, id)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*User), args.Error(1)
}

func (m *RepositoryMock) UpdateUser(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	userArg := args.Get(0)
	if userArg == nil {
		return args.Error(0)
	}
	return args.Error(0)
}

func (m *RepositoryMock) DeleteUser(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *RepositoryMock) ListUsers(ctx context.Context) ([]*User, error) {
	args := m.Called(ctx)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.([]*User), args.Error(1)
}

func (m *RepositoryMock) GetUserByEmailAndPassword(ctx context.Context, email, password string) (*User, error) {
	args := m.Called(ctx, email, password)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*User), args.Error(1)
}
