package service

import (
	"oj-back/internal/db"
	"oj-back/internal/db/models"
)

func GetAllProblems() ([]models.Problems, error) {
	var problems []models.Problems

	err := db.DB.Find(&problems).Error
	if err != nil {
		return nil, err
	}

	return problems, nil
}

func GetProblemByID(id int) (*models.Problems, error) {
	var problem models.Problems

	err := db.DB.Where("id = ?", id).First(&problem).Error
	if err != nil {
		return nil, err
	}

	return &problem, nil
}
