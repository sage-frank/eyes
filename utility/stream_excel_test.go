package utility

import (
	"database/sql"
	"fmt"
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

	query := "SELECT id, title, content, is_published, created_at, updated_at FROM `articles` "
	options := ExportOptions{
		DB:        db,
		Query:     query,
		OutPath:   "output.xlsx",
		BatchSize: 1000,
		SheetName: "Sheet1",
	}

	if err := ExportDataToExcel(options); err != nil {
		fmt.Printf("导出数据失败: %v\n", err)
	}
}
