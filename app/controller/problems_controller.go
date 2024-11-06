package controller

import (
	"oj-back/app/service"
	"oj-back/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func GetAllProblems(c *fiber.Ctx) error {
	page := c.QueryInt("page")
	size := c.QueryInt("size")
	keyword := c.Query("keyword")
	if page == 0 || size == 0 {
		page = 1
		size = 10
	}
	problems, err := service.GetAllProblems(page, size, keyword)
	if err != nil {
		utils.HandleError(c, err, "获取题目列表失败")
	}

	return c.Status(fiber.StatusOK).JSON(problems)
}

func GetProblemByID(c *fiber.Ctx) error {
	id := c.QueryInt("id")
	problem, err := service.GetProblemByID(id)
	if err != nil {
		utils.HandleError(c, err, "获取题目失败")
	}
	return c.Status(fiber.StatusOK).JSON(problem)
}
