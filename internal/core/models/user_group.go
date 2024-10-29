package models

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type UserGroup struct {
	UserID   uint      `gorm:"primaryKey" json:"user_id"`
	GroupID  uint      `gorm:"primaryKey" json:"group_id"`
	Role     string    `gorm:"size:10;not null" json:"role" validate:"required,oneof=admin editor viewer"`
	JoinedAt time.Time `gorm:"autoCreateTime" json:"joined_at"`
}

type UserGroupRepository interface {
	AddMember(ctx context.Context, userGroup *UserGroup) (*UserGroup, error)
	AddMemberTx(ctx context.Context, tx *gorm.DB, group *UserGroup) error
	RemoveMember(ctx context.Context, userID uint, groupID uint) error
	UpdateMemberRole(ctx context.Context, userID uint, groupID uint, role string) error
	GetUserRole(ctx context.Context, userID uint, groupID uint) (string, error)
	GetGroupMembers(ctx context.Context, groupID uint) ([]*UserGroup, error)
}
