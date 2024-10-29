package models

import (
	"context"
	"time"
)

type Apartment struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name" validate:"required"`
	Description string    `gorm:"type:text" json:"description"`
	GroupID     uint      `gorm:"not null" json:"group_id" validate:"required"`
	CreatedBy   uint      `gorm:"not null" json:"created_by"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type ApartmentRepository interface {
	GetManyByGroup(ctx context.Context, groupID uint) ([]*Apartment, error)
	GetOne(ctx context.Context, apartmentID uint) (*Apartment, error)
	CreateOne(ctx context.Context, apartment *Apartment) (*Apartment, error)
	UpdateOne(ctx context.Context, apartmentID uint, updateData map[string]interface{}) (*Apartment, error)
	DeleteOne(ctx context.Context, apartmentID uint) error
}
