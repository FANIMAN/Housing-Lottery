package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/FANIMAN/housing-lottery/internal/config"
	"github.com/FANIMAN/housing-lottery/internal/delivery/http"
	"github.com/FANIMAN/housing-lottery/internal/infrastructure/persistence"
	"github.com/FANIMAN/housing-lottery/internal/usecase"
	"github.com/FANIMAN/housing-lottery/internal/delivery/middleware"

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
	adminUsecase := usecase.NewAdminUsecase(adminRepo, os.Getenv("JWT_SECRET"))

	// Handler
	adminHandler := http.NewAdminHandler(adminUsecase)

	app := fiber.New()

	// Routes
	app.Post("/admin/register", adminHandler.Register)
	app.Post("/admin/login", adminHandler.Login)

	// Protected routes group
	api := app.Group("/api", middleware.JWTMiddleware())
	
	// Subcity
	subcityRepo := persistence.NewSubcityRepository(db)
	subcityUsecase := usecase.NewSubcityUsecase(subcityRepo)
	subcityHandler := http.NewSubcityHandler(subcityUsecase)

	api.Post("/subcities", subcityHandler.Create)
	api.Get("/subcities", subcityHandler.List)
	api.Put("/subcities/:id", subcityHandler.Update)
	api.Delete("/subcities/:id", subcityHandler.Delete)



	uploadBatchRepo := persistence.NewUploadBatchRepository(db)
	uploadService := usecase.NewUploadService(
		persistence.NewApplicantRepository(db),
		uploadBatchRepo,
	)
	uploadHandler := http.NewUploadHandler(uploadService)

	api.Post("/subcities/:id/upload", uploadHandler.UploadApplicants)


	log.Println("Server running on :8080")
	app.Listen(":8080")
}
