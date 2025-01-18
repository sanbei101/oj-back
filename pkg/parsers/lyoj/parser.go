package lyoj

import (
	"encoding/json"
	"fmt"
	"oj-back/app/model"
	"os"
	"path/filepath"
)

// ParseFromFolder 接受一个表示存放LyojProblem的文件夹的路径,解析出对应LyojProblem结构体
func ParseFromFolder(path string) (*LyojProblem, error) {
	metaPath := filepath.Join(path, "meta.json")
	meta := Meta{}
	if err := loadAsStruct(metaPath, &meta); err != nil {
		return nil, fmt.Errorf("无法解析元数据: %v", err)
	}

	descriptionPath := filepath.Join(path, "problem.md")
	desc, err := loadAsString(descriptionPath)
	if err != nil {
		return nil, fmt.Errorf("无法解析问题描述: %v", err)
	}

	configPath := filepath.Join(path, "config.json")
	config := Config{}
	if err := loadAsStruct(configPath, &config); err != nil {
		return nil, fmt.Errorf("无法解析配置: %v", err)
	}

	testCases, err := loadTestCases(path, config)
	if err != nil {
		return nil, fmt.Errorf("无法解析测试用例: %v", err)
	}

	problem := LyojProblem{
		Meta:        meta,
		Description: desc,
		Config:      config,
		TestCases:   testCases,
	}
	return &problem, nil
}

func loadAsStruct(path string, v any) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("无法打开文件: %v", err)
	}
	err = json.Unmarshal(content, v)
	if err != nil {
		return fmt.Errorf("无法解析json: %v", err)
	}
	return nil
}

func loadAsString(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("无法打开文件: %v", err)
	}

	return string(content), nil
}

// 接受一个config结构体和存放测试用例的文件夹路径,读取config中指定的所有测试用例。
func loadTestCases(path string, c Config) ([]model.Case, error) {
	var cases = make([]model.Case, len(c.Datas))
	for i, data := range c.Datas {
		inputPath := filepath.Join(path, data.Input)
		outputPath := filepath.Join(path, data.Output)

		input, err := loadAsString(inputPath)
		if err != nil {
			return nil, fmt.Errorf("无法读取测试用例: %v", err)
		}

		output, err := loadAsString(outputPath)
		if err != nil {
			return nil, fmt.Errorf("无法读取测试用例: %v", err)
		}

		testCase := model.Case{
			Input:          input,
			ExpectedOutput: output,
		}
		cases[i] = testCase
	}
	return cases, nil
}

// ExportToProblem 接受一个LyojProblem和问题ID,构造一个model.Problem
func ExportToProblem(problem *LyojProblem, ID uint64) model.Problem {
	return model.Problem{
		ID:          ID,
		Name:        problem.Meta.Title,
		Description: problem.Description,
		Tags:        problem.Meta.Tags,
		TestCase: &model.TestCase{
			ProblemID: ID,
			Cases:     problem.TestCases,
		},
	}
}
