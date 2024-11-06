package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"oj-back/app/db"
	"oj-back/app/models"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/panjf2000/ants/v2"
)

var (
	// 设置协程池的最大容量，可以根据您的服务器性能进行调整
	maxPoolSize = 10
	pool        *ants.PoolWithFunc
)

type runCodeTask struct {
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
	if t.language != "c" {
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

	// 编译代码
	outputFile := fmt.Sprintf("./user_code_out_%d", time.Now().UnixNano())
	defer os.Remove(outputFile)

	cmd := exec.Command("gcc", codeFile.Name(), "-o", outputFile)

	// 捕获标准错误输出
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// 执行编译命令
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("编译失败: %v, %s", err, stderr.String())
	}

	// 运行编译后的可执行文件
	runCmd := exec.Command(outputFile)
	var out, runStderr bytes.Buffer
	runCmd.Stdin = strings.NewReader(t.input)
	runCmd.Stdout = &out
	runCmd.Stderr = &runStderr

	// 执行运行命令
	if err := runCmd.Run(); err != nil {
		return "", fmt.Errorf("执行代码错误: %v, %s", err, runStderr.String())
	}

	return out.String(), nil
}

func RunCode(language string, codeContent string, input string) (string, error) {
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

// 比较实际输出与预期输出是否一致
func CompareOutput(actualOutput string, expectedOutput string) bool {
	return strings.TrimSpace(actualOutput) == strings.TrimSpace(expectedOutput)
}

// 从数据库中获取测试用例
func GetTestCases(problemID int) ([]models.Case, error) {
	var record models.TestCase

	err := db.DB.Where("problem_id = ?", problemID).First(&record).Error
	if err != nil {
		return nil, fmt.Errorf("查询测试用例失败: %v", err)
	}
	cases := record.Cases

	return cases, nil
}
