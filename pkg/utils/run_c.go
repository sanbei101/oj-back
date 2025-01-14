package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 运行单一C语言源代码文件
func runCCode(codeFile *os.File, input string, threadID int) (string, error) {
	// 计算出编译目标文件绝对路径
	outputFile := fmt.Sprintf("user_code_out_%d", threadID)
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
