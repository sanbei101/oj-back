package service

import (
	"oj-back/app/model"
	"oj-back/pkg/utils"
	"sync"
)

type JudgeService struct{}

var JudgeServiceApp = new(JudgeService)

func (js *JudgeService) EvaluateProblem(language string, codeContent string, testCases []model.Case) (*model.EvaluationResult, error) {
	var wg sync.WaitGroup
	results := make([]model.TestResult, len(testCases))

	for i, testCase := range testCases {
		wg.Add(1)
		go func(i int, tc model.Case) {
			defer wg.Done()

			// 执行用户代码并获取输出
			output, err := utils.RunCode(i, language, codeContent, tc.Input)
			if err != nil {
				results[i] = model.TestResult{
					IsSuccess:      false,
					ExpectedOutput: tc.ExpectedOutput,
					ActualOutput:   err.Error(),
				}
				return
			}
			isCorrect, isStrictlyCorrect := utils.CompareOutput(output, tc.ExpectedOutput)
			results[i] = model.TestResult{
				IsSuccess:       isCorrect,
				IsStrictSuccess: isStrictlyCorrect,
				ExpectedOutput:  tc.ExpectedOutput,
				ActualOutput:    output,
			}
		}(i, testCase)
	}

	wg.Wait()

	evaluation := &model.EvaluationResult{
		Count:   len(results),
		Results: results,
	}

	return evaluation, nil
}
