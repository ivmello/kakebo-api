package entity

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	ADMIN   Role = "admin"
	REGULAR Role = "regular"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUser(id int, name, email, password string) *User {
	return &User{
		ID:        id,
		Name:      name,
		Email:     email,
		Password:  HashPassword(password),
		CreatedAt: time.Now(),
	}
}

func (u *User) UpdatePassword(password string) {
	u.Password = HashPassword(password)
}

func HashPassword(password string) string {
	if password == "" {
		return ""
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(encryptedPassword)
}

func CheckPassword(userPassword, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(inputPassword))
	return err == nil
}
