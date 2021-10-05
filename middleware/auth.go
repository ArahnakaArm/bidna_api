package middleware

import (
	"context"
	"fmt"
	"gofiber/db"
	"gofiber/models"
	"gofiber/responseMessage"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SuccessValidate(c *fiber.Ctx) error {
	fmt.Println("WORK BEFORE0")
	return c.Next()
}

func FailAuth(c *fiber.Ctx, e error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"resultCode":    strconv.Itoa(fiber.StatusUnauthorized * 100),
		"resultMessage": responseMessage.RESULT_UNAUTHORIZED,
	})
}

var AuthConfig = jwtware.New(jwtware.Config{
	SigningMethod:  "HS256",
	SigningKey:     []byte("secret"),
	SuccessHandler: SuccessValidate,
	ErrorHandler:   FailAuth,
})

func CheckAuthFromId(c *fiber.Ctx) error {
	const COLLECTION_USERS = "users"
	splitToken := strings.Split(c.Get("authorization"), "Bearer ")
	reqToken := splitToken[1]
	token, err := jwt.Parse(reqToken, nil)
	if token == nil {
		return err
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	id := claims["id"].(string)
	var fuser models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := db.ConnectMongoDB()
	collection := client.Collection(COLLECTION_USERS)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusInternalServerError * 100),
			"resultMessage": responseMessage.RESULT_MESSAGE_INTERNAL_ERROR,
		})
	}

	var exist = collection.FindOne(ctx, bson.M{
		"_id": objectId,
	})

	err = exist.Decode(&fuser)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusUnauthorized * 100),
			"resultMessage": responseMessage.RESULT_UNAUTHORIZED,
		})
	}

	return c.Next()

}
