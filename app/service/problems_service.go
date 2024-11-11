package service

import (
	"fmt"
	"oj-back/app/db"
	"oj-back/app/models"

	"github.com/lib/pq"
)

type Page[T any] struct {
	Total int64 `json:"total"`
	Data  []T   `json:"data"`
}
type ProblemDTO struct {
	ID   uint           `json:"problemId"`
	Name string         `json:"title"`
	Tags pq.StringArray `gorm:"type:text[];" json:"tags,omitempty"`
}
type ProblemDetailDTO struct {
	ID          uint           `json:"problemId"`
	Name        string         `json:"title"`
	Description string         `json:"description"`
	Tags        pq.StringArray `gorm:"type:text[];" json:"tags,omitempty"`
}

func GetAllProblems(page int, size int, keyword string) (*Page[ProblemDTO], error) {
	var problems []ProblemDTO
	var total int64
	query := db.DB.Model(&models.Problem{})
	if keyword == "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("获取题目总数失败: %w", err)
	}
	err := query.Offset((page - 1) * size).Limit(size).Find(&problems).Error
	if err != nil {
		return nil, fmt.Errorf("获取题目列表失败: %w", err)
	}

	return &Page[ProblemDTO]{
		Total: total,
		Data:  problems,
	}, nil

}

func GetProblemByID(id int) (*ProblemDetailDTO, error) {
	var problem ProblemDetailDTO
	err := db.DB.Model(&models.Problem{}).Where("id = ?", id).First(&problem).Error
	if err != nil {
		return nil, fmt.Errorf("获取题目详情失败: %w", err)
	}

	return &problem, nil
}
