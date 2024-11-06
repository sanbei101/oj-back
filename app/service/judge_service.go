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

// 评测函数，循环遍历每个测试用例并进行评测
func EvaluateProblem(language string, codeContent string, testCases []models.Case) (*EvaluationResult, error) {
	var results []TestResult
	var mu sync.Mutex
	var wg sync.WaitGroup
	var firstErr error
	var once sync.Once

	for _, testCase := range testCases {
		wg.Add(1)
		go func(tc models.Case) {
			defer wg.Done()

			// 执行用户代码并获取输出
			output, err := utils.RunCode(language, codeContent, tc.Input)
			if err != nil {
				once.Do(func() {
					firstErr = err
				})
				return
			}

			// 比对输出
			isCorrect := utils.CompareOutput(output, tc.ExpectedOutput)

			// 记录每个测试结果
			mu.Lock()
			results = append(results, TestResult{
				IsSuccess:      isCorrect,
				ExpectedOutput: tc.ExpectedOutput,
				ActualOutput:   output,
			})
			mu.Unlock()
		}(testCase)
	}

	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	// 生成总的评测结果
	evaluation := &EvaluationResult{
		Count:   len(results),
		Results: results,
	}

	return evaluation, nil
}
