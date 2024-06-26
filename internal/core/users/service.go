package users

import (
	"context"
)

type Service interface {
	CreateUser(ctx context.Context, input CreateUserInput) (CreateUserOutput, error)
	UpdateUser(ctx context.Context, id int, input UpdateUserInput) (UpdateUserOutput, error)
	GetUser(ctx context.Context, id int) (UserOutput, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo,
	}
}

func (s *service) CreateUser(ctx context.Context, input CreateUserInput) (CreateUserOutput, error) {
	user, _ := s.repo.GetUserByEmail(ctx, input.Email)
	if user != nil {
		return CreateUserOutput{}, ErrUserAlreadyExists
	}
	newUser := NewUser(0, input.Name, input.Email, input.Password)
	userId, err := s.repo.SaveUser(ctx, newUser)
	if err != nil {
		return CreateUserOutput{}, err
	}
	return CreateUserOutput{
		ID:     userId,
		Status: "created",
	}, nil
}

func (s *service) UpdateUser(ctx context.Context, id int, input UpdateUserInput) (UpdateUserOutput, error) {
	if id <= 0 {
		return UpdateUserOutput{}, ErrInvalidUserID
	}
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return UpdateUserOutput{}, err
	}
	user.Name = input.Name
	user.Email = input.Email
	user.UpdatePassword(input.Password)
	err = s.repo.UpdateUser(ctx, user)
	if err != nil {
		return UpdateUserOutput{}, err
	}
	return UpdateUserOutput{
		ID:     user.ID,
		Status: "updated",
	}, nil
}

func (s *service) GetUser(ctx context.Context, id int) (UserOutput, error) {
	if id <= 0 {
		return UserOutput{}, ErrInvalidUserID
	}
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return UserOutput{}, err
	}
	return UserOutput{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
