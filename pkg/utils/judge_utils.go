package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"oj-back/app/db"
	"oj-back/app/models"
	"os"
	"os/exec"
	"strings"
)

// 题目配置
type Question struct {
	ID          int
	Title       string
	Description string
}

// 测试用例
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
	var record models.TestCase

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
