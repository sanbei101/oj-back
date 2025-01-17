package service

import (
	"encoding/base64"
	"oj-back/app/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		assert.NoError(b, err)
	}
}

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
func BenchmarkBasicPy(b *testing.B) {
	var JudgeServiceApp = new(JudgeService)

	// 定义测试用例
	testCases := []model.Case{
		{
			Input:          "1 2",
			ExpectedOutput: "3",
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
	code := `
	a, b = map(int, input().split())
	print(a + b)
	`
	encodedCode := base64.StdEncoding.EncodeToString([]byte(code))
	for i := 0; i < b.N; i++ {
		_, err := JudgeServiceApp.EvaluateProblem("python", encodedCode, testCases)
		assert.NoError(b, err)
	}
}
func TestBasicCpp(t *testing.T) {
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
	#include <iostream>
	using namespace std;
	int main() {
		int a, b;
		cin >> a >> b;
		cout << a + b << endl;
		return 0;
	}
	`
	codeContent = base64.StdEncoding.EncodeToString([]byte(codeContent))
	evaluationResult, err := JudgeServiceApp.EvaluateProblem("cpp", codeContent, testCases)
	assert.NoError(t, err)
	assert.Equal(t, len(testCases), evaluationResult.Count)
	for i, result := range evaluationResult.Results {
		assert.True(t, result.IsSuccess, "测试用例 %d 未通过, 输入: %s, 预期: %s, 实际: %s", i+1, testCases[i].Input, testCases[i].ExpectedOutput, result.ActualOutput)
	}
}
func BenchmarkBasicCpp(b *testing.B) {
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
	#include <iostream>
	using namespace std;
	int main() {
		int a, b;
		cin >> a >> b;
		cout << a + b << endl;
		return 0;
	}
	`
	codeContent = base64.StdEncoding.EncodeToString([]byte(codeContent))
	for i := 0; i < b.N; i++ {
		_, err := JudgeServiceApp.EvaluateProblem("cpp", codeContent, testCases)
		assert.NoError(b, err)
	}
}
func TestBasicPy(t *testing.T) {
	var JudgeServiceApp = new(JudgeService)

	// 定义测试用例
	testCases := []model.Case{
		{
			Input:          "1 2",
			ExpectedOutput: "3",
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
a, b = map(int, input().split())
print(a + b)
`
	codeContent = base64.StdEncoding.EncodeToString([]byte(codeContent))
	evaluationResult, err := JudgeServiceApp.EvaluateProblem("python", codeContent, testCases)
	assert.NoError(t, err)
	assert.Equal(t, len(testCases), evaluationResult.Count)
	for i, result := range evaluationResult.Results {
		assert.True(t, result.IsSuccess, "测试用例 %d 未通过, 输入: %s, 预期: %s, 实际: %s", i+1, testCases[i].Input, testCases[i].ExpectedOutput, result.ActualOutput)
	}
}

// 测试Py模块能够兼容未格式化的期望输出
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

// 测试Py模块能够兼容未格式化的实际
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
