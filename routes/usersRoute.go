package routes

import (
	"gofiber/controllers"
	"gofiber/middleware"

	"github.com/gofiber/fiber/v2"
)

func AddUsersRoute(app *fiber.App) {

	apix := app.Group("/apix")
	v1 := apix.Group("/v1")
	users := v1.Group("/users")

	users.Get("/all", controllers.GetAllUser)
	users.Post("/login", controllers.Login)
	users.Post("/register", controllers.Register)
	users.Get("/restricted", controllers.Restricted)
	users.Get("/", controllers.Accessible)

	users.Use("/me", middleware.AuthConfig, middleware.CheckAuthFromId)
	users.Get("/me", controllers.GetUserByMe)
}
