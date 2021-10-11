package controllers

import (
	"bytes"
	"fmt"
	"gofiber/responseMessage"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/spf13/viper"
)

var (
	s3session *s3.S3
)

func UploadFile(c *fiber.Ctx) error {
	fmt.Println(viper.GetString("s3.region"))
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(viper.GetString("s3.region")),
		Credentials: credentials.NewStaticCredentials(viper.GetString("s3.accessKeyId"), viper.GetString("s3.accessSecretKey"), ""),
	},
	)
	/* sess := session.Must(session.NewSession()) */

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	/* 	f, err := os.Open(c.FormFile("file")) */

	data, err := c.FormFile("file")

	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusForbidden * 100),
			"resultMessage": responseMessage.RESULT_MESSAGE_MISSING_PARAMETER,
		})
	}

	fileType := data.Header.Get("Content-Type")

	fileType = strings.Split(fileType, "/")[1]

	file, err := data.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	uId, err := uuid.NewV4()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusInternalServerError * 100),
			"resultMessage": responseMessage.RESULT_MESSAGE_INTERNAL_ERROR,
		})
	}

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(viper.GetString("s3.bucketName")),
		Key:    aws.String(fmt.Sprintf("profileImage/%s.%s", uId.String(), fileType)),
		Body:   bytes.NewReader(fileData),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))

	path := &UploadFileRes{
		Path: aws.StringValue(&result.Location),
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"resultCode":    strconv.Itoa(fiber.StatusOK * 100),
		"resultMessage": responseMessage.RESULT_MESSAGE_SUCCESS,
		"resultData":    path,
	})
	/* uploadPath := c.Params("path")

	file, err := c.FormFile("file")

	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusForbidden * 100),
			"resultMessage": responseMessage.RESULT_MESSAGE_MISSING_PARAMETER,
		})
	}

	uId, err := uuid.NewV4()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"resultCode":    strconv.Itoa(fiber.StatusInternalServerError * 100),
			"resultMessage": responseMessage.RESULT_MESSAGE_INTERNAL_ERROR,
		})
	}

	c.SaveFile(file, fmt.Sprintf("./uploads/%s/%s.png", uploadPath, uId))

	path := &UploadFileRes{
		Path: fmt.Sprintf("/uploads/%s/%s.png", uploadPath, uId),
	}


	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"resultCode":    strconv.Itoa(fiber.StatusOK * 100),
		"resultMessage": responseMessage.RESULT_MESSAGE_SUCCESS,
		"resultData":    path,
	})

	*/
}

type UploadFileRes struct {
	Path string `json:"path"`
}
