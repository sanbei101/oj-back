package controller

import (
	"oj-back/app/service"
	"oj-back/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

const (
	defaultPage = 1
	defaultSize = 10
)

func GetAllProblems(c *fiber.Ctx) error {
	page, size, keyword := c.QueryInt("page"), c.QueryInt("size"), c.Query("keyword")
	if page == 0 || size == 0 {
		page = defaultPage
		size = defaultSize
	}
	problems, err := service.ProblemServiceApp.GetAllProblems(page, size, keyword)
	if err != nil {
		utils.HandleError(c, err, "获取题目列表失败")
	}

	return c.Status(fiber.StatusOK).JSON(problems)
}

func GetProblemByID(c *fiber.Ctx) error {
	id := c.QueryInt("id")
	problem, err := service.ProblemServiceApp.GetProblemByID(id)
	if err != nil {
		utils.HandleError(c, err, "获取题目失败")
	}

	return c.Status(fiber.StatusOK).JSON(problem)
}

func GetProblemTestCase(c *fiber.Ctx) error {
	problemID := c.QueryInt("problem_id")
	cases, err := service.ProblemServiceApp.GetProblemTestCase(uint64(problemID))
	if err != nil {
		utils.HandleError(c, err, "获取测试用例失败")
	}

	return c.Status(fiber.StatusOK).JSON(cases)
}
