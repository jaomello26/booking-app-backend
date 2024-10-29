//go:build test

package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jaomello26/booking-app-backend/internal/core/models"
	"github.com/jaomello26/booking-app-backend/internal/core/models/mocks"
	"go.uber.org/mock/gomock"
)

func TestAuthHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockAuthService(ctrl)
	handler := &AuthHandler{
		service: mockService,
	}

	app := fiber.New()
	app.Post("/register", handler.Register)

	registerData := models.AuthCredentials{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	requestBody, _ := json.Marshal(registerData)
	req := httptest.NewRequest("POST", "/register", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	token := "testtoken"
	user := &models.User{
		ID:    1,
		Name:  registerData.Name,
		Email: registerData.Email,
	}

	mockService.EXPECT().
		Register(gomock.Any(), &registerData).
		Return(token, user, nil)

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

func TestAuthHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockAuthService(ctrl)
	handler := &AuthHandler{
		service: mockService,
	}

	app := fiber.New()
	app.Post("/login", handler.Login)

	loginData := models.AuthCredentials{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	requestBody, _ := json.Marshal(loginData)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	token := "testtoken"
	user := &models.User{
		ID:    1,
		Name:  "John Doe",
		Email: loginData.Email,
	}

	mockService.EXPECT().
		Login(gomock.Any(), &loginData).
		Return(token, user, nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var responseBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseBody)

	if responseBody["status"] != "success" {
		t.Errorf("Expected status 'success', got '%v'", responseBody["status"])
	}
}
