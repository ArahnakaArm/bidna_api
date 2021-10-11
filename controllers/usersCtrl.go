package controllers

import (
	"context"
	"fmt"
	"gofiber/models"
	"gofiber/responseMessage"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserController interface {
	GetAllUser(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error
	GetUserByMe(ctx *fiber.Ctx) error
	ChangePassword(ctx *fiber.Ctx) error
}

type userController struct {
	database *mongo.Database
}

func NewUserController(database *mongo.Database) UserController {
	return &userController{database}
}

const COLLECTION_USERS = "users"
const COLLECTION_USERS_TOKENS = "users_tokens"

func (r *userController) GetAllUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := r.database.Collection(COLLECTION_USERS)
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

func (r *userController) Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user := new(models.User)
	var fuser models.User
	collection := r.database.Collection(COLLECTION_USERS)
	tokens := r.database.Collection(COLLECTION_USERS_TOKENS)
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
	claims["id"] = fuser.Id
	claims["name"] = fuser.Username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(viper.GetString("appAuth.tokenSecret")))
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

func (r *userController) Register(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})

	}

	hash, _ := hashPassword(user.Password)
	collection := r.database.Collection(COLLECTION_USERS)
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

func (r *userController) GetUserByMe(c *fiber.Ctx) error {
	splitToken := strings.Split(c.Get("authorization"), "Bearer ")
	reqToken := splitToken[1]
	token, err := jwt.Parse(reqToken, nil)
	if token == nil {
		return err
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	id := claims["id"].(string)
	fmt.Println(id)
	var fuser models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := r.database.Collection(COLLECTION_USERS)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusInternalServerError * 100),
			"resultMessage": responseMessage.RESULT_MESSAGE_INTERNAL_ERROR,
		})
	}

	var result = collection.FindOne(ctx, bson.M{
		"_id": objectId,
	})

	err = result.Decode(&fuser)

	if err != nil {
		return c.Status(200).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusNoContent * 100),
			"resultMessage": responseMessage.RESULT_MESSAGE_DATA_NOT_FOUND,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"resultCode": strconv.Itoa(fiber.StatusOK * 100),
		"resultData": fuser,
	})
}

func (r *userController) ChangePassword(c *fiber.Ctx) error {
	splitToken := strings.Split(c.Get("authorization"), "Bearer ")
	reqToken := splitToken[1]
	token, err := jwt.Parse(reqToken, nil)
	if token == nil {
		return err
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	id := claims["id"].(string)
	fmt.Println(id)
	var fuser models.User
	changePassRequest := new(models.ChangePasswordRequest)

	if err := c.BodyParser(changePassRequest); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})

	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.database.Collection(COLLECTION_USERS)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusInternalServerError * 100),
			"resultMessage": responseMessage.RESULT_MESSAGE_INTERNAL_ERROR,
		})
	}

	var result = collection.FindOne(ctx, bson.M{
		"_id": objectId,
	})

	if err != nil {
		return c.Status(200).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusNoContent * 100),
			"resultMessage": responseMessage.RESULT_MESSAGE_DATA_NOT_FOUND,
		})
	}

	err = result.Decode(&fuser)

	if err == nil {
		err := bcrypt.CompareHashAndPassword([]byte(fuser.Password), []byte(changePassRequest.OldPassword))
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

	newPasswordHash, _ := hashPassword(changePassRequest.NewPassword)

	update := bson.M{
		"$set": bson.M{"password": newPasswordHash},
	}

	updateResult := collection.FindOneAndUpdate(ctx, bson.M{"_id": objectId}, update)

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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"resultCode":    strconv.Itoa(fiber.StatusOK * 100),
		"resultMessage": responseMessage.RESULT_MESSAGE_SUCCESS,
	})
}

type CustomClaimsExample struct {
	*jwt.StandardClaims
	TokenType string
	CustomerInfo
}
type CustomerInfo struct {
	Name string
	Kind string
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
