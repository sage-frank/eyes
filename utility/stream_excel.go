package utility

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // 导入 MySQL 驱动
	"github.com/xuri/excelize/v2"
)

// ExportOptions 定义了导出数据到 Excel 的选项

type ExportOptions struct {
	DB             *sql.DB  // 已经建立的数据库连接
	TableName      string   // 表名
	Columns        []string // 需要导出的列名
	Title          []any    // 需要导出的列名
	QueryCondition string   // 查询条件
	OutPath        string   // 输出文件路径
	BatchSize      int      // 批量处理大小
	SheetName      string   // sheetName
}

// ExportDataToExcel 根据提供的选项导出数据到 Excel 文件
func ExportDataToExcel(options ExportOptions) (int64, error) {
	var (
		id          int64
		title       string
		content     string
		isPublished int
		createdAt   time.Time
		updatedAt   time.Time

		maxID int64 = 0
		row         = 2
	)

	start := time.Now() // 记录程序开始时间
	f := excelize.NewFile()
	sheet := options.SheetName
	if sheet == "" {
		sheet = "Sheet1"
	}

	_, err := f.NewSheet(sheet)
	if err != nil {
		return 0, fmt.Errorf("新建sheet失败: %w", err)
	}

	// 获取流式写入器
	streamWriter, err := f.NewStreamWriter(sheet)
	if err != nil {
		return 0, fmt.Errorf("无法创建流式写入器: %v", err)
	}

	// 写入表头
	if err = streamWriter.SetRow("A1", options.Title); err != nil {
		return 0, fmt.Errorf("写入表头失败: %v", err)
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE id > ? ORDER BY id LIMIT ?", strings.Join(options.Columns, ", "), options.TableName)

	for {
		// 执行查询
		rows, err := options.DB.Query(query, maxID, options.BatchSize)
		if err != nil {
			return 0, fmt.Errorf("查询失败: %v", err)
		}
		count := 0
		for rows.Next() {
			// 这里需要根据实际的列名和数据类型来扫描数据
			// 假设我们有以下结构体

			err := rows.Scan(&id, &title, &content, &isPublished, &createdAt, &updatedAt)
			if err != nil {
				return 0, fmt.Errorf("读取行数据失败: %v", err)
			}

			if err := streamWriter.SetRow(fmt.Sprintf("A%d", row), []any{
				fmt.Sprintf("%d", id), title, content, isPublished, createdAt, updatedAt,
			}); err != nil {
				return 0, fmt.Errorf("写入行数据失败: %v", err)
			}
			count++
			row++
			maxID = id
		}
		if err = rows.Close(); err != nil {
			return 0, fmt.Errorf("关闭行数据失败: %v", err)
		}
		if count == 0 {
			break // 如果没有更多数据，则退出循环
		}
	}
	// 结束流式写入
	if err = streamWriter.Flush(); err != nil {
		return 0, fmt.Errorf("流式写入失败: %v", err)
	}

	// 保存 Excel 文件
	if err := f.SaveAs(options.OutPath); err != nil {
		return 0, fmt.Errorf("保存 Excel 文件失败: %v", err)
	}

	elapsed := time.Since(start) // 计算程序执行时长
	fmt.Printf("数据成功导出到 %s，程序执行时长：%s\n", options.OutPath, elapsed)

	return maxID, nil
}
