package routes

import (
	"gofiber/controllers"

	"github.com/gofiber/fiber/v2"
)

func AddUsersRoute(app *fiber.App) {

	apix := app.Group("/apix")
	v1 := apix.Group("/v1")
	products := v1.Group("/users")

	products.Get("/all", controllers.GetAllUser)
}
