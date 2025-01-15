package service

import (
	"oj-back/app/model"
	"oj-back/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicC(t *testing.T) {
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

	codeContent := `
	#include <stdio.h>
	int main() {
		int a, b;
		scanf("%d %d", &a, &b);
		printf("%d\n", a + b);
		return 0;
	}
	`
	testAllSuccess(t, codeContent, "c", testCases, 2)
}

// 测试Py模块能够兼容未格式化的期望输出
func TestBadlyFormattedExpectedOut(t *testing.T) {

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

	testAllSuccess(t, codeContent, "python", testCases, 2)
}

// 测试Py模块能够兼容未格式化的实际
func TestBadlyFormattedActualOut(t *testing.T) {
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

	for _, codeContent := range codeContents {
		testAllSuccess(t, codeContent, "python", testCases, 2)
	}
}

// 辅助函数,用来测试给定代码在一组测试用例下能否全部通过
func testAllSuccess(t *testing.T, code string, language string, testCases []model.Case, repeat int) {
	for range repeat {
		codeFilePath := utils.TmpFile(t, code)
		evaluationResult := JudgeServiceApp.Evaluate(language, codeFilePath, testCases)
		assert.Equal(t, len(testCases), evaluationResult.Count)
		for i, result := range evaluationResult.Results {
			assert.True(t, result.IsSuccess,
				"测试用例 %d 未通过, 输入: %s, 预期: %s, 实际: %s",
				i+1, testCases[i].Input, testCases[i].ExpectedOutput, result.ActualOutput)
		}
	}
}
