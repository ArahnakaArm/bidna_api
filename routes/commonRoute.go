package routes

import (
	"github.com/gofiber/fiber/v2"
)

func AddCommonRoute(app *fiber.App) {

	apix := app.Group("/apix")

	v1 := apix.Group("/v1")

	v1.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Helloo")
	})
}
