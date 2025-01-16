package model

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
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
	fileContent, err := os.ReadFile("data/1.toml")
	if err != nil {
		log.Fatalf("Failed to read TOML file: %v", err)
	}

	// 修改数据结构定义
	type TomlProblem struct {
		Name        string   `toml:"name"`
		Description string   `toml:"description"`
		Tags        []string `toml:"tags"`
		Test        []struct {
			Input  string `toml:"input"`
			Output string `toml:"output"`
		} `toml:"test"`
	}
	type TomlFile struct {
		TomlProblems []TomlProblem `toml:"problem"`
	}
	var problems TomlFile

	err = toml.Unmarshal(fileContent, &problems)
	if err != nil {
		log.Fatalf("Failed to parse TOML file: %v", err)
	}
	fmt.Println(problems.TomlProblems[0].Name)
	for _, problem := range problems.TomlProblems {
		var InsertProblem Problem
		InsertProblem.Name = problem.Name
		InsertProblem.Description = problem.Description
		InsertProblem.Tags = problem.Tags
		if err := db.Create(&InsertProblem).Error; err != nil {
			log.Fatalf("Failed to insert problem: %v", err)
		}
		var InsertTestCase TestCase
		InsertTestCase.ProblemID = InsertProblem.ID
		for _, test := range problem.Test {
			InsertTestCase.Cases = append(InsertTestCase.Cases, Case{
				Input:          test.Input,
				ExpectedOutput: test.Output,
			})
		}
		if err := db.Create(&InsertTestCase).Error; err != nil {
			log.Fatalf("Failed to insert test case: %v", err)
		}
	}
}
