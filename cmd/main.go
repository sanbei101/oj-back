package main

import (
	"oj-back/app/controller"
	"oj-back/app/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

func main() {
	// 初始化数据库
	db.InitDB()
	app := fiber.New()
	app.Use(pprof.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Post("/judge", controller.JudgeCode)
	problemsGroup := app.Group("/problems")

	problemsGroup.Get("/", controller.GetAllProblems)
	problemsGroup.Get("/problem", controller.GetProblemByID)
	app.Listen("0.0.0.0:3000")
}
