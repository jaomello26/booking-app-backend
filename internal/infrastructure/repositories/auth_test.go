//go:build test

package repositories

import (
	"context"
	"testing"

	"github.com/jaomello26/booking-app-backend/internal/core/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestAuthRepository(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestAuthRepository_RegisterUser(t *testing.T) {
	db := setupTestAuthRepository(t)
	repo := NewAuthRepository(db)
	ctx := context.Background()

	userData := &models.AuthCredentials{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	user, err := repo.RegisterUser(ctx, userData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.ID == 0 {
		t.Errorf("Expected user ID to be set")
	}

	if user.Email != userData.Email {
		t.Errorf("Expected email %s, got %s", userData.Email, user.Email)
	}
}

func TestAuthRepository_GetUser(t *testing.T) {
	db := setupTestAuthRepository(t)
	repo := NewAuthRepository(db)
	ctx := context.Background()

	userData := &models.AuthCredentials{
		Name:     "Jane Doe",
		Email:    "jane@example.com",
		Password: "password123",
	}
	_, err := repo.RegisterUser(ctx, userData)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	user, err := repo.GetUser(ctx, "email = ?", userData.Email)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Email != userData.Email {
		t.Errorf("Expected email %s, got %s", userData.Email, user.Email)
	}
}
