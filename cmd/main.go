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
	judgeGroup.Post("/submit-code", controller.JudgeCode)

	//
	problemGroup := app.Group("/problem")
	problemGroup.Get("/get-all-problem", controller.GetAllProblems)
	problemGroup.Get("/get-problem-by-id", controller.GetProblemByID)
	problemGroup.Get("/get-problem-test-case", controller.GetProblemTestCase)
	app.Listen("0.0.0.0:3000")
}
