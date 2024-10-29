package models

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Group struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name" validate:"required,min=1,max=255"`
	CreatedBy uint      `gorm:"not null" json:"created_by" validate:"required"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type GroupRepository interface {
	GetMany(ctx context.Context, userId uint) ([]*Group, error)
	GetOne(ctx context.Context, groupId uint) (*Group, error)
	CreateOne(ctx context.Context, group *Group) (*Group, error)
	CreateOneTx(ctx context.Context, tx *gorm.DB, group *Group) (*Group, error)
	UpdateOne(ctx context.Context, groupId uint, updateData map[string]interface{}) (*Group, error)
	DeleteOne(ctx context.Context, groupId uint) error
}
