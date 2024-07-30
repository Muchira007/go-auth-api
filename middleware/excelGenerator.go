package middleware

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func DownloadExcelFile(c *gin.Context) {
	// Create a new Excel file
	f := excelize.NewFile()

	// Create a new sheet
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Unable to create new sheet")
		return
	}

	// Set value of a cell
	f.SetCellValue("Sheet1", "A1", "Name")
	f.SetCellValue("Sheet1", "B1", "Age")
	f.SetCellValue("Sheet1", "A2", "John")
	f.SetCellValue("Sheet1", "B2", 30)
	f.SetCellValue("Sheet1", "A3", "Doe")
	f.SetCellValue("Sheet1", "B3", 25)

	// Set active sheet
	f.SetActiveSheet(index)

	var buffer bytes.Buffer
	if err := f.Write(&buffer); err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Unable to generate Excel file")
		return
	}

	// Set the buffer as the response data
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=report.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer.Bytes())
}
