package controller

import (
	"oj-back/internal/service"
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
		utils.HandleError(c, err, "请求格式错误")
	}

	problemCases, err := utils.GetTestCases(payload.ProblemID)
	if err != nil {
		utils.HandleError(c, err, "获取测试用例失败")
	}

	evaluation, err := service.EvaluateProblem(payload.Language, payload.Code, problemCases)
	if err != nil {
		utils.HandleError(c, err, "评测失败")
	}

	response := utils.Response{
		Success: true,
		Message: "评测完成",
		Data:    evaluation,
	}

	return c.Status(fiber.StatusOK).JSON(response)

}
