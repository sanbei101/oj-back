package service

import (
	"fmt"
	"log"
	"math/rand"
	"oj-back/app/db"
	"oj-back/app/model"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	TestDB         *gorm.DB
	InsertDataSize = 100000
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
	for i := 0; i < InsertDataSize; i++ {
		problems = append(problems, model.Problem{
			Name:        fmt.Sprintf("test%d", i+1),
			Description: fmt.Sprintf("description%d", i+1),
			Tags:        []string{fmt.Sprintf("tag%d", i%10+1), fmt.Sprintf("tag%d", (i+1)%10+1)},
		})
	}

	var testCases []model.TestCase
	for i := 0; i < InsertDataSize; i++ {
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
	err := InitTestDB()
	assert.NoError(t, err, "初始化测试数据库失败: %v", err)
	InsertTestData()

	page, size := rand.Intn(10)+10, rand.Intn(100)+10
	keyword := "test"
	result, err := ProblemServiceApp.GetAllProblems(page, size, keyword)
	assert.NoError(t, err, "查询题目失败: %v", err)
	assert.Equal(t, int64(InsertDataSize), result.Total, "查询题目总数错误")
	assert.Equal(t, size, len(result.Data), "查询题目数量错误")

	RemoveTestData()
}

func BenchmarkGetAllProblems(b *testing.B) {
	err := InitTestDB()
	assert.NoError(b, err, "初始化测试数据库失败: %v", err)
	InsertTestData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		page, size := rand.Intn(1000), rand.Intn(90)+10
		keyword := ""
		_, err := ProblemServiceApp.GetAllProblems(page, size, keyword)
		assert.NoError(b, err, "查询题目失败: %v", err)
	}

	RemoveTestData()
}

func TestGetProblemByID(t *testing.T) {
	err := InitTestDB()
	assert.NoError(t, err, "初始化测试数据库失败: %v", err)
	InsertTestData()

	var id = rand.Intn(InsertDataSize) + 1
	problem, err := ProblemServiceApp.GetProblemByID(id)
	assert.NoError(t, err, "查询题目详情失败: %v", err)
	assert.Equal(t, uint64(id), problem.ID, "查询题目 ID 错误")

	RemoveTestData()
}

func BenchmarkGetProblemByID(b *testing.B) {
	err := InitTestDB()
	assert.NoError(b, err, "初始化测试数据库失败: %v", err)
	InsertTestData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ProblemServiceApp.GetProblemByID(rand.Intn(InsertDataSize) + 1)
		assert.NoError(b, err, "查询题目详情失败: %v", err)
	}

	RemoveTestData()
}

func TestGetProblemTestCase(t *testing.T) {
	err := InitTestDB()
	assert.NoError(t, err, "初始化测试数据库失败: %v", err)
	InsertTestData()

	testCase, err := ProblemServiceApp.GetProblemTestCase(uint64(rand.Intn(InsertDataSize) + 1))

	assert.NoError(t, err, "查询测试用例失败: %v", err)
	assert.Equal(t, "1 2", testCase[0].Input, "测试用例输入错误")
	assert.Equal(t, "3", testCase[0].ExpectedOutput, "测试用例输出错误")

	RemoveTestData()
}

func BenchmarkGetProblemTestCase(b *testing.B) {
	err := InitTestDB()
	assert.NoError(b, err, "初始化测试数据库失败: %v", err)
	InsertTestData()
	b.ResetTimer()

	// 运行基准测试
	for i := 0; i < b.N; i++ {
		_, err := ProblemServiceApp.GetProblemTestCase(uint64(rand.Intn(InsertDataSize) + 1))
		assert.NoError(b, err, "查询测试用例失败: %v", err)
	}

	RemoveTestData()
}
