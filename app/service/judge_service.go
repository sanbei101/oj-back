package service

import (
	"oj-back/app/models"
	"oj-back/pkg/utils"
	"sync"
)

// 评测结果结构体
type TestResult struct {
	IsSuccess      bool   `json:"is_success"`
	ExpectedOutput string `json:"expected_output"`
	ActualOutput   string `json:"actual_output"`
}

type EvaluationResult struct {
	Count   int          `json:"count"`
	Results []TestResult `json:"results"`
}

func EvaluateProblem(language string, codeContent string, testCases []models.Case) (*EvaluationResult, error) {
	var wg sync.WaitGroup
	results := make([]TestResult, len(testCases))

	for i, testCase := range testCases {
		wg.Add(1)
		go func(i int, tc models.Case) {
			defer wg.Done()

			// 执行用户代码并获取输出
			output, err := utils.RunCode(language, codeContent, tc.Input)
			if err != nil {
				results[i] = TestResult{
					IsSuccess:      false,
					ExpectedOutput: tc.ExpectedOutput,
					ActualOutput:   err.Error(),
				}
				return
			}
			isCorrect := utils.CompareOutput(output, tc.ExpectedOutput)
			results[i] = TestResult{
				IsSuccess:      isCorrect,
				ExpectedOutput: tc.ExpectedOutput,
				ActualOutput:   output,
			}
		}(i, testCase)
	}

	wg.Wait()

	evaluation := &EvaluationResult{
		Count:   len(results),
		Results: results,
	}

	return evaluation, nil
}
