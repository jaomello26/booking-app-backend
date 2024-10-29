package repositories

import (
	"context"

	"github.com/jaomello26/booking-app-backend/internal/core/models"

	"gorm.io/gorm"
)

type GroupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) models.GroupRepository {
	return &GroupRepository{
		db: db,
	}
}

func (r *GroupRepository) GetMany(ctx context.Context, userID uint) ([]*models.Group, error) {
	groups := []*models.Group{}

	err := r.db.WithContext(ctx).Model(&models.Group{}).
		Where("created_by = ?", userID).
		Find(&groups).
		Error

	return groups, err
}

func (r *GroupRepository) GetOne(ctx context.Context, groupID uint) (*models.Group, error) {
	group := &models.Group{}

	err := r.db.WithContext(ctx).Model(group).
		Where("id = ?", groupID).
		First(group).
		Error

	return group, err
}

func (r *GroupRepository) CreateOne(ctx context.Context, group *models.Group) (*models.Group, error) {
	err := r.db.WithContext(ctx).Model(group).
		Create(group).
		Error

	return group, err
}

func (r *GroupRepository) CreateOneTx(ctx context.Context, tx *gorm.DB, group *models.Group) (*models.Group, error) {
	err := tx.WithContext(ctx).Model(group).
		Create(group).
		Error

	return group, err
}

func (r *GroupRepository) UpdateOne(ctx context.Context, groupID uint, updateData map[string]interface{}) (*models.Group, error) {
	group := &models.Group{}

	err := r.db.WithContext(ctx).Model(group).
		Where("id = ?", groupID).
		Updates(updateData).
		Error
	if err != nil {
		return nil, err
	}

	err = r.db.Model(group).Where("id = ?", groupID).
		First(group).
		Error

	return group, err
}

func (r *GroupRepository) DeleteOne(ctx context.Context, groupId uint) error {
	err := r.db.WithContext(ctx).
		Delete(&models.Group{}, groupId).
		Error

	return err
}
