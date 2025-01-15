package controller

import (
	"github.com/google/uuid"
	"oj-back/app/model"
	"oj-back/app/service"
	"oj-back/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func JudgeCode(c *fiber.Ctx) error {
	payload := struct {
		UserID    uuid.UUID `json:"user_id"`
		Language  string    `json:"language"`
		ProblemID uint64    `json:"problem_id"`
		Code      string    `json:"code"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return utils.HandleError(c, err, "请求格式错误")
	}

	submit := model.Submit{
		SubmitID:    uuid.New(),
		UserID:      payload.UserID,
		Language:    payload.Language,
		ProblemID:   payload.ProblemID,
		CodeContent: payload.Code,
	}

	service.JudgeServiceApp.EvaluateSubmit(&submit)

	return c.Status(fiber.StatusOK).JSON(submit.EvaResult)
}
