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

func (r *AuthRepository) RegisterUser(ctx context.Context, registerData *models.AuthCredentials) (*models.User, error) {
	user := &models.User{
		Name:     registerData.Name,
		Email:    registerData.Email,
		Password: registerData.Password,
	}

	err := r.db.WithContext(ctx).Model(&models.User{}).
		Create(user).
		Error

	return user, err
}

func (r *AuthRepository) GetUser(ctx context.Context, query interface{}, args ...interface{}) (*models.User, error) {
	user := &models.User{}

	err := r.db.WithContext(ctx).Model(user).
		Where(query, args...).
		First(user).
		Error

	return user, err
}
