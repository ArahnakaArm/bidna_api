package controllers

import (
	"context"
	"fmt"
	"gofiber/db"
	"gofiber/models"
	"gofiber/responseMessage"
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

	if c.Query("category") != "" {
		query["category"] = c.Query("category")
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

	if products == nil {
		products = []models.Product{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"resultCode":    strconv.Itoa(fiber.StatusOK * 100),
		"resultMessage": responseMessage.RESULT_MESSAGE_SUCCESS,
		"resultData":    products,
		"rowCount":      len(products),
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
		return c.Status(200).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusNoContent * 100),
			"resultMessage": responseMessage.RESULT_MESSAGE_DATA_NOT_FOUND,
		})
	}

	err = findResult.Decode(&product)
	if err != nil {
		fmt.Println(err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"resultCode":    strconv.Itoa(fiber.StatusOK * 100),
		"resultMessage": responseMessage.RESULT_MESSAGE_SUCCESS,
		"resultData":    product,
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
			"resultMessage": responseMessage.RESULT_MESSAGE_MISSING_PARAMETER,
		})
	}

	if errs := validator.Validate(product); errs != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusForbidden * 100),
			"resultMessage": responseMessage.RESULT_MESSAGE_MISSING_PARAMETER,
		})
	}

	productCheckConflictCount, err := collection.CountDocuments(ctx, bson.M{"name": product.Name})

	if err != nil {
		panic(err)
	}
	if productCheckConflictCount >= 1 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusConflict * 100),
			"resultMessage": responseMessage.RESULT_MESSAGE_CONFLICT,
		})
	}

	result, err := collection.InsertOne(ctx, product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusInternalServerError * 100),
			"resultMessage": responseMessage.RESULT_MESSAGE_INTERNAL_ERROR,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"resultCode":    strconv.Itoa(fiber.StatusCreated * 100),
		"resultMessage": responseMessage.RESULT_MESSAGE_CREATED,
		"resultData":    result,
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
			"resultMessage": responseMessage.RESULT_MESSAGE_MISSING_PARAMETER,
		})
	}

	if errs := validator.Validate(product); errs != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusForbidden * 100),
			"resultMessage": responseMessage.RESULT_MESSAGE_MISSING_PARAMETER,
		})
	}

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		fmt.Println(err)
	}

	update := bson.M{
		"$set": product,
	}

	updateResult := collection.FindOneAndUpdate(ctx, bson.M{"_id": objId}, update)

	if updateResult.Err() != nil {
		return c.Status(200).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusNoContent * 100),
			"resultMessage": responseMessage.RESULT_MESSAGE_DATA_NOT_FOUND,
		})
	}

	if updateResult == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusInternalServerError * 100),
			"resultMessage": responseMessage.RESULT_MESSAGE_INTERNAL_ERROR,
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusInternalServerError * 100),
			"resultMessage": responseMessage.RESULT_MESSAGE_INTERNAL_ERROR,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"resultCode":    strconv.Itoa(fiber.StatusOK * 100),
		"resultMessage": responseMessage.RESULT_MESSAGE_SUCCESS,
		"resultData":    product,
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

	deleteResult := collection.FindOneAndDelete(ctx, bson.M{"_id": objId})

	if deleteResult.Err() != nil {
		return c.Status(200).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusNoContent * 100),
			"resultMessage": responseMessage.RESULT_MESSAGE_DATA_NOT_FOUND,
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusInternalServerError * 100),
			"resultMessage": responseMessage.RESULT_MESSAGE_INTERNAL_ERROR,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"resultCode":    strconv.Itoa(fiber.StatusOK * 100),
		"resultMessage": responseMessage.RESULT_MESSAGE_SUCCESS,
	})
}
