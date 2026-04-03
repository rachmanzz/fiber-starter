package routes

import (
	"github.com/gofiber/fiber/v3"
)

func ApiRoute(app *fiber.App) {
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hi, we up")
	})
}
