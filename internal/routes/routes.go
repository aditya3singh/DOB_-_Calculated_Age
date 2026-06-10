package routes

import (
	"user-api/internal/handler"
	"user-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// Register wires all application routes and middleware onto the Fiber app.
func Register(app *fiber.App, userHandler *handler.UserHandler) {
	// Global middleware
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// User routes
	users := app.Group("/users")
	users.Post("/", userHandler.CreateUser)
	users.Get("/", userHandler.ListUsers)
	users.Get("/:id", userHandler.GetUserByID)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
}
