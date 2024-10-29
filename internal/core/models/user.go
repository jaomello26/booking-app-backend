package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name" validate:"required,min=3,max=255"`
	Email     string    `gorm:"size:255;unique;not null" json:"email" validate:"required,email,max=255"`
	Password  string    `gorm:"size:255;not null" json:"-" validate:"required,min=4,max=255"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
