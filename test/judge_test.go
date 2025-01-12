package test

import (
	"encoding/base64"
	"oj-back/app/model"
	"oj-back/app/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 EvaluateProblem 函数性能
func BenchmarkEvaluateProblem(b *testing.B) {
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
		_, err := service.JudgeServiceApp.EvaluateProblem("c", encodedCodeContent, testCases)
		if err != nil {
			b.Fatalf("评测失败: %v", err)
		}
	}
}
func TestEvaluateProblem(t *testing.T) {
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
	//base64编码
	codeContent = base64.StdEncoding.EncodeToString([]byte(codeContent))
	evaluationResult, err := service.JudgeServiceApp.EvaluateProblem("c", codeContent, testCases)
	assert.NoError(t, err)
	assert.Equal(t, len(testCases), evaluationResult.Count)
	for i, result := range evaluationResult.Results {
		assert.True(t, result.IsSuccess, "测试用例 %d 未通过, 输入: %s, 预期: %s, 实际: %s", i+1, testCases[i].Input, testCases[i].ExpectedOutput, result.ActualOutput)
	}
}
