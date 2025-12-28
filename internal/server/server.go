package server

import (
	"github.com/gofiber/fiber/v2"

	"api-stockflow/internal/database"
	"api-stockflow/internal/handler"
	"api-stockflow/internal/middleware"
	"api-stockflow/internal/repository"
	"api-stockflow/internal/security"
	"api-stockflow/internal/usecase"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "api-stockflow",
			AppName:      "api-stockflow",
		}),

		db: database.New(),
	}

	return server
}

// SetupAuthRoutes sets up authentication routes
func (s *FiberServer) SetupAuthRoutes() {
	// Initialize dependencies
	tokenManager := security.NewJWTTokenManager()
	passwordHash := security.NewBcryptPasswordHash()
	
	userRepo := repository.NewUserRepository(s.db.GetDB())
	authRepo := repository.NewAuthenticationRepository(s.db.GetDB())
	
	authUseCase := usecase.NewAuthUseCase(userRepo, authRepo, tokenManager, passwordHash)
	authHandler := handler.NewAuthHandler(authUseCase)

	// Public routes
	auth := s.App.Group("/api/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Put("/refresh", authHandler.RefreshToken)
	auth.Delete("/logout", authHandler.Logout)

	// Protected routes
	authMiddleware := middleware.AuthMiddleware(tokenManager, userRepo)
	auth.Get("/profile", authMiddleware, authHandler.GetProfile)
}
