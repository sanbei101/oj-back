package service

import (
	"oj-back/app/db"
	"oj-back/app/models"
)

func GetAllProblems() ([]models.Problem, error) {
	var problems []models.Problem

	err := db.DB.Find(&problems).Error
	if err != nil {
		return nil, err
	}

	return problems, nil
}

func GetAllProblemsExceptDesc() ([]models.Problem, error) {
	var problems []models.Problem

	err := db.DB.Select("id", "name", "tags").Find(&problems).Error
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
