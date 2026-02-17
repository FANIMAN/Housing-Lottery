package middleware

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/FANIMAN/housing-lottery/internal/repository/interfaces"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Check if Authorization header exists
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Println("Authorization header missing")
			return fiber.ErrUnauthorized
		}

		// Expecting: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Println("Authorization header malformed:", authHeader)
			return fiber.ErrUnauthorized
		}

		tokenString := parts[1]

		// Parse the token using the secret
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		// Log token parsing errors
		if err != nil {
			log.Println("JWT parsing error:", err)
			return fiber.ErrUnauthorized
		}

		// Check if the token is valid
		if !token.Valid {
			log.Println("Invalid token")
			return fiber.ErrUnauthorized
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Failed to extract claims from token")
			return fiber.ErrUnauthorized
		}

		// Log claims to verify admin_id
		log.Println("JWT claims:", claims)

		// Store admin_id in context
		c.Locals("admin_id", claims["admin_id"])

		return c.Next()
	}
}



func AuditMiddleware(auditRepo interfaces.AuditRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Run the next handler
		err := c.Next()

		// Capture response details
		status := c.Response().StatusCode()
		ip := c.IP()
		userAgent := c.Get("User-Agent")
		duration := time.Since(start)

		var errorMessage string
		if err != nil {
			errorMessage = err.Error()
		}

		// Get admin_id if exists (public routes may not have it)
		adminID := ""
		if id, ok := c.Locals("admin_id").(string); ok {
			adminID = id
		}

		// Action = HTTP method + path
		action := c.Method() + " " + c.Path()

		// Include duration in error_message if needed
		if duration > 0 {
			if errorMessage != "" {
				errorMessage += " | "
			}
			errorMessage += "duration=" + duration.String()
		}

		// Log the HTTP request in audit_logs
		_ = auditRepo.Log(
			c.Context(),
			adminID,
			action,
			"http_request",
			"",       // entity_id empty for generic requests
			status,
			ip,
			userAgent,
			errorMessage,
		)

		return err
	}
}
