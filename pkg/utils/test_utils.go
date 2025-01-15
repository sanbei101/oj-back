package utils

import (
	"os"
	"testing"
)

// 存放一些专门用于单元测试的函数

// TmpFile 接受一个字符串,返回一个内容为该字符串的临时文件。该文件需要调用者手动删除。
func TmpFile(t *testing.T, content string) string {

	codeFile, _ := os.CreateTemp("", "test_*.tmp")

	_, err := codeFile.Write([]byte(content))
	if err != nil {
		t.Fatalf("创建临时文件失败:%v", err)
	}

	err = codeFile.Close()
	if err != nil {
		t.Fatalf("创建临时文件失败:%v", err)
	}

	return codeFile.Name()
}
