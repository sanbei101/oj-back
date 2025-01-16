package judge

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/panjf2000/ants/v2"
)

var (
	InterpretedJudgePool *ants.PoolWithFunc
)

type InterpretedJudge struct{}

var InterpretedJudgeApp = new(InterpretedJudge)

type InterpretedJudgeTask struct {
	language string
	code     []byte
	input    string
	output   string
	err      error
	done     chan bool
}

func init() {
	InterpretedJudgePool, _ = ants.NewPoolWithFunc(maxPoolSize, func(task interface{}) {
		t := task.(*InterpretedJudgeTask)
		// 执行任务
		t.output, t.err = t.run()
		t.done <- true
	})
}
func (t *InterpretedJudgeTask) run() (string, error) {
	codeFile, err := os.CreateTemp(".", "user_code_*.py")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %v", err)
	}
	defer os.Remove(codeFile.Name())
	codeFile.Write(t.code)
	codeFile.Close()

	if t.language == "python" {
		return InterpretedJudgeApp.runPyCode(codeFile, t.input)
	}

	return "", fmt.Errorf("不支持的语言")
}

func (j *InterpretedJudge) SubmitJudge(language string, code []byte, input string) (string, error) {
	task := &InterpretedJudgeTask{
		language: language,
		code:     code,
		input:    input,
		done:     make(chan bool),
	}

	// 将任务提交到协程池
	err := InterpretedJudgePool.Invoke(task)
	if err != nil {
		return "", fmt.Errorf("无法提交任务到协程池: %v", err)
	}

	// 等待任务完成
	<-task.done

	return task.output, task.err
}

// 运行单一Python源代码文件
func (j *InterpretedJudge) runPyCode(codeFile *os.File, input string) (string, error) {
	cmd := exec.Command("python3", codeFile.Name())

	var stdout, stderr bytes.Buffer
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("运行时错误: %v, %s", err, stderr.String())
	}

	return stdout.String(), nil
}
