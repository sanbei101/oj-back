package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// 检测Py运行模块是否能正常执行
func TestHelloWorld(t *testing.T) {
	codeContent := `print("Hello World!")`

	//创建临时代码文件
	codeFile, err := os.CreateTemp("", "temp.py")
	if err != nil {
		t.Error(err)
	}
	_, err = codeFile.Write([]byte(codeContent))
	if err != nil {
		return
	}
	err = codeFile.Close()
	if err != nil {
		return
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Error(err)
		}
	}(codeFile.Name())

	//运行
	output, err := runPyCode(codeFile, "")
	assert.NoError(t, err)
	same, _ := CompareOutput(output, "Hello World!")
	assert.True(t, same)
}

// 检测Py运行模块是否能执行带import的Python代码
func TestImport(t *testing.T) {
	codeContent := `import thisIsANonExistingPackage`

	codeFile, err := os.CreateTemp("", "test_*.py")
	if err != nil {
		t.Error(err)
	}
	_, err = codeFile.Write([]byte(codeContent))
	if err != nil {
		return
	}
	err = codeFile.Close()
	if err != nil {
		return
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Error(err)
		}
	}(codeFile.Name())

	_, err = runPyCode(codeFile, "")
	assert.Error(t, err)
}
