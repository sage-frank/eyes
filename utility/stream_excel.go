package utility

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // 导入 MySQL 驱动
	"github.com/xuri/excelize/v2"
)

// ExportOptions 定义了导出数据到 Excel 的选项
type ExportOptions struct {
	DB        *sql.DB // 已经建立的数据库连接
	Query     string  // 查询语句
	OutPath   string  // 输出文件路径
	BatchSize int     // 批量处理大小
	SheetName string  // sheetName
}

// ExportDataToExcel 根据提供的选项导出数据到 Excel 文件
func ExportDataToExcel(options ExportOptions) error {
	start := time.Now() // 记录程序开始时间
	f := excelize.NewFile()
	sheet := options.SheetName
	if sheet == "" {
		sheet = "Sheet1"
	}

	_, err := f.NewSheet(sheet)
	if err != nil {
		return fmt.Errorf("新建sheet失败: %w", err)
	}

	// 获取流式写入器
	streamWriter, err := f.NewStreamWriter(sheet)
	if err != nil {
		return fmt.Errorf("无法创建流式写入器: %v", err)
	}

	row := 1

	// 执行查询
	rows, err := options.DB.Query(options.Query)
	if err != nil {
		return fmt.Errorf("查询失败: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("获取列名失败: %v", err)
	}

	// 写入表头
	header := make([]interface{}, len(columns))
	for i, col := range columns {
		header[i] = col
	}
	if err = streamWriter.SetRow("A1", header); err != nil {
		return fmt.Errorf("写入表头失败: %v", err)
	}

	// 处理查询结果
	for rows.Next() {
		columnPointers := make([]interface{}, len(columns))
		columnValues := make([]interface{}, len(columns))
		for i := range columnPointers {
			columnPointers[i] = &columnValues[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return fmt.Errorf("读取行数据失败: %v", err)
		}

		rowValues := make([]interface{}, len(columns))
		for i, col := range columnValues {
			rowValues[i] = col
		}

		if err := streamWriter.SetRow(fmt.Sprintf("A%d", row+1), rowValues); err != nil {
			return fmt.Errorf("写入行数据失败: %v", err)
		}

		row++
	}

	// 结束流式写入
	if err = streamWriter.Flush(); err != nil {
		return fmt.Errorf("流式写入失败: %v", err)
	}

	// 保存 Excel 文件
	if err := f.SaveAs(options.OutPath); err != nil {
		return fmt.Errorf("保存 Excel 文件失败: %v", err)
	}

	elapsed := time.Since(start) // 计算程序执行时长
	fmt.Printf("数据成功导出到 %s，程序执行时长：%s\n", options.OutPath, elapsed)

	return nil
}
