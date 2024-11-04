package models

import (
	"context"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

type RegisterCredentials struct {
	FirstName   string `json:"first_name" validate:"required,min=2,max=255"`
	LastName    string `json:"last_name" validate:"required,min=2,max=255"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Email       string `json:"email" validate:"required,email,max=255"`
	Password    string `json:"password" validate:"required,min=6,max=255"`
}

type LoginCredentials struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

type AuthRepository interface {
	RegisterUser(ctx context.Context, registerData *RegisterCredentials) (*User, error)
	GetUser(ctx context.Context, query interface{}, args ...interface{}) (*User, error)
}

type AuthService interface {
	Login(ctx context.Context, loginData *LoginCredentials) (string, *User, error)
	Register(ctx context.Context, registerData *RegisterCredentials) (string, *User, error)
}

func MatchesHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
