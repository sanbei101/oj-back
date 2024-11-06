package utils

import (
	"github.com/gofiber/fiber/v2"
)

// Response 结构体
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// 通用错误处理函数
func HandleError(c *fiber.Ctx, err error, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error":   message,
		"details": err.Error(),
	})
}
