package controller

import (
	"oj-back/app/service"
	"oj-back/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func JudgeCode(c *fiber.Ctx) error {
	payload := struct {
		Language  string `json:"language"`
		ProblemID int    `json:"problem_id"`
		Code      string `json:"code"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return utils.HandleError(c, err, "请求格式错误")
	}

	problemCases, err := utils.GetTestCases(payload.ProblemID)
	if err != nil {
		return utils.HandleError(c, err, "获取测试用例失败")
	}

	evaluation, err := service.EvaluateProblem(payload.Language, payload.Code, problemCases)
	if err != nil {
		return utils.HandleError(c, err, "评测失败")
	}

	return c.Status(fiber.StatusOK).JSON(evaluation)
}
