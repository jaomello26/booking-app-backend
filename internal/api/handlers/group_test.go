//go:build test

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jaomello26/booking-app-backend/internal/core/models"
	"github.com/jaomello26/booking-app-backend/internal/core/models/mocks"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func authMiddlewareMock(userID uint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Simula um usuário autenticado definindo o userID no contexto
		c.Locals("userId", userID)
		return c.Next()
	}
}

func TestGroupHandler_CreateOne(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGroupRepo := mocks.NewMockGroupRepository(ctrl)
	mockUserGroupRepo := mocks.NewMockUserGroupRepository(ctrl)
	db := &gorm.DB{} // Pode ser nil, pois não será usado diretamente

	handler := &GroupHandler{
		GroupRepo:     mockGroupRepo,
		userGroupRepo: mockUserGroupRepo,
		db:            db,
	}

	app := fiber.New()
	userID := uint(1)

	// Adicionar o middleware de autenticação mock
	app.Use(authMiddlewareMock(userID))

	app.Post("/groups", handler.CreateOne)

	groupData := models.Group{
		Name: "Family Group",
	}

	requestBody, _ := json.Marshal(groupData)
	req := httptest.NewRequest("POST", "/groups", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	// Não é necessário definir o header Authorization, pois o middleware mock já cuida disso

	// Configurar os mocks
	mockGroupRepo.EXPECT().
		CreateOneTx(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, tx *gorm.DB, group *models.Group) (*models.Group, error) {
			group.ID = 1
			return group, nil
		})

	mockUserGroupRepo.EXPECT().
		AddMemberTx(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var responseBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseBody)

	if responseBody["status"] != "success" {
		t.Errorf("Expected status 'success', got '%v'", responseBody["status"])
	}
}
