package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// 运行单一C语言源代码文件
func runCCode(codeFilePath string, input string) (string, error) {
	// 给源码文件加后缀名,否则gcc报错
	err := os.Rename(codeFilePath, codeFilePath+".c")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %v", err)
	}
	codeFilePath = codeFilePath + ".c"

	// 计算出编译目标文件绝对路径
	execPath, err := os.CreateTemp("", "c_out_*.tmp")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %v", err)
	}
	err = execPath.Close()
	if err != nil {
		return "", fmt.Errorf("关闭临时文件失败: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Printf("删除临时文件失败: %v", err)
		}
	}(execPath.Name())

	// 编译
	cmd := exec.Command("gcc", codeFilePath, "-o", execPath.Name())

	var compileStderr bytes.Buffer // 捕获标准错误输出
	cmd.Stderr = &compileStderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("编译失败: %v, %s", err, compileStderr.String())
	}

	// 运行编译后的可执行文件
	runCmd := exec.Command(execPath.Name())

	var stdout, stderr bytes.Buffer
	runCmd.Stdin = strings.NewReader(input)
	runCmd.Stdout = &stdout
	runCmd.Stderr = &stderr

	if err := runCmd.Run(); err != nil {
		return "", fmt.Errorf("运行时错误: %v, %s", err, stderr.String())
	}

	return stdout.String(), nil
}
