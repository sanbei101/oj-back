package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/panjf2000/ants/v2"
)

var (
	// 设置协程池的最大容量，可以根据您的服务器性能进行调整
	maxPoolSize = 10
	pool        *ants.PoolWithFunc
	poolExec    *ants.PoolWithFunc
)

type runExecutableTask struct {
	executablePath string
	input          string
	output         string
	err            error
	done           chan bool
}

type runCodeTask struct {
	language    string
	codeContent string
	input       string
	output      string
	err         error
	done        chan bool
}

func init() {
	// 初始化带有返回值的协程池
	pool, _ = ants.NewPoolWithFunc(maxPoolSize, func(task interface{}) {
		t := task.(*runCodeTask)
		// 执行任务
		t.output, t.err = t.run()
		t.done <- true
	})
	poolExec, _ = ants.NewPoolWithFunc(maxPoolSize, func(task interface{}) {
		t := task.(*runExecutableTask)
		// 执行任务
		t.output, t.err = RunExecutable(t.executablePath, t.input)
		t.done <- true
	})
}

// RunExecutable 运行已编译的可执行文件并返回输出
func RunExecutable(executablePath string, input string) (string, error) {
	cmd := exec.Command(executablePath)
	var stdout, stderr bytes.Buffer
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("运行时错误: %v, %s", err, stderr.String())
	}

	return stdout.String(), nil
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
	codeFile, err := os.CreateTemp(".", "user_code_*.py")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %v", err)
	}
	defer os.Remove(codeFile.Name())
	codeFile.Write(decodedCode)
	codeFile.Close()

	if t.language == "c" {
		return runCCode(codeFile, t.input)
	}
	if t.language == "python" {
		return runPyCode(codeFile, t.input)
	}

	return "", fmt.Errorf("不支持的语言")
}
func SubmitExecutableTask(executablePath string, input string) (string, error) {
	task := &runExecutableTask{
		executablePath: executablePath,
		input:          input,
		done:           make(chan bool),
	}

	// 将任务提交到协程池
	err := poolExec.Invoke(task)
	if err != nil {
		return "", fmt.Errorf("无法提交任务到协程池: %v", err)
	}

	// 等待任务完成
	<-task.done

	return task.output, task.err
}

func SubmitCodeExecution(language string, codeContent string, input string) (string, error) {
	task := &runCodeTask{
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
