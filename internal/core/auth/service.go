package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ivmello/kakebo-go-api/internal/core/auth/dto"
	"github.com/ivmello/kakebo-go-api/internal/core/users"
	"github.com/ivmello/kakebo-go-api/internal/core/users/entity"
)

type Service interface {
	Login(input dto.LoginInput) (dto.LoginOutput, error)
	VerifyToken(tokenString string) error
	GetUserFromToken(tokenString string) (*entity.User, error)
}

type service struct {
	JWTSecret string
	repo      users.Repository
}

func NewService(jwtSecret string, repo users.Repository) Service {
	return &service{
		JWTSecret: jwtSecret,
		repo:      repo,
	}
}

func (s *service) createToken(email, password string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    email,
		"password": password,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *service) VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.JWTSecret), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil
}

func (s *service) Login(input dto.LoginInput) (dto.LoginOutput, error) {
	ctx := context.Background()
	user, err := s.repo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return dto.LoginOutput{}, ErrInvalidUser
	}
	checked := entity.CheckPassword(user.Password, input.Password)
	if !checked {
		return dto.LoginOutput{}, ErrInvalidUser
	}
	token, err := s.createToken(input.Email, input.Password)
	if err != nil {
		return dto.LoginOutput{}, err
	}
	return dto.LoginOutput{
		Token: token,
	}, nil
}

func (s *service) GetUserFromToken(tokenString string) (*entity.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}
	email, ok := claims["email"].(string)
	if !ok {
		return nil, ErrInvalidUser
	}
	password, ok := claims["password"].(string)
	if !ok {
		return nil, ErrInvalidUser
	}
	ctx := context.Background()
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	checked := entity.CheckPassword(user.Password, password)
	if !checked {
		return nil, ErrInvalidUser
	}
	return user, nil
}
