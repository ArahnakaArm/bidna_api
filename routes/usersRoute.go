package routes

import (
	"gofiber/controllers"
	"gofiber/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddUsersRoute(app *fiber.App, db *mongo.Database) {

	apix := app.Group("/apix")
	v1 := apix.Group("/v1")
	users := v1.Group("/users")

	users.Get("/all", controllers.NewUserController(db).GetAllUser)
	users.Post("/login", controllers.NewUserController(db).Login)
	users.Post("/register", controllers.NewUserController(db).Register)
	users.Get("/restricted", controllers.Restricted)
	users.Get("/", controllers.Accessible)

	users.Use("/me", middleware.AuthConfig, middleware.CheckAuthFromId)
	users.Get("/me", controllers.NewUserController(db).GetUserByMe)

	users.Use("/changepassword", middleware.AuthConfig, middleware.CheckAuthFromId)
	users.Post("/changepassword", controllers.NewUserController(db).ChangePassword)

}
