package handlers

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jaomello26/booking-app-backend/internal/core/models"
	"github.com/jaomello26/booking-app-backend/pkg/utils"
)

type BookingHandler struct {
	bookingRepo   models.BookingRepository
	apartmentRepo models.ApartmentRepository
	userGroupRepo models.UserGroupRepository
}

func NewBookingHandler(router fiber.Router, bookingRepo models.BookingRepository, apartmentRepo models.ApartmentRepository, userGroupRepo models.UserGroupRepository) {
	handler := &BookingHandler{
		bookingRepo:   bookingRepo,
		apartmentRepo: apartmentRepo,
		userGroupRepo: userGroupRepo,
	}

	router.Get("/", handler.GetMany)
	router.Post("/", handler.CreateOne)
	router.Get("/:bookingId", handler.GetOne)
	router.Put("/:bookingId", handler.UpdateOne)
	router.Delete("/:bookingId", handler.DeleteOne)
}

func (h *BookingHandler) GetMany(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
		})
	}

	apartmentIDQuery := c.Query("apartment_id", "")
	apartmentID, err := utils.StringToUint(apartmentIDQuery)
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
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Apartment not found",
		})
	}

	_, err = h.userGroupRepo.GetUserRole(ctx, userID, apartment.GroupID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "You do not have permission to view bookings for this apartment",
		})
	}

	bookings, err := h.bookingRepo.GetMany(ctx, apartmentID)
	if err != nil {
		log.Printf("Error retrieving bookings: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not retrieve bookings",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   bookings,
	})
}

func (h *BookingHandler) GetOne(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
		})
	}

	bookingIDParam := c.Params("bookingId")
	bookingID, err := utils.StringToUint(bookingIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid booking ID",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	booking, err := h.bookingRepo.GetOne(ctx, bookingID)
	if err != nil {
		log.Printf("Error retrieving booking: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Booking not found",
		})
	}

	apartment, err := h.apartmentRepo.GetOne(ctx, booking.ApartmentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not retrieve apartment",
		})
	}

	_, err = h.userGroupRepo.GetUserRole(ctx, userID, apartment.GroupID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "You do not have permission to view this booking",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   booking,
	})
}

func (h *BookingHandler) CreateOne(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
		})
	}

	booking := new(models.Booking)
	if err := c.BodyParser(booking); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request body",
		})
	}

	validate := utils.GetValidator()
	if err := validate.Struct(booking); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	booking.CreatedBy = userID

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	apartment, err := h.apartmentRepo.GetOne(ctx, booking.ApartmentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Apartment not found",
		})
	}

	role, err := h.userGroupRepo.GetUserRole(ctx, userID, apartment.GroupID)
	if err != nil || (role != "admin" && role != "editor") {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "You do not have permission to create bookings for this apartment",
		})
	}

	overlapping, err := h.bookingRepo.IsOverlapping(ctx, booking)
	if err != nil {
		log.Printf("Error checking booking overlap: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not create booking",
		})
	}
	if overlapping {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "fail",
			"message": "Booking dates overlap with an existing booking",
		})
	}

	newBooking, err := h.bookingRepo.CreateOne(ctx, booking)
	if err != nil {
		log.Printf("Error creating booking: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not create booking",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Booking created",
		"data":    newBooking,
	})
}

func (h *BookingHandler) UpdateOne(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
		})
	}

	bookingIDParam := c.Params("bookingId")
	bookingID, err := utils.StringToUint(bookingIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid booking ID",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	booking, err := h.bookingRepo.GetOne(ctx, bookingID)
	if err != nil {
		log.Printf("Error retrieving booking: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Booking not found",
		})
	}

	apartment, err := h.apartmentRepo.GetOne(ctx, booking.ApartmentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not retrieve apartment",
		})
	}

	role, err := h.userGroupRepo.GetUserRole(ctx, userID, apartment.GroupID)
	if err != nil || (role != "admin" && role != "editor") {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "You do not have permission to update bookings for this apartment",
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
	delete(updateData, "apartment_id")

	validate := utils.GetValidator()
	updatedBooking := &models.Booking{}
	if err := utils.MapToStruct(updateData, updatedBooking); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid update data",
		})
	}
	if err := validate.StructPartial(updatedBooking, "Name", "Phone", "CheckIn", "CheckOut", "NumberOfGuests", "Source", "AmountCharged"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	if checkIn, ok := updateData["check_in"]; ok {
		booking.CheckIn, _ = time.Parse(time.RFC3339, checkIn.(string))
	}
	if checkOut, ok := updateData["check_out"]; ok {
		booking.CheckOut, _ = time.Parse(time.RFC3339, checkOut.(string))
	}
	overlapping, err := h.bookingRepo.IsOverlapping(ctx, booking)
	if err != nil {
		log.Printf("Error checking booking overlap: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not update booking",
		})
	}
	if overlapping {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "fail",
			"message": "Booking dates overlap with an existing booking",
		})
	}

	updatedBooking, err = h.bookingRepo.UpdateOne(ctx, bookingID, updateData)
	if err != nil {
		log.Printf("Error updating booking: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not update booking",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Booking updated",
		"data":    updatedBooking,
	})
}

func (h *BookingHandler) DeleteOne(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
		})
	}

	bookingIDParam := c.Params("bookingId")
	bookingID, err := utils.StringToUint(bookingIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid booking ID",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	booking, err := h.bookingRepo.GetOne(ctx, bookingID)
	if err != nil {
		log.Printf("Error retrieving booking: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Booking not found",
		})
	}

	apartment, err := h.apartmentRepo.GetOne(ctx, booking.ApartmentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not retrieve apartment",
		})
	}

	role, err := h.userGroupRepo.GetUserRole(ctx, userID, apartment.GroupID)
	if err != nil || (role != "admin" && role != "editor") {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "You do not have permission to delete bookings for this apartment",
		})
	}

	err = h.bookingRepo.DeleteOne(ctx, bookingID)
	if err != nil {
		log.Printf("Error deleting booking: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Could not delete booking",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
