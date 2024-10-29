package utils

import (
	"context"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func GetUserID(ctx *fiber.Ctx) (uint, error) {
	userIDInterface := ctx.Locals("userId")

	if userIDInterface == nil {
		return 0, errors.New("user ID not found in context")
	}

	userIDFloat, ok := userIDInterface.(float64)
	if !ok {
		return 0, errors.New("user ID has invalid type")
	}

	userID := uint(userIDFloat)

	return userID, nil
}

// only for tests
func SetUserID(req *http.Request, userID uint) {
	ctx := req.Context()
	ctx = context.WithValue(ctx, "userId", userID)
	*req = *req.WithContext(ctx)
}
