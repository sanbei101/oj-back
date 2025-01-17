package service

import (
	"encoding/base64"
	"fmt"
	"oj-back/app/model"
	"oj-back/pkg/judge"
	"os"
	"strings"
	"sync"
)

var compiledLanguages = map[string]bool{
	"c":    true,
	"cpp":  true,
	"java": true,
	"go":   true,
}

type JudgeService struct{}

var JudgeServiceApp = new(JudgeService)

func CompareOutput(actualOutput string, expectedOutput string) (same bool, strictlySame bool) {
	// TODO: 似乎逻辑有点不足,例如每一行行尾的/r/n或/n无法区分,需要改进一下.
	actual := strings.TrimSpace(actualOutput)
	expected := strings.TrimSpace(expectedOutput)
	return actual == expected, expectedOutput == actualOutput
}

func (js *JudgeService) EvaluateProblem(language string, codeContent string, testCases []model.Case) (*model.EvaluationResult, error) {
	// 客户端发送的代码是经过base64编码的
	decodedCode, err := base64.StdEncoding.DecodeString(codeContent)

	if err != nil {
		return nil, fmt.Errorf("解码代码失败: %v", err)
	}

	// 针对编译型语言进行编译
	var executablePath string
	if compiledLanguages[language] {
		switch language {
		case "c":
			executablePath, err = judge.CompiledJudgeApp.CompileCCode(decodedCode)
		case "cpp":
			executablePath, err = judge.CompiledJudgeApp.CompileCppCode(decodedCode)
		default:
			return nil, fmt.Errorf("还不支持这个语言:%s", language)
		}
		// 完成之后删除可执行文件
		defer os.Remove(executablePath)
		if err != nil {
			results := make([]model.TestResult, len(testCases))
			for i := range results {
				results[i] = model.TestResult{
					IsSuccess:      false,
					ExpectedOutput: testCases[i].ExpectedOutput,
					ActualOutput:   err.Error(),
				}
			}
			evaluation := &model.EvaluationResult{
				Count:   len(results),
				Results: results,
			}
			return evaluation, nil
		}

		// 并发执行测试用例
		var wg sync.WaitGroup
		results := make([]model.TestResult, len(testCases))
		wg.Add(len(testCases))
		for i, testCase := range testCases {
			go func(i int, tc model.Case) {
				defer wg.Done()
				// 对编译型语言语言，运行已编译的可执行文件
				output, err := judge.CompiledJudgeApp.SubmitJudge(executablePath, tc.Input)
				if err != nil {
					output = err.Error()
				}
				isCorrect, isStrictlyCorrect := CompareOutput(output, tc.ExpectedOutput)
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
	} else {
		// 对解释型语言，直接运行代码
		var wg sync.WaitGroup
		results := make([]model.TestResult, len(testCases))
		wg.Add(len(testCases))
		for i, testCase := range testCases {
			go func(i int, tc model.Case) {
				defer wg.Done()
				output, err := judge.InterpretedJudgeApp.SubmitJudge(language, decodedCode, tc.Input)
				if err != nil {
					output = err.Error()
				}
				isCorrect, isStrictlyCorrect := CompareOutput(output, tc.ExpectedOutput)
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
}
