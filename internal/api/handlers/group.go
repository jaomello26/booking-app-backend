package handlers

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jaomello26/booking-app-backend/internal/core/models"
	"github.com/jaomello26/booking-app-backend/pkg/utils"
	"gorm.io/gorm"
)

type GroupHandler struct {
	GroupRepo     models.GroupRepository
	userGroupRepo models.UserGroupRepository
	db            *gorm.DB
}

func NewGroupHandler(router fiber.Router, GroupRepo models.GroupRepository, UserGroupRepo models.UserGroupRepository, db *gorm.DB) {
	handler := &GroupHandler{
		GroupRepo:     GroupRepo,
		userGroupRepo: UserGroupRepo,
		db:            db,
	}

	router.Get("/", handler.GetMany)
	router.Post("/", handler.CreateOne)
	router.Put("/:groupId", handler.UpdateOne)
	router.Delete("/:groupId", handler.DeleteOne)
}

func (h *GroupHandler) GetMany(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	groups, err := h.GroupRepo.GetMany(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  "success",
		"message": "",
		"data":    groups,
	})
}

func (h *GroupHandler) CreateOne(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
		})
	}

	group := &models.Group{}
	if err := c.BodyParser(group); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
			"data":    nil,
		})
	}

	validate := utils.GetValidator()
	if err := validate.Struct(group); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	group.CreatedBy = userID

	tx := h.db.Begin()
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not start transaction",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	group, err = h.GroupRepo.CreateOneTx(ctx, tx, group)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
			"data":    nil,
		})
	}

	userGroup := &models.UserGroup{
		UserID:  userID,
		GroupID: group.ID,
		Role:    "admin",
	}

	err = h.userGroupRepo.AddMemberTx(ctx, tx, userGroup)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not add user to group",
		})
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not commit transaction",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Group created",
		"data":    group,
	})
}

func (h *GroupHandler) UpdateOne(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
		})
	}

	groupIDParam := c.Params("groupId")
	groupID, err := utils.StringToUint(groupIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	group, err := h.GroupRepo.GetOne(ctx, groupID)

	if err != nil {
		log.Printf("Error fetching group: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Group not found",
		})
	}

	if group.CreatedBy != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "You do not have permission to update this group",
		})
	}

	updateData := make(map[string]interface{})
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request body",
		})
	}

	delete(updateData, "id")
	delete(updateData, "created_by")
	delete(updateData, "created_at")
	delete(updateData, "updated_at")

	group, err = h.GroupRepo.UpdateOne(ctx, groupID, updateData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"status":  "success",
		"message": "Group updated",
		"data":    group,
	})
}

func (h *GroupHandler) DeleteOne(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
		})
	}

	groupIDParam := c.Params("groupId")
	groupID, err := utils.StringToUint(groupIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	group, err := h.GroupRepo.GetOne(ctx, groupID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Group not found",
		})
	}

	if group.CreatedBy != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "You do not have permission to delete this group",
		})
	}

	err = h.GroupRepo.DeleteOne(ctx, groupID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
