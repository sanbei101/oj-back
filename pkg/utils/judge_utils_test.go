package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
)

func TestCompareOutput(t *testing.T) {
	//测试不匹配的预期输入和实际输入
	same, strictlySame := CompareOutput("Hello,world!", "Hello,word!")
	assert.False(t, same)
	assert.False(t, strictlySame)

	//测试错误格式的预期输出
	same, strictlySame = CompareOutput("Hello,world!", "Hello,world!")
	assert.True(t, same)
	assert.True(t, strictlySame)

	same, strictlySame = CompareOutput("Hello,world!", "Hello,world!\n")
	assert.True(t, same)
	assert.False(t, strictlySame)

	same, strictlySame = CompareOutput("Hello,world!", "Hello,world!\t")
	assert.True(t, same)
	assert.False(t, strictlySame)

	same, strictlySame = CompareOutput("Hello,world!", "Hello,world!\r\n")
	assert.True(t, same)
	assert.False(t, strictlySame)

	same, strictlySame = CompareOutput("Hello,world!", "\r\nHello,world!\r\n")
	assert.True(t, same)
	assert.False(t, strictlySame)

	same, strictlySame = CompareOutput("Hello,world!", "\n\t\rHello,world!\n\t\r")
	assert.True(t, same)
	assert.False(t, strictlySame)

	//测试错误格式的实际输出
	same, strictlySame = CompareOutput("Hello,world!", "Hello,world!")
	assert.True(t, same)
	assert.True(t, strictlySame)

	same, strictlySame = CompareOutput("Hello,world!\n", "Hello,world!")
	assert.True(t, same)
	assert.False(t, strictlySame)

	same, strictlySame = CompareOutput("Hello,world!\t", "Hello,world!")
	assert.True(t, same)
	assert.False(t, strictlySame)

	same, strictlySame = CompareOutput("Hello,world!\r\n", "Hello,world!")
	assert.True(t, same)
	assert.False(t, strictlySame)

	same, strictlySame = CompareOutput("\r\nHello,world!\r\n", "Hello,world!")
	assert.True(t, same)
	assert.False(t, strictlySame)

	same, strictlySame = CompareOutput("\n\t\rHello,world!\n\t\r", "Hello,world!")
	assert.True(t, same)
	assert.False(t, strictlySame)
}

func TestGetCopyOfSrcCode(t *testing.T) {
	// 创建要复制的文件
	temp, err := os.CreateTemp("", "ThisIsASrcCode")
	if err != nil {
		t.Fatal(err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatal(err)
		}
	}(temp.Name())
	defer func(temp *os.File) {
		err := temp.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(temp)

	// 写入
	content := []byte("Hello world!")
	err = os.WriteFile(temp.Name(), content, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 测试
	code, err := getCopyOfSrcCode(temp.Name(), &sync.Mutex{})
	assert.NoError(t, err)
	assert.NotNil(t, code)

	actualContent, err := os.ReadFile(code)
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatal(err)
		}
	}(code)
	assert.NoError(t, err)
	assert.Equal(t, content, actualContent)
}
