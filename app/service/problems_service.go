package service

import (
	"fmt"
	"oj-back/app/db"
	"oj-back/app/model"
)

type ProblemService struct{}

var ProblemServiceApp = new(ProblemService)

// GetAllProblems 查询所有题目并分页
func (ps *ProblemService) GetAllProblems(page int, size int, keyword string) (*model.Page[model.Problem], error) {
	var problems []model.Problem
	var total int64
	query := db.DB.Model(&model.Problem{})
	// 如果关键词不为空，进行匹配
	if keyword != "" {
		// 模糊匹配 name 和 tags
		query = query.Where("name ILIKE ?", "%"+keyword+"%").
			Or("tags ILIKE ?", "%"+keyword+"%") // 使用 ILIKE 进行大小写不敏感的匹配
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("获取题目总数失败: %w", err)
	}

	// 获取题目数据
	err := query.Offset((page - 1) * size).Limit(size).Find(&problems).Error
	if err != nil {
		return nil, fmt.Errorf("获取题目列表失败: %w", err)
	}

	// 返回分页结果
	return &model.Page[model.Problem]{
		Total: total,
		Data:  problems,
	}, nil
}

// GetProblemByID 查询指定 ID 的题目详情
func (ps *ProblemService) GetProblemByID(id int) (*model.Problem, error) {
	var problem model.Problem
	if err := db.DB.First(&problem, id).Error; err != nil {
		return nil, fmt.Errorf("获取题目详情失败: %w", err)
	}

	// 返回题目详情
	return &problem, nil
}

// GetProblemTestCase 查询指定题目的测试用例
func (ps *ProblemService) GetProblemTestCase(problemID uint64) ([]model.Case, error) {
	var record model.TestCase

	if err := db.DB.Where("problem_id = ?", problemID).First(&record).Error; err != nil {
		return nil, fmt.Errorf("查询测试用例失败: %v", err)
	}

	return record.Cases, nil
}
