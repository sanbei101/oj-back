package main

import (
	"oj-back/internal/controller"
	"oj-back/internal/db"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// 初始化数据库
	db.InitDB()
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Post("/judge", controller.JudgeCode)

	app.Listen("0.0.0.0:3000")
}
