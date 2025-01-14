package utils

import (
	"encoding/base64"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"os"
	"strings"
)

var (
	// 设置协程池的最大容量，可以根据您的服务器性能进行调整
	maxPoolSize = 10
	pool        *ants.PoolWithFunc
)

type runCodeTask struct {
	threadID    int
	language    string
	codeContent string
	input       string
	output      string
	err         error
	done        chan bool
}

func init() {
	var err error
	// 初始化一个带有返回值的协程池
	pool, err = ants.NewPoolWithFunc(maxPoolSize, func(task interface{}) {
		t := task.(*runCodeTask)
		t.output, t.err = t.run()
		t.done <- true
	})
	if err != nil {
		panic(fmt.Sprintf("无法创建协程池: %v", err))
	}
	fmt.Println("协程池初始化成功")
}

func (t *runCodeTask) run() (string, error) {
	if !(t.language == "c" || t.language == "python") {
		return "", fmt.Errorf("不支持的语言")
	}

	decodedCode, err := base64.StdEncoding.DecodeString(t.codeContent)
	if err != nil {
		return "", fmt.Errorf("解码代码失败: %v", err)
	}

	// 将解码后的代码写入临时文件
	codeFile, err := os.CreateTemp("", "user_code_*.c")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %v", err)
	}
	defer os.Remove(codeFile.Name())
	codeFile.Write(decodedCode)
	codeFile.Close()

	if t.language == "c" {
		return runCCode(codeFile, t.input, t.threadID)
	}
	if t.language == "python" {
		return runPyCode(codeFile, t.input)
	}

	return "", fmt.Errorf("不支持的语言")
}

func RunCode(ThreadID int, language string, codeContent string, input string) (string, error) {
	task := &runCodeTask{
		threadID:    ThreadID,
		language:    language,
		codeContent: codeContent,
		input:       input,
		done:        make(chan bool),
	}

	// 将任务提交到协程池
	err := pool.Invoke(task)
	if err != nil {
		return "", fmt.Errorf("无法提交任务到协程池: %v", err)
	}

	// 等待任务完成
	<-task.done

	return task.output, task.err
}

// CompareOutput 比较实际输出与预期输出是否一致,
// 第二个返回值表示实际输出是否需要修剪空格才能与预期值一致.
func CompareOutput(actualOutput string, expectedOutput string) (same bool, strictlySame bool) {
	// TODO: 似乎逻辑有点不足,例如每一行行尾的/r/n或/n无法区分,需要改进一下.
	actual := strings.TrimSpace(actualOutput)
	expected := strings.TrimSpace(expectedOutput)
	return actual == expected, expectedOutput == actualOutput
}
