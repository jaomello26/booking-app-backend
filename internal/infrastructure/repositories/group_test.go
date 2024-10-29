package repositories

import (
	"context"
	"testing"

	"github.com/jaomello26/booking-app-backend/internal/core/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestGroupRepository(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&models.Group{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestGroupRepository_CreateOne(t *testing.T) {
	db := setupTestGroupRepository(t)
	repo := NewGroupRepository(db)
	ctx := context.Background()

	group := &models.Group{
		Name:      "Family Group",
		CreatedBy: 1,
	}

	createdGroup, err := repo.CreateOne(ctx, group)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if createdGroup.ID == 0 {
		t.Errorf("Expected group ID to be set")
	}

	if createdGroup.Name != group.Name {
		t.Errorf("Expected name %s, got %s", group.Name, createdGroup.Name)
	}
}

func TestGroupRepository_GetOne(t *testing.T) {
	db := setupTestGroupRepository(t)
	repo := NewGroupRepository(db)
	ctx := context.Background()

	group := &models.Group{
		Name:      "Family Group",
		CreatedBy: 1,
	}
	repo.CreateOne(ctx, group)

	fetchedGroup, err := repo.GetOne(ctx, group.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if fetchedGroup.ID != group.ID {
		t.Errorf("Expected group ID %d, got %d", group.ID, fetchedGroup.ID)
	}
}

func TestGroupRepository_UpdateOne(t *testing.T) {
	db := setupTestGroupRepository(t)
	repo := NewGroupRepository(db)
	ctx := context.Background()

	group := &models.Group{
		Name:      "Family Group",
		CreatedBy: 1,
	}
	repo.CreateOne(ctx, group)

	updateData := map[string]interface{}{
		"Name": "Friends Group",
	}
	updatedGroup, err := repo.UpdateOne(ctx, group.ID, updateData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updatedGroup.Name != "Friends Group" {
		t.Errorf("Expected name 'Friends Group', got '%s'", updatedGroup.Name)
	}
}

func TestGroupRepository_DeleteOne(t *testing.T) {
	db := setupTestGroupRepository(t)
	repo := NewGroupRepository(db)
	ctx := context.Background()

	group := &models.Group{
		Name:      "Family Group",
		CreatedBy: 1,
	}
	repo.CreateOne(ctx, group)

	err := repo.DeleteOne(ctx, group.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = repo.GetOne(ctx, group.ID)
	if err == nil {
		t.Fatalf("Expected error, got none")
	}
	if err != gorm.ErrRecordNotFound {
		t.Fatalf("Expected ErrRecordNotFound, got %v", err)
	}
}
