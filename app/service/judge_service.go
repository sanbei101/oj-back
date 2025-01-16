package service

import (
	"encoding/base64"
	"fmt"
	"oj-back/app/model"
	"oj-back/pkg/utils"
	"sync"
)

type JudgeService struct{}

var JudgeServiceApp = new(JudgeService)

func (js *JudgeService) EvaluateProblem(language string, codeContent string, testCases []model.Case) (*model.EvaluationResult, error) {
	var executablePath string
	var compileErr error

	// 针对编译型语言进行编译
	if language == "c" {
		// 解码代码
		decodedCode, err := base64.StdEncoding.DecodeString(codeContent)
		if err != nil {
			return nil, fmt.Errorf("解码代码失败: %v", err)
		}

		// 编译C代码
		executablePath, compileErr = utils.CompileCCode(decodedCode)
		if compileErr != nil {
			// 如果编译失败，所有测试用例均标记为编译错误
			results := make([]model.TestResult, len(testCases))
			for i := range results {
				results[i] = model.TestResult{
					IsSuccess:      false,
					ExpectedOutput: testCases[i].ExpectedOutput,
					ActualOutput:   compileErr.Error(),
				}
			}
			evaluation := &model.EvaluationResult{
				Count:   len(results),
				Results: results,
			}
			return evaluation, nil
		}
	}

	var wg sync.WaitGroup
	results := make([]model.TestResult, len(testCases))

	for i, testCase := range testCases {
		wg.Add(1)
		go func(i int, tc model.Case) {
			defer wg.Done()

			var output string
			var err error

			if language == "c" {
				// 对编译型语言语言，运行已编译的可执行文件
				output, err = utils.SubmitExecutableTask(executablePath, tc.Input)
			} else if language == "python" {
				// 对于解释型语言，使用现有的RunCode函数
				output, err = utils.SubmitCodeExecution(language, codeContent, tc.Input)
			} else {
				err = fmt.Errorf("不支持的语言")
			}

			if err != nil {
				results[i] = model.TestResult{
					IsSuccess:      false,
					ExpectedOutput: tc.ExpectedOutput,
					ActualOutput:   err.Error(),
				}
				return
			}

			// 比较输出
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
