package service

import "oj-back/pkg/utils"

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
func EvaluateProblem(language string, codeContent string, testCases []utils.TestCase) (*EvaluationResult, error) {
	var results []TestResult

	for _, testCase := range testCases {
		// 执行用户代码并获取输出
		output, err := utils.RunCode(language, codeContent, testCase.Input)
		if err != nil {
			return nil, err
		}

		// 比对输出
		isCorrect := utils.CompareOutput(output, testCase.ExpectedOutput)

		// 记录每个测试结果
		results = append(results, TestResult{
			IsSuccess:      isCorrect,
			ExpectedOutput: testCase.ExpectedOutput,
			ActualOutput:   output,
		})
	}

	// 生成总的评测结果
	evaluation := &EvaluationResult{
		Count:   len(results),
		Results: results,
	}

	return evaluation, nil
}
