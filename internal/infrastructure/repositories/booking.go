package repositories

import (
	"context"

	"github.com/jaomello26/booking-app-backend/internal/core/models"
	"gorm.io/gorm"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) models.BookingRepository {
	return &BookingRepository{
		db: db,
	}
}

func (r *BookingRepository) GetMany(ctx context.Context, apartmentID uint) ([]*models.Booking, error) {
	bookings := []*models.Booking{}

	err := r.db.WithContext(ctx).Model(&models.Booking{}).
		Where("apartment_id = ?", apartmentID).
		Order("check_in asc").
		Find(&bookings).
		Error

	return bookings, err
}

func (r *BookingRepository) GetOne(ctx context.Context, bookingID uint) (*models.Booking, error) {
	booking := &models.Booking{}

	err := r.db.WithContext(ctx).Model(&models.Booking{}).
		Where("id = ?", bookingID).
		First(booking).
		Error

	return booking, err
}

func (r *BookingRepository) CreateOne(ctx context.Context, booking *models.Booking) (*models.Booking, error) {
	err := r.db.WithContext(ctx).Model(booking).
		Create(booking).
		Error

	return booking, err
}

func (r *BookingRepository) UpdateOne(ctx context.Context, bookingID uint, updateData map[string]interface{}) (*models.Booking, error) {
	booking := &models.Booking{}

	err := r.db.WithContext(ctx).Model(booking).
		Where("id = ?", bookingID).
		Updates(updateData).
		Error

	if err != nil {
		return nil, err
	}

	err = r.db.WithContext(ctx).Model(booking).
		Where("id = ?", bookingID).
		First(booking).
		Error

	return booking, err
}

func (r *BookingRepository) DeleteOne(ctx context.Context, bookingID uint) error {
	err := r.db.WithContext(ctx).
		Delete(&models.Booking{}, bookingID).
		Error

	return err
}

func (r *BookingRepository) IsOverlapping(ctx context.Context, booking *models.Booking) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&models.Booking{}).
		Where("apartment_id = ?", booking.ApartmentID).
		Where("id != ?", booking.ID).
		Where("check_in < ? AND check_out > ?", booking.CheckOut, booking.CheckIn).
		Count(&count).
		Error

	return count > 0, err
}
