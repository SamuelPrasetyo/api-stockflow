package middleware

import (
	"api-stockflow/internal/domain"
	"api-stockflow/internal/repository"
	"api-stockflow/internal/security"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware creates authentication middleware
func AuthMiddleware(tokenManager security.TokenManager, userRepo repository.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "Missing authorization header",
			})
		}

		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "Invalid authorization header format",
			})
		}

		// Verify token
		token := parts[1]
		payload, err := tokenManager.VerifyAccessToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "Invalid or expired token",
			})
		}

		// Get user from database
		user, err := userRepo.GetByID(c.Context(), payload.UserID)
		if err != nil {
			if err == domain.ErrUserNotFound {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"status":  "fail",
					"message": "User not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Internal server error",
			})
		}

		// Set user to context
		c.Locals("user", user)

		return c.Next()
	}
}

// RoleMiddleware creates role-based authorization middleware
func RoleMiddleware(allowedRoles ...domain.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*domain.User)

		// Check if user role is allowed
		for _, role := range allowedRoles {
			if user.Role == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "Insufficient permissions",
		})
	}
}

// ApprovalMiddleware checks if user can approve based on amount
func ApprovalMiddleware(amountField string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*domain.User)

		// Parse request body to get amount
		var body map[string]interface{}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": "Invalid request body",
			})
		}

		// Get amount from body
		amount, ok := body[amountField].(float64)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": "Invalid amount field",
			})
		}

		// Check if user can approve
		if !user.CanApprove(amount) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "fail",
				"message": "You cannot approve purchases over 10,000,000. Please contact Admin/Manager",
			})
		}

		return c.Next()
	}
}
