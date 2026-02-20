package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"github.com/FANIMAN/housing-lottery/internal/config"
	"github.com/FANIMAN/housing-lottery/internal/delivery/http"
	"github.com/FANIMAN/housing-lottery/internal/delivery/middleware"
	"github.com/FANIMAN/housing-lottery/internal/infrastructure/persistence"
	"github.com/FANIMAN/housing-lottery/internal/usecase"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db := config.NewDB()

	// Repositories
	adminRepo := persistence.NewAdminRepository(db)
	auditRepo := persistence.NewAuditRepository(db)
	subcityRepo := persistence.NewSubcityRepository(db) // exported type SubcityRepo
	applicantRepo := persistence.NewApplicantRepository(db)
	uploadBatchRepo := persistence.NewUploadBatchRepository(db)
	lotteryRepo := persistence.NewLotteryRepository(db)
	lotteryWinnerRepo := persistence.NewLotteryWinnerRepository(db)

	// Usecases
	adminUsecase := usecase.NewAdminUsecase(adminRepo, auditRepo, os.Getenv("JWT_SECRET"))
	subcityUsecase := usecase.NewSubcityUsecase(subcityRepo, auditRepo)
	uploadService := usecase.NewUploadService(applicantRepo, uploadBatchRepo, auditRepo)
	lotteryService := usecase.NewLotteryService(lotteryRepo, applicantRepo, lotteryWinnerRepo, auditRepo)

	// Handlers
	adminHandler := http.NewAdminHandler(adminUsecase)
	subcityHandler := http.NewSubcityHandler(subcityUsecase)
	uploadHandler := http.NewUploadHandler(uploadService)
	lotteryHandler := http.NewLotteryHandler(lotteryService)

	// Dashboard
	dashboardRepo := persistence.NewDashboardRepository(db)
	dashboardUsecase := usecase.NewDashboardUsecase(dashboardRepo, subcityRepo, lotteryRepo)
	dashboardHandler := http.NewDashboardHandler(dashboardUsecase)

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))
	app.Use(middleware.AuditMiddleware(auditRepo))

	// Public routes
	app.Post("/admin/login", adminHandler.Login)

	// Protected routes
	api := app.Group("/api", middleware.JWTMiddleware())

	// Admin
	api.Post("/admin/register", adminHandler.Register)

	// Dashboard endpoints
	api.Get("/dashboard/summary", dashboardHandler.GetSummary)
	api.Get("/subcities", dashboardHandler.ListSubcities)
	api.Get("/lotteries", dashboardHandler.ListLotteries)

	// Subcity CRUD
	api.Post("/subcities", subcityHandler.Create)
	api.Get("/subcities/list", subcityHandler.List) // keep previous endpoint
	api.Put("/subcities/:id", subcityHandler.Update)
	api.Delete("/subcities/:id", subcityHandler.Delete)

	// Upload
	api.Post("/subcities/:id/upload", uploadHandler.UploadApplicants)
	api.Get("/applicants", uploadHandler.ListApplicants)

	// Lottery
	api.Post("/subcities/:id/lottery/start", lotteryHandler.Start)
	api.Post("/lotteries/:id/spin", lotteryHandler.Spin)
	api.Post("/lotteries/:id/close", lotteryHandler.Close)

	log.Println("Server running on :8080")
	log.Fatal(app.Listen(":8080"))
}
