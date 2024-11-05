package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"oj-back/internal/db"
	"oj-back/internal/db/models"
	"os"
	"os/exec"
	"strings"
)

// 编译并运行 C 代码字符串，并传入测试输入，返回代码的输出结果
func RunCode(language string, codeContent string, input string) (string, error) {
	if language != "c" {
		return "", fmt.Errorf("不支持的语言")
	}

	// 创建临时文件保存 C 代码
	tmpFile, err := os.CreateTemp("", "user_code_*.c")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// 将代码内容写入临时文件
	if _, err := tmpFile.WriteString(codeContent); err != nil {
		return "", fmt.Errorf("写入代码到临时文件失败: %v", err)
	}
	tmpFile.Close()

	// 编译 C 代码
	outputFile := tmpFile.Name() + "_out"
	defer os.Remove(outputFile)

	if err := exec.Command("gcc", tmpFile.Name(), "-o", outputFile).Run(); err != nil {
		return "", fmt.Errorf("编译失败: %v", err)
	}

	// 运行编译后的可执行文件
	cmd := exec.Command(outputFile)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("执行代码错误: %v, %s", err, stderr.String())
	}

	return out.String(), nil
}

// 比较实际输出与预期输出是否一致
func CompareOutput(actualOutput string, expectedOutput string) bool {
	return strings.TrimSpace(actualOutput) == strings.TrimSpace(expectedOutput)
}

// 从数据库中获取测试用例
func GetTestCases(problemID int) ([]TestCase, error) {
	var record models.TestCases

	// 使用 GORM 查询
	err := db.DB.Where("problem_id = ?", problemID).First(&record).Error
	if err != nil {
		return nil, fmt.Errorf("查询测试用例失败: %v", err)
	}

	// 将 JSON 字符串解析为 TestCase 切片
	var testCases []TestCase
	err = json.Unmarshal([]byte(record.Cases), &testCases)
	if err != nil {
		return nil, fmt.Errorf("解析测试用例 JSON 失败: %v", err)
	}

	return testCases, nil
}

// 评测函数，循环遍历每个测试用例并进行评测
func EvaluateProblem(language string, codeContent string, testCases []TestCase) (*EvaluationResult, error) {
	var results []TestResult

	for _, testCase := range testCases {
		// 执行用户代码并获取输出
		output, err := RunCode(language, codeContent, testCase.Input)
		if err != nil {
			return nil, err
		}

		// 比对输出
		isCorrect := CompareOutput(output, testCase.ExpectedOutput)

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
