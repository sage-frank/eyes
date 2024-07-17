package utility

import (
	"database/sql"
	"log"
	"testing"
)

func TestExportDataToExcel(t *testing.T) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:33061)/bigdata?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		// 在这里记录日志或者处理错误
		log.Printf("无法连接数据库: %v", err)
		return
	}

	defer db.Close()

	options := ExportOptions{
		DB:             db,
		TableName:      "article",
		Title:          []any{"ID", "标题", "内容", "是否发布", "创建时间", "更新时间"},
		Columns:        []string{"id", "title", "content", "is_published", "created_at", "updated_at"},
		QueryCondition: "id > 0",
		OutPath:        "output.xlsx",
		BatchSize:      10000,
	}

	n, err := ExportDataToExcel(options)
	if err != nil {
		// 在这里记录日志或者处理错误
		log.Printf("导出数据时发生错误: %v", err)
	}
	log.Println(n)
}
