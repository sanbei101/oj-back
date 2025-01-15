package model

import (
	"log"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Problem struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement;" json:"problemId,omitempty"`
	Name        string         `gorm:"not null;index" json:"title"`
	Description string         `gorm:"not null" json:"description"`
	Tags        pq.StringArray `gorm:"type:text[]" json:"tags"`
	TestCase    *TestCase      `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE" json:"test_case,omitempty"`
}

// TestCase 包含某一问题的所有测试用例的结构体
type TestCase struct {
	ProblemID uint64 `gorm:"primaryKey;index" json:"problem_id"`
	Cases     []Case `gorm:"type:jsonb;serializer:json" json:"cases"`
}

// Case 单一测试用例结构体,包括输入与输出
type Case struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
}

func InsertData(db *gorm.DB) {
	problems := []Problem{
		{
			Name:        "加法运算",
			Description: "**梦开始的地方**\\n```c\\n#include <stdio.h>\\nint main() {\\n    int num1, num2, sum;\\n\\n    scanf(\\\"%d\\\", &num1);\\n    scanf(\\\"%d\\\", &num2);\\n    sum = num1 + num2;\\n    printf(\\\"%d\\n\\\", sum);\\n\\n    return 0;\\n}\\n```",
			Tags:        []string{"数学", "简单"},
		},
		{
			Name:        "字符串拼接",
			Description: "**题目描述**\\n给定两个字符串，返回它们拼接后的结果。\\n\\n**示例**\\n```c\\n#include <stdio.h>\\n#include <string.h>\\n\\nint main() {\\n    char str1[100], str2[100];\\n\\n    scanf(\\\"%s\\\", str1);\\n    scanf(\\\"%s\\\", str2);\\n\\n    printf(\\\"%s%s\\n\\\", str1, str2);\\n\\n    return 0;\\n}\\n```",
			Tags:        []string{"字符串", "简单"},
		},
		{
			Name:        "最大数查找",
			Description: "**题目描述**\\n给定一个整数数组，返回其中的最大值。\\n\\n**示例**\\n```c\\n#include <stdio.h>\\n\\nint main() {\\n    int n, i, max, num;\\n\\n    scanf(\\\"%d\\\", &n); // 输入数组长度\\n\\n    if (n > 0) {\\n        scanf(\\\"%d\\\", &max);\\n        for (i = 1; i < n; i++) {\\n            scanf(\\\"%d\\\", &num);\\n            if (num > max) {\\n                max = num;\\n            }\\n        }\\n    }\\n\\n    printf(\\\"%d\\n\\\", max);\\n\\n    return 0;\\n}\\n```",
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
