package repositories

import (
	"context"

	"github.com/jaomello26/booking-app-backend/internal/core/models"
	"gorm.io/gorm"
)

type ApartmentRepository struct {
	db *gorm.DB
}

func NewApartmentRepository(db *gorm.DB) models.ApartmentRepository {
	return &ApartmentRepository{
		db: db,
	}
}

func (r *ApartmentRepository) GetManyByGroup(ctx context.Context, groupID uint) ([]*models.Apartment, error) {
	apartments := []*models.Apartment{}

	err := r.db.WithContext(ctx).Model(&models.Apartment{}).
		Where("group_id = ?", groupID).
		Find(&apartments).
		Error

	return apartments, err
}

func (r *ApartmentRepository) GetOne(ctx context.Context, apartmentID uint) (*models.Apartment, error) {
	apartment := &models.Apartment{}

	err := r.db.WithContext(ctx).Model(apartment).
		Where("id = ?", apartmentID).
		First(apartment).
		Error

	return apartment, err
}

func (r *ApartmentRepository) CreateOne(ctx context.Context, apartment *models.Apartment) (*models.Apartment, error) {
	err := r.db.WithContext(ctx).Model(&models.Apartment{}).
		Create(apartment).
		Error

	return apartment, err
}

func (r *ApartmentRepository) UpdateOne(ctx context.Context, apartmentID uint, updateData map[string]interface{}) (*models.Apartment, error) {
	apartment := &models.Apartment{}

	err := r.db.WithContext(ctx).Model(apartment).
		Where("id = ?", apartmentID).
		Updates(updateData).Error

	if err != nil {
		return nil, err
	}

	err = r.db.WithContext(ctx).Model(apartment).
		Where("id = ?", apartmentID).
		First(apartment, apartmentID).
		Error

	return apartment, err
}

func (r *ApartmentRepository) DeleteOne(ctx context.Context, apartmentID uint) error {
	return r.db.WithContext(ctx).
		Delete(&models.Apartment{}, apartmentID).
		Error
}
