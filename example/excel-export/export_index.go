package main

import (
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

func main() {
	t1 := time.Now()
	f := excelize.NewFile()
	sheet := "Sheet1"

	for i := 1; i <= 1000000; i++ {
		cell := "A" + strconv.Itoa(i)
		f.SetCellValue(sheet, cell, i)
	}

	if err := f.SaveAs("test_golang.xlsx"); err != nil {
		println(err.Error())
	}
	println(time.Since(t1).Milliseconds())
}
