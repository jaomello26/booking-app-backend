package handlers

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jaomello26/booking-app-backend/internal/core/models"
	"github.com/jaomello26/booking-app-backend/pkg/utils"
)

type ApartmentHandler struct {
	apartmentRepo models.ApartmentRepository
	userGroupRepo models.UserGroupRepository
}

func NewApartmentHandler(router fiber.Router, apartmentRepo models.ApartmentRepository, userGroupRepo models.UserGroupRepository) {
	handler := &ApartmentHandler{
		apartmentRepo: apartmentRepo,
		userGroupRepo: userGroupRepo,
	}

	router.Get("/", handler.GetMany)
	router.Get("/:apartmentId", handler.GetOne)
	router.Post("/", handler.CreateOne)
	router.Put("/:apartmentId", handler.UpdateOne)
	router.Delete("/:apartmentId", handler.DeleteOne)
}

func (h *ApartmentHandler) GetMany(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
		})
	}

	groupIDStr := c.Query("group_id")
	if groupIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Group ID is required",
		})
	}

	groupID, err := utils.StringToUint(groupIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid group ID",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = h.userGroupRepo.GetUserRole(ctx, userID, groupID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "You do not have access to this group",
		})
	}

	apartments, err := h.apartmentRepo.GetManyByGroup(ctx, groupID)
	if err != nil {
		log.Printf("Error retrieving apartments: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not retrieve apartments",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   apartments,
	})
}

func (h *ApartmentHandler) GetOne(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
		})
	}

	apartmentIDStr := c.Params("apartmentId")
	apartmentID, err := utils.StringToUint(apartmentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid apartment ID",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	apartment, err := h.apartmentRepo.GetOne(ctx, apartmentID)
	if err != nil {
		log.Printf("Error retrieving apartment: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Apartment not found",
		})
	}

	_, err = h.userGroupRepo.GetUserRole(ctx, userID, apartment.GroupID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "You do not have access to this apartment",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   apartment,
	})
}

func (h *ApartmentHandler) CreateOne(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
		})
	}

	apartment := new(models.Apartment)
	if err := c.BodyParser(apartment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request body",
		})
	}

	apartment.CreatedBy = userID
	apartment.ID = 0

	validate := utils.GetValidator()
	if err := validate.Struct(apartment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	role, err := h.userGroupRepo.GetUserRole(ctx, userID, apartment.GroupID)
	if err != nil || role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "You do not have permission to create apartments in this group",
		})
	}

	newApartment, err := h.apartmentRepo.CreateOne(ctx, apartment)
	if err != nil {
		log.Printf("Error creating apartment: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not create apartment",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Apartment created",
		"data":    newApartment,
	})
}

func (h *ApartmentHandler) UpdateOne(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
		})
	}

	apartmentIDStr := c.Params("apartmentId")
	apartmentID, err := utils.StringToUint(apartmentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid apartment ID",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	apartment, err := h.apartmentRepo.GetOne(ctx, apartmentID)
	if err != nil {
		log.Printf("Error retrieving apartment: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Apartment not found",
		})
	}

	role, err := h.userGroupRepo.GetUserRole(ctx, userID, apartment.GroupID)
	if err != nil || role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "You do not have permission to update apartments in this group",
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
	delete(updateData, "group_id")

	updatedApartment, err := h.apartmentRepo.UpdateOne(ctx, apartmentID, updateData)
	if err != nil {
		log.Printf("Error updating apartment: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not update apartment",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Apartment updated",
		"data":    updatedApartment,
	})
}

func (h *ApartmentHandler) DeleteOne(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
		})
	}

	apartmentIDStr := c.Params("apartmentId")
	apartmentID, err := utils.StringToUint(apartmentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid apartment ID",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	apartment, err := h.apartmentRepo.GetOne(ctx, apartmentID)
	if err != nil {
		log.Printf("Error retrieving apartment: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Apartment not found",
		})
	}

	role, err := h.userGroupRepo.GetUserRole(ctx, userID, apartment.GroupID)
	if err != nil || role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "You do not have permission to delete apartments in this group",
		})
	}

	err = h.apartmentRepo.DeleteOne(ctx, apartmentID)
	if err != nil {
		log.Printf("Error deleting apartment: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not delete apartment",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
