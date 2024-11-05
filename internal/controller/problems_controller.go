package controller

import (
	"oj-back/internal/service"
	"oj-back/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func GetAllProblems(c *fiber.Ctx) error {
	problems, err := service.GetAllProblems()
	if err != nil {
		utils.HandleError(c, err, "获取题目列表失败")
	}

	response := utils.Response{
		Success: true,
		Message: "获取题目列表成功",
		Data:    problems,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func GetProblemByID(c *fiber.Ctx) error {
	id := c.QueryInt("id")
	problem, err := service.GetProblemByID(id)
	if err != nil {
		utils.HandleError(c, err, "获取题目失败")
	}

	response := utils.Response{
		Success: true,
		Message: "获取题目成功",
		Data:    problem,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
