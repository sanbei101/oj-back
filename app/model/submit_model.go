package model

import "github.com/google/uuid"

// Submit 表示用户的一次提交
type Submit struct {
	// TODO:还可以增加提交时间等字段
	SubmitID    uuid.UUID        `json:"submit_id"`
	UserID      uuid.UUID        `json:"user_id"`
	ProblemID   uint64           `json:"problem_id"`
	Language    string           `json:"language"`
	CodeContent string           `json:"code_content"`
	EvaResult   EvaluationResult `json:"eva_result"`
}

// TestResult 表示某单一测试用例的测试结果
type TestResult struct {
	IsSuccess       bool   `json:"is_success"`
	IsStrictSuccess bool   `json:"is_strict_success"`
	ExpectedOutput  string `json:"expected_output"`
	ActualOutput    string `json:"actual_output"`
}

// EvaluationResult 包含某一提交的所有测试用例的测试结果
type EvaluationResult struct {
	Count   int          `json:"count"`
	Results []TestResult `json:"results"`
	Status  error        `json:"status"` // 评估结果,包括内部错误、编译错误、错误答案等等
	//Status字段好像没有很完善,后序需要再修改
}
