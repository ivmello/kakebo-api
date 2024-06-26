package users_test

import (
	"context"
	"testing"

	"github.com/ivmello/kakebo-go-api/internal/core/users"
	"github.com/ivmello/kakebo-go-api/internal/core/users/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	t.Run("Should create a user", func(t *testing.T) {
		input := users.CreateUserInput{
			Name:     "John Doe",
			Email:    "john@doe.com",
			Password: "123456",
		}
		expectedOutput := users.CreateUserOutput{
			ID:     1,
			Status: "created",
		}
		ctx := context.Background()
		repository := new(users.RepositoryMock)
		repository.On("GetUserByEmail", ctx, input.Email).Return(nil, nil)
		repository.On("SaveUser", ctx, mock.Anything).Return(1, nil)
		service := users.NewService(repository)
		output, err := service.CreateUser(ctx, input)
		assert.Nil(t, err)
		assert.Equal(t, expectedOutput, output)
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("Should update a user", func(t *testing.T) {
		userId := 2
		input := users.UpdateUserInput{
			Name:     "John Doe 2",
			Email:    "john2@doe.com",
			Password: "1234567",
		}
		existingUser := &entity.User{
			ID:       userId,
			Name:     "John Doe",
			Email:    "john2@doe.com",
			Password: "1234567",
		}
		expectedOutput := users.UpdateUserOutput{
			ID:     2,
			Status: "updated",
		}
		repository := new(users.RepositoryMock)
		repository.On("GetUserByID", context.Background(), mock.Anything).Return(existingUser, nil)
		repository.On("UpdateUser", context.Background(), mock.Anything).Return(nil)
		service := users.NewService(repository)
		output, err := service.UpdateUser(context.Background(), userId, input)
		assert.Nil(t, err)
		assert.Equal(t, expectedOutput, output)
	})
}
