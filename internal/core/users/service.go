package users

import (
	"context"
	"fmt"

	"github.com/ivmello/kakebo-go-api/internal/core/users/dto"
	"github.com/ivmello/kakebo-go-api/internal/core/users/entity"
)

type Service interface {
	CreateUser(ctx context.Context, input dto.CreateUserInput) (dto.CreateUserOutput, error)
	UpdateUser(ctx context.Context, id int, input dto.UpdateUserInput) (dto.UpdateUserOutput, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo,
	}
}

func (s *service) CreateUser(ctx context.Context, input dto.CreateUserInput) (dto.CreateUserOutput, error) {
	user := entity.NewUser(0, input.Name, input.Email, input.Password)
	_, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return dto.CreateUserOutput{}, ErrUserAlreadyExists
	}
	fmt.Println("Creating user")
	userId, err := s.repo.SaveUser(ctx, user)
	if err != nil {
		return dto.CreateUserOutput{}, err
	}
	return dto.CreateUserOutput{
		ID:     userId,
		Status: "created",
	}, nil
}

func (s *service) UpdateUser(ctx context.Context, id int, input dto.UpdateUserInput) (dto.UpdateUserOutput, error) {
	if id <= 0 {
		return dto.UpdateUserOutput{}, ErrInvalidUserID
	}
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return dto.UpdateUserOutput{}, err
	}
	user.Name = input.Name
	user.Email = input.Email
	user.UpdatePassword(input.Password)
	err = s.repo.UpdateUser(ctx, user)
	if err != nil {
		return dto.UpdateUserOutput{}, err
	}
	return dto.UpdateUserOutput{
		ID:     user.ID,
		Status: "updated",
	}, nil
}
