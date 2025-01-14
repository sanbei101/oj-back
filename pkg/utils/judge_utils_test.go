package utils

import (
	"github.com/stretchr/testify/assert"
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
