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
	"time"
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

	// 使用管道代替文件写入
	r, w, err := os.Pipe()
	if err != nil {
		return "", fmt.Errorf("创建管道失败: %v", err)
	}
	defer r.Close()

	// 启动一个 goroutine 写入代码到管道
	go func() {
		defer w.Close()
		w.Write(decodedCode)
	}()

	// 生成输出文件名并编译代码
	outputFile := fmt.Sprintf("./user_code_out_%d", time.Now().UnixNano())
	cmd := exec.Command("gcc", "-x", "c", "-", "-o", outputFile)
	defer os.Remove(outputFile)
	cmd.Stdin = r

	// 捕获标准错误输出
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// 执行编译命令
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("编译失败: %v, %s", err, stderr.String())
	}

	// 运行编译后的可执行文件
	runCmd := exec.Command(outputFile)
	var out, runStderr bytes.Buffer
	runCmd.Stdin = strings.NewReader(input)
	runCmd.Stdout = &out
	runCmd.Stderr = &runStderr

	// 使用 goroutine 并行执行运行命令
	runErrChan := make(chan error)
	go func() {
		runErrChan <- runCmd.Run()
	}()

	// 等待运行命令完成
	if err := <-runErrChan; err != nil {
		return "", fmt.Errorf("执行代码错误: %v, %s", err, runStderr.String())
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
