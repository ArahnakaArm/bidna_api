package routes

import (
	"gofiber/controllers"

	"github.com/gofiber/fiber/v2"
)

func AddUploadFileRoute(app *fiber.App) {

	apix := app.Group("/apix")
	v1 := apix.Group("/v1")
	upload := v1.Group("/upload")

	upload.Post("/:path", controllers.UploadFile)

}
