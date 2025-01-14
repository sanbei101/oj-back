package service

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"oj-back/app/model"
	"testing"
)

// BenchmarkBasicC benchmarks the performance of the JudgeService's EvaluateProblem method for a basic C program that adds two integers. It tests the method's efficiency by running multiple iterations with predefined test cases involving integer addition. The benchmark measures the time and resource consumption of evaluating a simple C code snippet across different input scenarios.
func BenchmarkBasicC(b *testing.B) {
	var JudgeServiceApp = new(JudgeService)

	// 定义测试用例
	testCases := []model.Case{
		{
			Input:          "2 3",
			ExpectedOutput: "5",
		},
		{
			Input:          "10 20",
			ExpectedOutput: "30",
		},
		{
			Input:          "-5 5",
			ExpectedOutput: "0",
		},
	}

	// 要测试的代码内容
	codeContent := `
	#include <stdio.h>
	int main() {
		int a, b;
		scanf("%d %d", &a, &b);
		printf("%d\n", a + b);
		return 0;
	}
	`
	// base64 编码
	encodedCodeContent := base64.StdEncoding.EncodeToString([]byte(codeContent))

	for i := 0; i < b.N; i++ {
		_, err := JudgeServiceApp.EvaluateProblem("c", encodedCodeContent, testCases)
		if err != nil {
			b.Fatalf("评测失败: %v", err)
		}
	}
}

// TestBasicC tests the JudgeService's ability to evaluate a basic C program that adds two integers. It verifies the correct handling of various input scenarios, including positive, negative, and zero values. The test checks that the EvaluateProblem method correctly processes the code, runs it against multiple test cases, and returns successful results for each input combination.
func TestBasicC(t *testing.T) {
	var JudgeServiceApp = new(JudgeService)

	// 定义测试用例
	testCases := []model.Case{
		{
			Input:          "2 3",
			ExpectedOutput: "5",
		},
		{
			Input:          "10 20",
			ExpectedOutput: "30",
		},
		{
			Input:          "-5 5",
			ExpectedOutput: "0",
		},
	}

	// 测试的代码内容
	codeContent := `
	#include <stdio.h>
	int main() {
		int a, b;
		scanf("%d %d", &a, &b);
		printf("%d\n", a + b);
		return 0;
	}
	`
	// base64编码
	codeContent = base64.StdEncoding.EncodeToString([]byte(codeContent))
	evaluationResult, err := JudgeServiceApp.EvaluateProblem("c", codeContent, testCases)
	assert.NoError(t, err)
	assert.Equal(t, len(testCases), evaluationResult.Count)
	for i, result := range evaluationResult.Results {
		assert.True(t, result.IsSuccess, "测试用例 %d 未通过, 输入: %s, 预期: %s, 实际: %s", i+1, testCases[i].Input, testCases[i].ExpectedOutput, result.ActualOutput)
	}
}

// TestBadlyFormattedExpectedOut tests the Python module's compatibility with unformatted expected outputs. It verifies that the JudgeService can handle various output formatting scenarios, including different line endings and whitespace characters.
func TestBadlyFormattedExpectedOut(t *testing.T) {
	var JudgeServiceApp = new(JudgeService)

	testCases := []model.Case{
		{
			Input:          "",
			ExpectedOutput: "Hello, World!", //正确输出
		},
		{
			Input:          "",
			ExpectedOutput: "Hello, World!\r\n", //测试CRLF
		},
		{
			Input:          "",
			ExpectedOutput: "Hello, World!\r\n", //测试LF
		},
		{
			Input:          "",
			ExpectedOutput: "\n\t\rHello, World!\n\t\r", //测试随机空格字符
		},
	}

	codeContent := `print("Hello, World!")`
	codeContent = base64.StdEncoding.EncodeToString([]byte(codeContent))
	evaluationResult, err := JudgeServiceApp.EvaluateProblem("python", codeContent, testCases)
	assert.NoError(t, err)
	assert.Equal(t, len(testCases), evaluationResult.Count)
	for i, result := range evaluationResult.Results {
		assert.True(t, result.IsSuccess, "测试用例 %d 未通过, 输入: %s, 预期: %s, 实际: %s", i+1, testCases[i].Input, testCases[i].ExpectedOutput, result.ActualOutput)
	}
}

// TestBadlyFormattedActualOut tests the Python module's compatibility with various output formatting scenarios. It verifies that the JudgeService can correctly evaluate Python code snippets with different whitespace and line ending variations while maintaining the expected output.
func TestBadlyFormattedActualOut(t *testing.T) {
	var JudgeServiceApp = new(JudgeService)

	testCases := []model.Case{
		{
			Input:          "",
			ExpectedOutput: "Hello, World!", //正确输出
		},
	}

	codeContents := []string{
		`print("Hello, World!")`,
		`print("Hello, World!\r\n")`,
		`print("Hello, World!\n")`,
		`print("\tHello, World!\t")`,
		`print("\tHello, World!\r\n\t")`,
	}

	for i, codeContent := range codeContents {
		codeContent = base64.StdEncoding.EncodeToString([]byte(codeContent))
		evaluationResult, err := JudgeServiceApp.EvaluateProblem("python", codeContent, testCases)
		assert.NoError(t, err)
		assert.Equal(t, len(testCases), evaluationResult.Count)
		for _, result := range evaluationResult.Results {
			assert.True(t, result.IsSuccess, "测试用例 %d 未通过, 源代码: %s, 预期: %s, 实际: %s", i+1, codeContents[i], testCases[0].ExpectedOutput, result.ActualOutput)
		}
	}
}
