package server

import (
	"github.com/gofiber/fiber/v2"

	"api-stockflow/internal/database"
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
