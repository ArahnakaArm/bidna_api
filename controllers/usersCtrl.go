package controllers

import (
	"context"
	"fmt"
	"gofiber/db"
	"gofiber/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const COLLECTION_USERS = "users"
const COLLECTION_USERS_TOKENS = "users_tokens"

func GetAllUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := db.ConnectMongoDB()
	collection := client.Collection(COLLECTION_USERS)
	var users []models.User
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
		fmt.Println("err")
	}
	if err = cursor.All(ctx, &users); err != nil {
		fmt.Println(err)
	}

	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"resultCode": strconv.Itoa(fiber.StatusOK * 100),
		"resultData": users,
		"rowCount":   len(users),
	})
}

func Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := db.ConnectMongoDB()
	user := new(models.User)
	var fuser models.User
	collection := client.Collection(COLLECTION_USERS)
	tokens := client.Collection(COLLECTION_USERS_TOKENS)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})

	}

	var exist = collection.FindOne(ctx, bson.M{
		"username": user.Username,
	})

	err := exist.Decode(&fuser)
	fmt.Println(fuser)
	if err == nil {
		err := bcrypt.CompareHashAndPassword([]byte(fuser.Password), []byte(user.Password))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "รหัสผ่านไม่ถูกต้อง",
			})
		}
	} else {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "ไม่มีบัญชีผู่ใช้นี้ในระบบ",
		})
	}

	// Throws Unauthorized error
	/* if username != "john" || passname != "doe" {
		return c.SendStatus(fiber.StatusUnauthorized)
	} */

	// Create token

	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = fuser.Username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	filter := bson.M{
		"user_id": fuser.Id,
	}
	efef := bson.M{
		"$set": bson.M{
			"token":      t,
			"updated_at": time.Now(),
		},
		"$setOnInsert": bson.M{
			"user_id":    fuser.Id,
			"created_at": time.Now(),
		},
	}
	result := tokens.FindOneAndUpdate(ctx, filter, efef, &opt)

	if result.Err() != nil {
		fmt.Println(result.Err())
	}

	if err != nil {
		fmt.Println(err)
		fmt.Println("err")
	}

	return c.JSON(fiber.Map{"token": t})
}

func Register(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := db.ConnectMongoDB()
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})

	}

	hash, _ := hashPassword(user.Password)
	collection := client.Collection(COLLECTION_USERS)
	var users []models.User

	var exist = collection.FindOne(ctx, bson.M{
		"username": user.Username,
	})

	err := exist.Decode(user)
	fmt.Println(user)
	if err == nil {
		fmt.Println("exist")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "มีบัญชีผู้ใช้นี้แล้วในระบบ",
		})
	}

	result, err := collection.InsertOne(ctx, bson.M{
		"username":   user.Username,
		"password":   hash,
		"first_name": user.First_name,
		"last_name":  user.Last_name,
		"credit":     0.00,
		"created_at": time.Now(),
	})

	if err != nil {
		fmt.Println(err)
		fmt.Println("err")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"resultCode": strconv.Itoa(fiber.StatusOK * 100),
		"resultData": result,
		"rowCount":   len(users),
	})
}

func Accessible(c *fiber.Ctx) error {
	return c.SendString("Accessible")
}

func Restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
