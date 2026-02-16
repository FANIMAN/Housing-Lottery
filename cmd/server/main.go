package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/FANIMAN/housing-lottery/internal/config"
	"github.com/FANIMAN/housing-lottery/internal/delivery/http"
	"github.com/FANIMAN/housing-lottery/internal/infrastructure/persistence"
	"github.com/FANIMAN/housing-lottery/internal/usecase"
)

func main() {

	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// Connect DB
	db := config.NewDB()

	// Repository
	adminRepo := persistence.NewAdminRepository(db)

	// Usecase
	adminUsecase := usecase.NewAdminUsecase(adminRepo)

	// Handler
	adminHandler := http.NewAdminHandler(adminUsecase)

	app := fiber.New()

	// Routes
	app.Post("/admin/register", adminHandler.Register)

	log.Println("Server running on :8080")
	app.Listen(":8080")
}
