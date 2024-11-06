package db

import (
	"fmt"
	"log"
	"oj-back/app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	var err error
	const dsn string = "host=125.91.126.101 user=ghr password=GZH031ghr dbname=mydatabase port=5432 sslmode=disable "
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	fmt.Println("数据库连接成功")
	migrator := DB.Migrator()
	// 清空所有表
	if err := migrator.DropTable(&models.Problem{}, &models.TestCase{}); err != nil {
		log.Fatalf("删除表失败: %v", err)
	}
	// 重新创建所有表
	if err := migrator.CreateTable(&models.Problem{}, &models.TestCase{}); err != nil {
		log.Fatalf("创建表失败: %v", err)
	}
	models.InsertData(DB)
	return DB
}
