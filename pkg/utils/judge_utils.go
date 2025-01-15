package utils

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/panjf2000/ants/v2"
)

var (
	// 设置协程池的最大容量，可以根据您的服务器性能进行调整
	maxPoolSize = 10
	pool        *ants.PoolWithFunc
)

type runCodeTask struct {
	language     string
	codeFilePath string
	readLock     *sync.Mutex
	input        string
	output       string
	err          error
	done         chan bool
}

func init() {
	var err error
	// 初始化一个带有返回值的协程池
	pool, err = ants.NewPoolWithFunc(maxPoolSize, func(task interface{}) {
		t := task.(*runCodeTask)
		// 执行任务
		t.output, t.err = t.run()
		t.done <- true
	})
	if err != nil {
		panic(fmt.Sprintf("无法创建协程池: %v", err))
	}
}

func (t *runCodeTask) run() (string, error) {
	if !(t.language == "c" || t.language == "python") {
		return "", fmt.Errorf("不支持的语言")
	}

	//将传入的文件复制一份,防止并发冲突
	copyOfSrc, err := getCopyOfSrcCode(t.codeFilePath, t.readLock)
	if err != nil {
		return "", fmt.Errorf("创建临时代码文件失败: %v", err)
	}

	if t.language == "c" {
		return runCCode(copyOfSrc, t.input)
	}
	if t.language == "python" {
		return runPyCode(copyOfSrc, t.input)
	}

	return "", fmt.Errorf("不支持的语言")
}

func RunCode(language string, codeFilePath string, readLock *sync.Mutex, input string) (string, error) {
	task := &runCodeTask{
		language:     language,
		codeFilePath: codeFilePath,
		readLock:     readLock,
		input:        input,
		done:         make(chan bool),
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

// 该函数接受一个表示文件路径的字符串和一把读写锁,创建一个该文件的新副本,并返回这个副本的路径.
// 副本文件需要上层函数自行删除.
func getCopyOfSrcCode(srcFilePath string, readLock *sync.Mutex) (string, error) {
	readLock.Lock()
	defer readLock.Unlock()

	// 打开源文件
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return "", err
	}

	// 创建副本文件
	copyOfSrc, err := os.CreateTemp("", "copy_of_src_*")
	if err != nil {
		return "", err
	}

	// 复制
	_, err = io.Copy(copyOfSrc, srcFile)
	if err != nil {
		return "", err
	}

	// 关闭文件
	err = srcFile.Close()
	if err != nil {
		return "", err
	}
	err = copyOfSrc.Close()
	if err != nil {
		return "", err
	}

	return copyOfSrc.Name(), nil
}
