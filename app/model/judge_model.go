package model

// 评测结果结构体
type TestResult struct {
	IsSuccess       bool   `json:"is_success"`
	IsStrictSuccess bool   `json:"is_strict_success"`
	ExpectedOutput  string `json:"expected_output"`
	ActualOutput    string `json:"actual_output"`
}

type EvaluationResult struct {
	Count   int          `json:"count"`
	Results []TestResult `json:"results"`
}
