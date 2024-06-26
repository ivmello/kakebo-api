package users

import (
	"context"

	"github.com/ivmello/kakebo-go-api/internal/adapters/database"
	"github.com/ivmello/kakebo-go-api/internal/core/users/entity"
)

type Repository interface {
	SaveUser(ctx context.Context, user *entity.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id int) error
	ListUsers(ctx context.Context) ([]*entity.User, error)
	GetUserByEmailAndPassword(ctx context.Context, email, password string) (*entity.User, error)
}

type repo struct {
	conn database.Connection
}

func NewRepository(conn database.Connection) Repository {
	return &repo{
		conn,
	}
}

func (r *repo) SaveUser(ctx context.Context, user *entity.User) (int, error) {
	var userID int
	err := r.conn.QueryRow(ctx,
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		user.Name, user.Email, user.Password).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (r *repo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}
	err := r.conn.QueryRow(ctx, "SELECT id, name, email, password FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repo) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	user := &entity.User{}
	err := r.conn.QueryRow(ctx, "SELECT id, name, email, password FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repo) UpdateUser(ctx context.Context, user *entity.User) error {
	_, err := r.conn.Exec(ctx, "UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4",
		user.Name, user.Email, user.Password, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) DeleteUser(ctx context.Context, id int) error {
	_, err := r.conn.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) ListUsers(ctx context.Context) ([]*entity.User, error) {
	rows, err := r.conn.Query(ctx, "SELECT id, name, email, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*entity.User, 0)
	for rows.Next() {
		user := &entity.User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *repo) GetUserByEmailAndPassword(ctx context.Context, email, password string) (*entity.User, error) {
	user := &entity.User{}
	err := r.conn.QueryRow(ctx, "SELECT id, name, email, password FROM users WHERE email = $1 AND password = $2",
		email, password).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
