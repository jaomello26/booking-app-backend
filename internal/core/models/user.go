package models

import (
	"time"
)

type User struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName   string    `gorm:"size:255;not null" json:"first_name" validate:"required,min=2,max=255"`
	LastName    string    `gorm:"size:255;not null" json:"last_name" validate:"required,min=2,max=255"`
	PhoneNumber string    `gorm:"size:20;not null" json:"phone_number" validate:"required"`
	Email       string    `gorm:"size:255;unique;not null" json:"email" validate:"required,email,max=255"`
	Password    string    `gorm:"size:255;not null" json:"-" validate:"required,min=6,max=255"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
