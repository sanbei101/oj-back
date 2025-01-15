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

type TestCase struct {
	ProblemID uint64 `gorm:"primaryKey;index" json:"problem_id"`
	Cases     []Case `gorm:"type:jsonb;serializer:json" json:"cases"`
}

type Case struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
}

func InsertData(db *gorm.DB) {
	problems := []Problem{
		{
			Name: "加法运算",
			Description: "**题目描述**\n" +
				"```c\n" +
				"#include <stdio.h>\n" +
				"int main() {\n" +
				"    int num1, num2, sum;\n" +
				"    scanf(\"%d\", &num1);\n" +
				"    scanf(\"%d\", &num2);\n" +
				"    sum = num1 + num2;\n" +
				"    printf(\"%d\\n\", sum);\n" +
				"    return 0;\n" +
				"}\n" +
				"```",
			Tags: []string{"数学", "简单"},
		},
		{
			Name: "冒泡排序",
			Description: "**题目描述**\n" +
				"```c\n" +
				"#include <stdio.h>\n" +
				"void bubbleSort(int arr[], int n) {\n" +
				"    int i, j, temp;\n" +
				"    for (i = 0; i < n-1; i++) {\n" +
				"        for (j = 0; j < n-i-1; j++) {\n" +
				"            if (arr[j] > arr[j+1]) {\n" +
				"                temp = arr[j];\n" +
				"                arr[j] = arr[j+1];\n" +
				"                arr[j+1] = temp;\n" +
				"            }\n" +
				"        }\n" +
				"    }\n" +
				"}\n" +
				"int main() {\n" +
				"    int arr[] = {64, 34, 25, 12, 22, 11, 90};\n" +
				"    int n = sizeof(arr)/sizeof(arr[0]);\n" +
				"    bubbleSort(arr, n);\n" +
				"    printf(\"Sorted array: \\n\");\n" +
				"    for (int i=0; i < n; i++)\n" +
				"        printf(\"%d \", arr[i]);\n" +
				"    printf(\"\\n\");\n" +
				"    return 0;\n" +
				"}\n" +
				"```",
			Tags: []string{"排序", "冒泡排序"},
		},
		{
			Name: "斐波那契数列",
			Description: "**题目描述**\n" +
				"```c\n" +
				"#include <stdio.h>\n" +
				"int fibonacci(int n) {\n" +
				"    if (n <= 1) return n;\n" +
				"    return fibonacci(n-1) + fibonacci(n-2);\n" +
				"}\n" +
				"int main() {\n" +
				"    int n;\n" +
				"    scanf(\"%d\", &n);\n" +
				"    printf(\"%d\\n\", fibonacci(n));\n" +
				"    return 0;\n" +
				"}\n" +
				"```",
			Tags: []string{"数学", "递归"},
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
				{Input: "64 34 25 12 22 11 90", ExpectedOutput: "11 12 22 25 34 64 90"},
				{Input: "5 1 4 2 8", ExpectedOutput: "1 2 4 5 8"},
				{Input: "3 0 2 5 -1 4 1", ExpectedOutput: "-1 0 1 2 3 4 5"},
			},
		},
		{
			ProblemID: 3,
			Cases: []Case{
				{Input: "5", ExpectedOutput: "5"},
				{Input: "10", ExpectedOutput: "55"},
				{Input: "15", ExpectedOutput: "610"},
			},
		},
	}

	if err := db.Create(&testCases).Error; err != nil {
		log.Fatal("插入测试用例数据失败:", err)
	}
}
