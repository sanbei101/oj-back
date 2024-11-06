package service

import (
	"oj-back/app/db"
	"oj-back/app/models"

	"github.com/lib/pq"
)

type ProblemDTO struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Tags        pq.StringArray `gorm:"type:text[];" json:"tags,omitempty"`
}

func GetAllProblems() ([]ProblemDTO, error) {
	var problems []ProblemDTO

	err := db.DB.Model(models.Problem{}).Find(&problems).Error
	if err != nil {
		return nil, err
	}

	return problems, nil
}

func GetProblemByID(id int) (*models.Problem, error) {
	var problem models.Problem

	err := db.DB.Where("id = ?", id).First(&problem).Error
	if err != nil {
		return nil, err
	}

	return &problem, nil
}
