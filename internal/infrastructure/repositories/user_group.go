package repositories

import (
	"context"

	"github.com/jaomello26/booking-app-backend/internal/core/models"

	"gorm.io/gorm"
)

type UserGroupRepository struct {
	db *gorm.DB
}

func NewUserGroupRepository(db *gorm.DB) models.UserGroupRepository {
	return &UserGroupRepository{
		db: db,
	}
}

func (r *UserGroupRepository) AddMember(ctx context.Context, userGroup *models.UserGroup) (*models.UserGroup, error) {
	err := r.db.WithContext(ctx).Model(&models.UserGroup{}).
		Create(userGroup).
		Error

	return userGroup, err
}

func (r *UserGroupRepository) AddMemberTx(ctx context.Context, tx *gorm.DB, userGroup *models.UserGroup) error {
	return tx.WithContext(ctx).Model(userGroup).
		Create(userGroup).
		Error
}

func (r *UserGroupRepository) RemoveMember(ctx context.Context, userID uint, groupID uint) error {
	return r.db.WithContext(ctx).Model(&models.UserGroup{}).
		Delete(&models.UserGroup{}, "user_id = ? AND group_id = ?", userID, groupID).
		Error
}

func (r *UserGroupRepository) UpdateMemberRole(ctx context.Context, userID uint, groupID uint, role string) error {
	return r.db.WithContext(ctx).Model(&models.UserGroup{}).
		Where("user_id = ? AND group_id = ?", userID, groupID).
		Update("role", role).
		Error
}

func (r *UserGroupRepository) GetUserRole(ctx context.Context, userID uint, groupID uint) (string, error) {
	userGroup := &models.UserGroup{}

	err := r.db.WithContext(ctx).Model(userGroup).
		First(&userGroup, "user_id = ? AND group_id = ?", userID, groupID).
		Error

	return userGroup.Role, err
}

func (r *UserGroupRepository) GetGroupMembers(ctx context.Context, groupID uint) ([]*models.UserGroup, error) {
	userGroups := []*models.UserGroup{}

	err := r.db.WithContext(ctx).Model(&models.UserGroup{}).
		Where("group_id = ?", groupID).
		Find(&userGroups).
		Error

	return userGroups, err
}
