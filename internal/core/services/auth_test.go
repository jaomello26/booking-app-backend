//go:build test

package services

import (
	"context"
	"testing"

	"github.com/jaomello26/booking-app-backend/internal/core/models"
	"github.com/jaomello26/booking-app-backend/internal/core/models/mocks"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestAuthService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	service := NewAuthService(mockRepo)
	ctx := context.Background()

	registerData := &models.AuthCredentials{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	mockRepo.EXPECT().
		GetUser(ctx, "email = ?", registerData.Email).
		Return(nil, gorm.ErrRecordNotFound)

	mockRepo.EXPECT().
		RegisterUser(ctx, registerData).
		Return(&models.User{
			ID:    1,
			Name:  registerData.Name,
			Email: registerData.Email,
		}, nil)

	token, user, err := service.Register(ctx, registerData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if token == "" {
		t.Error("Expected token to be generated")
	}

	if user.ID != 1 {
		t.Errorf("Expected user ID 1, got %d", user.ID)
	}
}

func TestAuthService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	service := NewAuthService(mockRepo)
	ctx := context.Background()

	loginData := &models.AuthCredentials{
		Email:    "john@example.com",
		Password: "password123",
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(loginData.Password), bcrypt.DefaultCost)

	mockRepo.EXPECT().
		GetUser(ctx, "email = ?", loginData.Email).
		Return(&models.User{
			ID:       1,
			Email:    loginData.Email,
			Password: string(hashedPassword),
		}, nil)

	token, user, err := service.Login(ctx, loginData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if token == "" {
		t.Error("Expected token to be generated")
	}

	if user.ID != 1 {
		t.Errorf("Expected user ID 1, got %d", user.ID)
	}
}
