package controllers

import (
	"context"
	"fmt"
	"gofiber/db"
	"gofiber/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

const COLLECTION_NAME = "products"

func GetAllProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := db.ConnectMongoDB()
	collection := client.Collection(COLLECTION_NAME)
	var products []models.Product
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
		fmt.Println("err")
	}
	if err = cursor.All(ctx, &products); err != nil {
		fmt.Println(err)
	}

	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"resultCode": strconv.Itoa(fiber.StatusOK * 100),
		"resultData": products,
	})
}
