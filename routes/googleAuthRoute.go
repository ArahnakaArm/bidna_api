package routes

import (
	"gofiber/controllers"

	"github.com/gofiber/fiber/v2"
)

func AddGoogleAuthRoute(app *fiber.App) {

	apix := app.Group("/apix")
	v1 := apix.Group("/v1")
	googleAuth := v1.Group("/google_auth")

	googleAuth.Get("/", controllers.HandleGoogleLogin)

	googleAuth.Get("/callback", controllers.HandleGoogleCallBack)

}
