package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// 运行单一Python源代码文件
func runPyCode(codeFilePath string, input string) (string, error) {
	cmd := exec.Command("python", codeFilePath)

	var stdout, stderr bytes.Buffer
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("运行时错误: %v, %s", err, stderr.String())
	}

	return stdout.String(), nil
}
