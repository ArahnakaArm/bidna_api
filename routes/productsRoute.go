package routes

import (
	"gofiber/controllers"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddProductsRoute(app *fiber.App, db *mongo.Database) {

	apix := app.Group("/apix")
	v1 := apix.Group("/v1")
	products := v1.Group("/products")

	products.Get("/", controllers.NewProductController(db).GetAllProduct)

	products.Get("/:id", controllers.NewProductController(db).GetProduct)

	products.Post("/", controllers.NewProductController(db).AddProduct)

	products.Put("/:id", controllers.NewProductController(db).UpdateProduct)

	products.Delete("/:id", controllers.NewProductController(db).DeleteProduct)

}
