package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// 运行单一C语言源代码文件
func runCCode(codeFile *os.File, input string) (string, error) {
	// 计算出编译目标文件绝对路径
	outputFile := fmt.Sprintf("c_out_%s", uuid.New().String())
	destPath := filepath.Join(os.TempDir(), outputFile)
	defer os.Remove(destPath)

	// 编译
	cmd := exec.Command("gcc", codeFile.Name(), "-o", destPath)

	// 捕获标准错误输出
	var compileStderr bytes.Buffer
	cmd.Stderr = &compileStderr

	// 执行编译命令
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("编译失败: %v, %s", err, compileStderr.String())
	}

	// 运行编译后的可执行文件
	runCmd := exec.Command(destPath)
	var stdout, stderr bytes.Buffer
	runCmd.Stdin = strings.NewReader(input)
	runCmd.Stdout = &stdout
	runCmd.Stderr = &stderr

	// 执行运行命令
	if err := runCmd.Run(); err != nil {
		return "", fmt.Errorf("运行时错误: %v, %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// CompileCCode 编译C代码并返回可执行文件的路径
func CompileCCode(codeContent []byte) (string, error) {
	// 生成唯一的临时文件名
	codeFile, err := os.CreateTemp("", "user_code_*.c")
	if err != nil {
		return "", fmt.Errorf("创建临时C文件失败: %v", err)
	}
	defer os.Remove(codeFile.Name())

	// 写入C代码
	if _, err := codeFile.Write([]byte(codeContent)); err != nil {
		return "", fmt.Errorf("写入C代码失败: %v", err)
	}
	codeFile.Close()

	// 生成可执行文件路径
	executableName := fmt.Sprintf("c_out_%s", uuid.New().String())
	executablePath := filepath.Join(os.TempDir(), executableName)
	defer os.Remove(executablePath)

	// 编译C代码
	cmd := exec.Command("gcc", codeFile.Name(), "-o", executablePath)
	var compileStderr bytes.Buffer
	cmd.Stderr = &compileStderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("编译失败: %v, %s", err, compileStderr.String())
	}

	return executablePath, nil
}
