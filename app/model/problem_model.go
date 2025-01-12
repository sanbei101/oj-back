package model

import (
	"log"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Problem struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement;" json:"problemId,omitempty"`
	Name        string         `gorm:"not null" json:"title"`
	Description string         `gorm:"not null" json:"description"`
	Tags        pq.StringArray `gorm:"type:text[]" json:"tags"`
	TestCase    TestCase       `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE" json:"test_case,omitempty"`
}

type TestCase struct {
	ProblemID uint64 `gorm:"primaryKey" json:"problem_id"`
	Cases     []Case `gorm:"type:jsonb;serializer:json" json:"cases"`
}

type Case struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
}

func InsertData(db *gorm.DB) {
	problems := []Problem{
		{
			Name:        "加法运算",
			Description: "给定两个数字，输出它们的和。",
			Tags:        []string{"数学", "简单"},
		},
		{
			Name:        "字符串拼接",
			Description: "给定两个字符串，返回它们拼接后的结果。",
			Tags:        []string{"字符串", "简单"},
		},
		{
			Name:        "最大数查找",
			Description: "给定一个整数数组，返回其中的最大值。",
			Tags:        []string{"数组", "查找", "简单"},
		},
	}

	if err := db.Create(&problems).Error; err != nil {
		log.Fatal("插入问题数据失败:", err)
	}

	testCases := []TestCase{
		{
			ProblemID: 1,
			Cases: []Case{
				{Input: "1 2", ExpectedOutput: "3"},
				{Input: "3 5", ExpectedOutput: "8"},
				{Input: "10 15", ExpectedOutput: "25"},
			},
		},
		{
			ProblemID: 2,
			Cases: []Case{
				{Input: "hello world", ExpectedOutput: "helloworld"},
				{Input: "foo bar", ExpectedOutput: "foobar"},
				{Input: "abc def", ExpectedOutput: "abcdef"},
			},
		},
		{
			ProblemID: 3,
			Cases: []Case{
				{Input: "1, 2, 3, 4, 5", ExpectedOutput: "5"},
				{Input: "-1, -2, -3, -4", ExpectedOutput: "-1"},
				{Input: "10, 100, 30, 40", ExpectedOutput: "100"},
			},
		},
	}

	if err := db.Create(&testCases).Error; err != nil {
		log.Fatal("插入测试用例数据失败:", err)
	}
}
