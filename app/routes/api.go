package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

func ApiRoute(app *fiber.App) {
	fmt.Println("we here")
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hi, we up")
	})
}
