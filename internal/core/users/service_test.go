package users_test

import (
	"context"
	"testing"

	"github.com/ivmello/kakebo-go-api/internal/core/users"
	"github.com/ivmello/kakebo-go-api/internal/core/users/dto"
	"github.com/ivmello/kakebo-go-api/internal/core/users/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	t.Run("Should create a user", func(t *testing.T) {
		input := dto.CreateUserInput{
			Name:     "John Doe",
			Email:    "john@doe.com",
			Password: "123456",
		}
		expectedOutput := dto.CreateUserOutput{
			ID:     1,
			Status: "created",
		}
		repository := new(users.RepositoryMock)
		repository.On("SaveUser", context.Background(), mock.Anything).Return(1, nil)
		service := users.NewService(repository)
		output, err := service.CreateUser(context.Background(), input)
		assert.Nil(t, err)
		assert.Equal(t, expectedOutput, output)
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("Should update a user", func(t *testing.T) {
		userId := 2
		input := dto.UpdateUserInput{
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
		expectedOutput := dto.UpdateUserOutput{
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
