package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jaomello26/booking-app-backend/config"
	"github.com/jaomello26/booking-app-backend/internal/api/handlers"
	"github.com/jaomello26/booking-app-backend/internal/api/middlewares"
	"github.com/jaomello26/booking-app-backend/internal/core/services"
	"github.com/jaomello26/booking-app-backend/internal/infrastructure/db"
	"github.com/jaomello26/booking-app-backend/internal/infrastructure/repositories"
)

func main() {
	envConfig := config.NewEnvConfig()
	database := db.Init(envConfig, db.DBMigrator)

	app := fiber.New(fiber.Config{
		AppName:      "Booking APP",
		ServerHeader: "Fiber",
	})

	// Repositories
	authRepository := repositories.NewAuthRepository(database)
	groupRepository := repositories.NewGroupRepository(database)
	userGroupRepository := repositories.NewUserGroupRepository(database)
	apartmentRepository := repositories.NewApartmentRepository(database)
	bookingRepository := repositories.NewBookingRepository(database)

	// Services
	authService := services.NewAuthService(authRepository)

	// Routing
	server := app.Group("/api")
	handlers.NewAuthHandler(server.Group("/auth"), authService)

	privateRoutes := server.Use(middlewares.AuthProtected(database))

	// Handlers
	handlers.NewGroupHandler(privateRoutes.Group("/group"), groupRepository, userGroupRepository, database)
	handlers.NewUserGroupHandler(privateRoutes.Group("/user-group"), userGroupRepository, groupRepository)
	handlers.NewApartmentHandler(privateRoutes.Group("/apartment"), apartmentRepository, userGroupRepository)
	handlers.NewBookingHandler(privateRoutes.Group("/booking"), bookingRepository, apartmentRepository, userGroupRepository)

	if err := app.Listen(":" + envConfig.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
