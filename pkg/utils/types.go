package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Question 结构体表示题目配置
type Question struct {
	ID          int
	Title       string
	Description string
}

// TestCase 结构体表示测试用例
type TestCase struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
}

func (t *TestCase) UnmarshalJSON(data []byte) error {
	// 使用匿名结构体解析 ExpectedOutput 和通用的 Input 字段
	aux := struct {
		Input          interface{} `json:"input"`
		ExpectedOutput string      `json:"expected_output"`
	}{}

	// 解析 JSON 数据
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// 根据 Input 的类型进行处理
	switch v := aux.Input.(type) {
	case string:
		t.Input = v // 如果是字符串，直接赋值
	case []interface{}:
		// 将 []interface{} 转为 []string 后拼接
		var strParts []string
		for _, item := range v {
			strParts = append(strParts, fmt.Sprintf("%v", item))
		}
		t.Input = strings.Join(strParts, " ")
	default:
		return fmt.Errorf("未知的 input 类型: %T", v)
	}

	// 赋值 ExpectedOutput
	t.ExpectedOutput = aux.ExpectedOutput
	return nil
}

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
