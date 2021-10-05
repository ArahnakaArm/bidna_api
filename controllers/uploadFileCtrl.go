package controllers

import (
	"fmt"
	"gofiber/responseMessage"
	"strconv"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/nu7hatch/gouuid"
)

func UploadFile(c *fiber.Ctx) error {

	uploadPath := c.Params("path")

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

	// Save file to root directory:
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"resultCode":    strconv.Itoa(fiber.StatusOK * 100),
		"resultMessage": responseMessage.RESULT_MESSAGE_SUCCESS,
		"resultData":    path,
	})

	/* 	c.SendString(fmt.Sprintf("/uploads/%s/%s.png", uploadPath, uId)) */
}

type UploadFileRes struct {
	Path string `json:"path"`
}
