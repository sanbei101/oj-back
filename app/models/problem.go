package models

import (
	"log"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Problem struct {
	ID          uint           `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `gorm:"not null" json:"description,omitempty"`
	Tags        pq.StringArray `gorm:"type:text[];" json:"tags,omitempty"`
	TestCases   TestCase       `gorm:"foreignKey:ProblemID" json:"test_cases,omitempty"`
}

type TestCase struct {
	ProblemID uint   `gorm:"primaryKey;not null;unique;" json:"problem_id"`
	Cases     string `gorm:"type:jsonb;not null" json:"cases"`
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
			Cases:     `[{"input": "1 2", "expected_output": "3"}, {"input": "3 5", "expected_output": "8"}, {"input": "10 15", "expected_output": "25"}]`,
		},
		{
			ProblemID: 2,
			Cases:     `[{"input": ["hello", "world"], "expected_output": "helloworld"}, {"input": ["foo", "bar"], "expected_output": "foobar"}, {"input": ["abc", "def"], "expected_output": "abcdef"}]`,
		},
		{
			ProblemID: 3,
			Cases:     `[{"input": "[1, 2, 3, 4, 5]", "expected_output": "5"}, {"input": "[-1, -2, -3, -4]", "expected_output": "-1"}, {"input": "[10, 100, 30, 40]", "expected_output": "100"}]`,
		},
	}

	if err := db.Create(&testCases).Error; err != nil {
		log.Fatal("插入测试用例数据失败:", err)
	}
}
