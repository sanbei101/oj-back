package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"oj-back/app/db"
	"oj-back/app/models"
	"os"
	"os/exec"
	"strings"
)

// 编译并运行 C 代码字符串，并传入测试输入，返回代码的输出结果
func RunCode(language string, codeContent string, input string) (string, error) {
	if language != "c" {
		return "", fmt.Errorf("不支持的语言")
	}

	decodedCode, err := base64.StdEncoding.DecodeString(codeContent)
	if err != nil {
		return "", fmt.Errorf("解码代码失败: %v", err)
	}

	// 创建临时文件保存 C 代码
	tmpFile, err := os.CreateTemp("", "user_code_*.c")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// 将代码内容写入临时文件
	if _, err := tmpFile.WriteString(string(decodedCode)); err != nil {
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
	var out, stderr bytes.Buffer
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("执行代码错误: %v, %s", err, stderr.String())
	}

	return out.String(), nil
}

// 比较实际输出与预期输出是否一致
func CompareOutput(actualOutput string, expectedOutput string) bool {
	return strings.TrimSpace(actualOutput) == strings.TrimSpace(expectedOutput)
}

// 从数据库中获取测试用例
func GetTestCases(problemID int) ([]models.Case, error) {
	var record models.TestCase

	err := db.DB.Where("problem_id = ?", problemID).First(&record).Error
	if err != nil {
		return nil, fmt.Errorf("查询测试用例失败: %v", err)
	}
	cases := record.Cases

	return cases, nil
}
