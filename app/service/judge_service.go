package service

import (
	"encoding/base64"
	"fmt"
	"log"
	"oj-back/app/model"
	"oj-back/pkg/utils"
	"os"
	"sync"
)

type JudgeService struct{}

var JudgeServiceApp = new(JudgeService)

// EvaluateSubmit 该函数接受一个Submit结构体,改变该结构体的evaResult字段。
func (js *JudgeService) EvaluateSubmit(s *model.Submit) {
	// 获取测试用例
	cases, err := ProblemServiceApp.GetProblemTestCase(s.ProblemID)
	if err != nil {
		s.EvaResult = model.EvaluationResult{
			Status: fmt.Errorf("内部错误: %v", err),
		}
	}

	// TODO:统一错误处理

	//解码Base64编码的源码
	decodedCode, err := base64.StdEncoding.DecodeString(s.CodeContent)
	if err != nil {
		s.EvaResult = model.EvaluationResult{
			Status: fmt.Errorf("内部错误: 解码Base64失败:%v", err),
		}
		return
	}

	//创建临时的代码文件
	codeFile, err := os.CreateTemp("", "submit_decoded_*.src")
	if err != nil {
		s.EvaResult = model.EvaluationResult{
			Status: fmt.Errorf("创建临时代码文件失败:%v", err),
		}
		return
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Printf("[WARN]临时文件未正常删除:%v", name)
		}
	}(codeFile.Name())

	_, err = codeFile.Write(decodedCode)
	if err != nil {
		s.EvaResult = model.EvaluationResult{
			Status: fmt.Errorf("内部错误: 无法写入临时代码文件:%v", err),
		}
		return
	}

	err = codeFile.Close()
	if err != nil {
		s.EvaResult = model.EvaluationResult{
			Status: fmt.Errorf("内部错误: 无法写入临时代码文件:%v", err),
		}
		return
	}

	s.EvaResult = *JudgeServiceApp.Evaluate(s.Language, codeFile.Name(), cases)
}

// Evaluate 该函数作为bootstrap,接受一个代码文件,将其复制为多个副本并进行并发测评。
func (js *JudgeService) Evaluate(language string, codeFilePath string, testCases []model.Case) *model.EvaluationResult {
	var wg sync.WaitGroup
	var lock = &sync.RWMutex{}
	results := make([]model.TestResult, len(testCases))

	for i, testCase := range testCases {
		wg.Add(1)
		go func(i int, tc model.Case) {
			defer wg.Done()
			rst, _ := js.judge(language, codeFilePath, lock, testCases[i])
			results[i] = rst
		}(i, testCase)
	}

	wg.Wait()

	evaluation := &model.EvaluationResult{
		Count:   len(results),
		Results: results,
	}

	return evaluation
}

// 判断输出是否符合给定的单个测试用例。该函数不返回error,若出现error则写入TestResult中.
func (js *JudgeService) judge(language string, codeFilePath string, readLock *sync.RWMutex, testCase model.Case) (model.TestResult, error) {
	output, err := utils.RunCode(language, codeFilePath, readLock, testCase.Input)
	if err != nil {
		returnVal := model.TestResult{
			IsSuccess:      false,
			ExpectedOutput: testCase.ExpectedOutput,
			ActualOutput:   err.Error(),
		}
		return returnVal, nil
	}
	isCorrect, isStrictlyCorrect := utils.CompareOutput(output, testCase.ExpectedOutput)
	returnVal := model.TestResult{
		IsSuccess:       isCorrect,
		IsStrictSuccess: isStrictlyCorrect,
		ExpectedOutput:  testCase.ExpectedOutput,
		ActualOutput:    output,
	}
	return returnVal, nil
}
