package types

type TestCases struct {
	ProblemID int    `gorm:"column:problem_id"`
	Cases     string `gorm:"column:cases"` // 注意 Cases 存储的是 JSON 字符串
}
