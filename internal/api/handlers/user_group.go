// internal/api/handlers/user_group.go
package handlers

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jaomello26/booking-app-backend/internal/core/models"
	"github.com/jaomello26/booking-app-backend/pkg/utils"
)

type UserGroupHandler struct {
	userGroupRepo models.UserGroupRepository
	groupRepo     models.GroupRepository
}

func NewUserGroupHandler(router fiber.Router, userGroupRepo models.UserGroupRepository, groupRepo models.GroupRepository) {
	handler := &UserGroupHandler{
		userGroupRepo: userGroupRepo,
		groupRepo:     groupRepo,
	}

	router.Post("/", handler.AddMember)
	router.Delete("/", handler.RemoveMember)
	router.Put("/role", handler.UpdateMemberRole)
	router.Get("/role", handler.GetUserRole)
}

type userGroupRequest struct {
	UserID  uint   `json:"user_id" validate:"required"`
	GroupID uint   `json:"group_id" validate:"required"`
	Role    string `json:"role" validate:"required,oneof=editor viewer"`
}

func (h *UserGroupHandler) AddMember(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
		})
	}

	request := userGroupRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request body",
		})
	}

	validate := utils.GetValidator()
	if err := validate.Struct(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	role, err := h.userGroupRepo.GetUserRole(ctx, userID, request.GroupID)
	if err != nil || role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "You do not have permission to add members to this group",
		})
	}

	_, err = h.userGroupRepo.GetUserRole(ctx, request.UserID, request.GroupID)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "User is already a member of the group",
		})
	}

	userGroup := &models.UserGroup{
		UserID:  request.UserID,
		GroupID: request.GroupID,
		Role:    request.Role,
	}

	_, err = h.userGroupRepo.AddMember(ctx, userGroup)
	if err != nil {
		log.Printf("Error adding member: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not add member to the group",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Member added to the group",
	})
}

func (h *UserGroupHandler) RemoveMember(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
		})
	}

	var request struct {
		UserID  uint `json:"user_id" validate:"required"`
		GroupID uint `json:"group_id" validate:"required"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request body",
		})
	}

	validate := utils.GetValidator()
	if err := validate.Struct(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	role, err := h.userGroupRepo.GetUserRole(ctx, userID, request.GroupID)
	if err != nil || role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "You do not have permission to remove members from this group",
		})
	}

	if request.UserID == userID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Admin cannot remove themselves from the group",
		})
	}

	err = h.userGroupRepo.RemoveMember(ctx, request.UserID, request.GroupID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not remove member from the group",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Member removed from the group",
	})
}

func (h *UserGroupHandler) UpdateMemberRole(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
		})
	}

	request := userGroupRequest{}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request body",
		})
	}

	validate := utils.GetValidator()
	if err := validate.Struct(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	role, err := h.userGroupRepo.GetUserRole(ctx, userID, request.GroupID)
	if err != nil || role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "You do not have permission to update member roles in this group",
		})
	}

	if request.UserID == userID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Admin cannot change their own role",
		})
	}

	err = h.userGroupRepo.UpdateMemberRole(ctx, request.UserID, request.GroupID, request.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not update member role",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Member role updated",
	})
}

func (h *UserGroupHandler) GetUserRole(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
		})
	}

	groupIdParam := c.Query("group_id")
	if groupIdParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Group ID is required",
		})
	}

	groupID, err := utils.StringToUint(groupIdParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid group ID",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	role, err := h.userGroupRepo.GetUserRole(ctx, userID, groupID)
	if err != nil {
		log.Printf("Error retrieving user role: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Role not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"role": role,
		},
	})
}
