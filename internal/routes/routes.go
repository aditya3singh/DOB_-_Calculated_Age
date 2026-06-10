package routes

import (
	"user-api/internal/handler"
	"user-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, userHandler *handler.UserHandler) {
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	users := app.Group("/users")
	users.Post("/", userHandler.CreateUser)
	users.Get("/", userHandler.ListUsers)
	users.Get("/:id", userHandler.GetUserByID)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
}
