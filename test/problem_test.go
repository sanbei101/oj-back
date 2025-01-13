package test

import (
	"fmt"
	"log"
	"math/rand"
	"oj-back/app/db"
	"oj-back/app/model"
	"oj-back/app/service"
	"os"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	TestDB            *gorm.DB
	ProblemServiceApp = service.ProblemServiceApp
)

func InitTestDB() error {
	var err error
	const dsn string = "host=localhost user=testuser password=justfortest dbname=testdatabase sslmode=disable "
	TestDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel:      logger.Warn,     // 设置日志级别
				SlowThreshold: time.Second * 2, // 慢查询的时间阈值
				Colorful:      false,           // 是否开启日志输出的彩色
			},
		),
	})
	db.DB = TestDB
	if err != nil {
		return err
	}
	migrator := db.DB.Migrator()
	if err := migrator.DropTable(&model.Problem{}, &model.TestCase{}); err != nil {
		return err
	}
	if err := migrator.CreateTable(&model.Problem{}, &model.TestCase{}); err != nil {
		return err
	}
	if err = db.DB.Exec(`CREATE INDEX IF NOT EXISTS idx_problems_tags_gin ON problems USING GIN (tags);`).Error; err != nil {
		return err
	}
	return nil
}

func InsertTestData() {
	var problems []model.Problem
	for i := 0; i < 100000; i++ {
		problems = append(problems, model.Problem{
			Name:        fmt.Sprintf("test%d", i+1),
			Description: fmt.Sprintf("description%d", i+1),
			Tags:        []string{fmt.Sprintf("tag%d", i%10+1), fmt.Sprintf("tag%d", (i+1)%10+1)},
		})
	}

	var testCases []model.TestCase
	for i := 0; i < 100000; i++ {
		testCases = append(testCases, model.TestCase{
			ProblemID: uint64(i + 1),
			Cases: []model.Case{
				{Input: "1 2", ExpectedOutput: "3"},
				{Input: "3 5", ExpectedOutput: "8"},
				{Input: "10 15", ExpectedOutput: "25"},
			},
		})
	}

	batchSize := 5000 // 每批次插入 5000 条
	if err := db.DB.CreateInBatches(&problems, batchSize).Error; err != nil {
		log.Fatalf("批量插入测试数据失败: %v", err)
	}
	if err := db.DB.CreateInBatches(&testCases, batchSize).Error; err != nil {
		log.Fatalf("批量插入测试数据失败: %v", err)
	}
}
func RemoveTestData() {
	migrator := db.DB.Migrator()
	migrator.DropTable(&model.Problem{}, &model.TestCase{})
}
func TestGetAllProblems(t *testing.T) {
	if err := InitTestDB(); err != nil {
		t.Fatalf("初始化测试数据库失败: %v", err)
	}
	problems := []model.Problem{
		{
			Name:        "test1",
			Description: "test1",
			Tags:        []string{"tag1", "tag2"},
		},
		{
			Name:        "test2",
			Description: "test2",
			Tags:        []string{"tag2", "tag3"},
		},
	}
	db.DB.CreateInBatches(&problems, len(problems))
	// 测试查询所有题目
	page, size := 1, 10
	keyword := ""
	result, err := ProblemServiceApp.GetAllProblems(page, size, keyword)
	if err != nil {
		t.Fatalf("查询题目失败: %v", err)
	}
	if result.Total != 2 {
		t.Fatalf("查询题目总数错误: %d", result.Total)
	}
	if len(result.Data) != 2 {
		t.Fatalf("查询题目数量错误: %d", len(result.Data))
	}
	// 清理测试数据
	if err := db.DB.Delete(&problems).Error; err != nil {
		t.Fatalf("删除测试数据失败: %v", err)
	}
}

func BenchmarkGetAllProblems(b *testing.B) {
	if err := InitTestDB(); err != nil {
		b.Fatalf("初始化测试数据库失败: %v", err)
	}
	InsertTestData()
	b.ResetTimer()

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		page, size := rand.Intn(1000), 10
		keyword := ""
		_, err := ProblemServiceApp.GetAllProblems(page, size, keyword)
		if err != nil {
			b.Fatalf("查询题目失败: %v", err)
		}
	}
	RemoveTestData()
}

func TestGetProblemByID(t *testing.T) {
	if err := InitTestDB(); err != nil {
		t.Fatalf("初始化测试数据库失败: %v", err)
	}
	problems := []model.Problem{
		{
			Name:        "test1",
			Description: "test1",
			Tags:        []string{"tag1", "tag2"},
		},
		{
			Name:        "test2",
			Description: "test2",
			Tags:        []string{"tag2", "tag3"},
		},
	}
	db.DB.CreateInBatches(&problems, len(problems))
	// 测试查询指定 ID 的题目
	problem, err := ProblemServiceApp.GetProblemByID(1)
	if err != nil {
		t.Fatalf("查询题目详情失败: %v", err)
	}
	if problem.ID != 1 {
		t.Fatalf("查询题目 ID 错误: %d", problem.ID)
	}
	// 清理测试数据
	if err := db.DB.Delete(&problems).Error; err != nil {
		t.Fatalf("删除测试数据失败: %v", err)
	}
}

func BenchmarkGetProblemByID(b *testing.B) {
	if err := InitTestDB(); err != nil {
		b.Fatalf("初始化测试数据库失败: %v", err)
	}

	InsertTestData()
	b.ResetTimer()

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		_, err := ProblemServiceApp.GetProblemByID(rand.Intn(100000) + 1)
		if err != nil {
			b.Fatalf("查询题目详情失败: %v", err)
		}
	}

	RemoveTestData()
}

func TestGetProblemTestCase(t *testing.T) {
	if err := InitTestDB(); err != nil {
		t.Fatalf("初始化测试数据库失败: %v", err)
	}
	problems := []model.Problem{
		{
			Name:        "加法运算",
			Description: "给定两个数字，输出它们的和。",
			Tags:        []string{"数学", "简单"},
		},
	}

	testCases := []model.TestCase{
		{
			ProblemID: 1,
			Cases: []model.Case{
				{Input: "1 2", ExpectedOutput: "3"},
				{Input: "3 5", ExpectedOutput: "8"},
				{Input: "10 15", ExpectedOutput: "25"},
			},
		},
	}

	db.DB.Create(&problems)
	db.DB.Create(&testCases)

	// 测试查询指定题目的测试用例
	testCase, err := ProblemServiceApp.GetProblemTestCase(1)
	if err != nil {
		t.Fatalf("查询测试用例失败: %v", err)
	}

	if testCase[0].Input != "1 2" || testCase[0].ExpectedOutput != "3" {
		t.Fatalf("测试用例错误: %v", testCase)
	}
}

func BenchmarkGetProblemTestCase(b *testing.B) {
	if err := InitTestDB(); err != nil {
		b.Fatalf("初始化测试数据库失败: %v", err)
	}

	InsertTestData()

	b.ResetTimer()

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		_, err := ProblemServiceApp.GetProblemTestCase(uint64(rand.Intn(100000) + 1))
		if err != nil {
			b.Fatalf("查询测试用例失败: %v", err)
		}
	}

	RemoveTestData()
}
