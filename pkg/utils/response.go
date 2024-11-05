package utils

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Response 结构体
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// 通用错误处理函数
func HandleError(c *fiber.Ctx, err error, message string) error {
	if err != nil {
		log.Println(err)
		var statusCode int

		switch {
		case err == gorm.ErrRecordNotFound:
			statusCode = fiber.StatusNotFound
		case err.Error() == "validation error":
			statusCode = fiber.StatusBadRequest
		default:
			statusCode = fiber.StatusInternalServerError
		}

		response := Response{
			Success: false,
			Message: message,
			Data:    err.Error(),
		}
		return c.Status(statusCode).JSON(response)
	} else {
		response := Response{
			Success: false,
			Message: message,
		}
		return c.JSON(response)
	}
}
