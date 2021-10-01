package controllers

import (
	"context"
	"fmt"
	"gofiber/db"
	"gofiber/models"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/validator.v2"
)

const COLLECTION_NAME = "products"

func GetAllProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := db.ConnectMongoDB()
	collection := client.Collection(COLLECTION_NAME)
	var products []models.Product

	query := bson.M{}
	if c.Query("name") != "" {
		query["name"] = c.Query("name")
	}

	cursor, err := collection.Find(ctx, query)
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
		"rowCount":   len(products),
	})
}

func GetProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := db.ConnectMongoDB()
	collection := client.Collection(COLLECTION_NAME)
	var product models.Product
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		fmt.Println(err)
	}
	findResult := collection.FindOne(ctx, bson.M{"_id": objId})
	if err := findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusNoContent * 100),
			"resultMessage": "Not Found",
		})
	}

	err = findResult.Decode(&product)
	if err != nil {
		fmt.Println(err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"resultCode": strconv.Itoa(fiber.StatusOK * 100),
		"resultData": product,
	})
}

func AddProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := db.ConnectMongoDB()
	collection := client.Collection(COLLECTION_NAME)
	product := new(models.Product)

	if err := c.BodyParser(product); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusForbidden * 100),
			"resultMessage": "Missing Or Invalid Parameter",
		})
	}

	if errs := validator.Validate(product); errs != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusForbidden * 100),
			"resultMessage": "Missing Or Invalid Parameter",
		})
	}

	productCheckConflictCount, err := collection.CountDocuments(ctx, bson.M{"name": product.Name})

	if err != nil {
		panic(err)
	}
	if productCheckConflictCount >= 1 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusConflict * 100),
			"resultMessage": "Conflict",
		})
	}

	result, err := collection.InsertOne(ctx, product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusInternalServerError * 100),
			"resultMessage": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"resultCode": strconv.Itoa(fiber.StatusCreated * 100),
		"resultData": result,
	})
}

func UpdateProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := db.ConnectMongoDB()
	collection := client.Collection(COLLECTION_NAME)

	product := new(models.Product)

	if err := c.BodyParser(product); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusForbidden * 100),
			"resultMessage": "Missing Or Invalid Parameter",
		})
	}

	if errs := validator.Validate(product); errs != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusForbidden * 100),
			"resultMessage": "Missing Or Invalid Parameter",
		})
	}

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		fmt.Println(err)
	}

	update := bson.M{
		"$set": product,
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objId}, update)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusInternalServerError * 100),
			"resultMessage": "Internal Server Error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"resultCode": strconv.Itoa(fiber.StatusOK * 100),
		"resultData": product,
	})
}

func DeleteProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := db.ConnectMongoDB()
	collection := client.Collection(COLLECTION_NAME)

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		fmt.Println(err)
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objId})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusInternalServerError * 100),
			"resultMessage": "Internal Server Error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"resultCode":    strconv.Itoa(fiber.StatusOK * 100),
		"resultMessage": "Success",
	})
}
