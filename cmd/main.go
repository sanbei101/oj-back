package main

import (
	"oj-back/app/controller"
	"oj-back/app/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// 初始化数据库
	db.InitDB()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	// 判题逻辑
	judgeGroup := app.Group("/judge")
	judgeGroup.Post("/submit", controller.JudgeCode)

	// 核心逻辑
	coreGroup := app.Group("/core")
	coreGroup.Get("/problem", controller.GetAllProblems)

	app.Listen("0.0.0.0:3000")
}
