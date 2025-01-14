package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

// It ensures that the program correctly processes input, produces the expected result, and cleans up temporary files.
func TestBasic(t *testing.T) {
	codeContent := `
	#include <stdio.h>
	int main() {
		int a, b;
		scanf("%d %d", &a, &b);
		printf("%d\n", a + b);
		return 0;
	}
	`

	// 创建临时代码文件
	codeFile, err := os.CreateTemp("", "test_*.c")
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

	// 执行
	output, err := runCCode(codeFile, "10 20", 0)
	expected := "30"

	// 检查输出
	assert.NoError(t, err)
	same, _ := CompareOutput(output, expected)
	assert.True(t, same)

	// 检测编译后的可执行文件是否被删除
	assert.NoFileExists(t, path.Join(os.TempDir(), "user_code_out_0"))
}
