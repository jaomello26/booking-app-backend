package models

import (
	"context"
	"time"
)

type Booking struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string    `gorm:"size:255;not null" json:"name" validate:"required"`
	Phone           *string   `gorm:"size:20" json:"phone" validate:"omitempty,max=20"`
	CheckIn         time.Time `gorm:"not null" json:"check_in" validate:"required"`
	CheckOut        time.Time `gorm:"not null" json:"check_out" validate:"required,gtfield=CheckIn"`
	NumberOfGuests  int       `gorm:"not null" json:"number_of_guests"`
	CreatedBy       uint      `gorm:"not null" json:"created_by"`
	Source          *string   `gorm:"size:50" json:"source" validate:"omitempty,max=50"`
	AmountCharged   *float64  `gorm:"type:decimal(10,2)" json:"amount_charged"`
	ReservationDate time.Time `gorm:"autoCreateTime" json:"reservation_date"`
	ApartmentID     uint      `gorm:"not null" json:"apartment_id" validate:"required"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type BookingRepository interface {
	GetMany(ctx context.Context, apartmentID uint) ([]*Booking, error)
	GetOne(ctx context.Context, bookingID uint) (*Booking, error)
	CreateOne(ctx context.Context, booking *Booking) (*Booking, error)
	UpdateOne(ctx context.Context, bookingID uint, updateData map[string]interface{}) (*Booking, error)
	DeleteOne(ctx context.Context, bookingID uint) error
	IsOverlapping(ctx context.Context, booking *Booking) (bool, error)
}
