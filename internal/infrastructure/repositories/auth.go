package repositories

import (
	"context"

	"github.com/jaomello26/booking-app-backend/internal/core/models"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) models.AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) RegisterUser(ctx context.Context, registerData *models.RegisterCredentials) (*models.User, error) {
	user := &models.User{
		FirstName:   registerData.FirstName,
		LastName:    registerData.LastName,
		PhoneNumber: registerData.PhoneNumber,
		Email:       registerData.Email,
		Password:    registerData.Password,
	}

	err := r.db.WithContext(ctx).
		Create(user).
		Error

	return user, err
}

func (r *AuthRepository) GetUser(ctx context.Context, query interface{}, args ...interface{}) (*models.User, error) {
	user := &models.User{}

	err := r.db.WithContext(ctx).
		Where(query, args...).
		First(user).
		Error

	return user, err
}
