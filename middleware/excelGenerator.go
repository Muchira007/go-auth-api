package middleware

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muchira007/jambo-green-go/initializers"
	"github.com/xuri/excelize/v2"
)

func DownloadExcelFile(c *gin.Context) {
	// Get all sales data
	var salesData []struct {
		Gender   string `json:"gender"`
		Count    int64  `json:"count"`
		Product  string `json:"product"`
		Customer string `json:"customer"`
	}

	// Retrieve sales data from the database
	if err := initializers.DB.Table("sales").
		Select("customers.gender, COUNT(*) as count").
		Joins("JOIN customers ON sales.customer_id = customers.id").
		Joins("JOIN products ON sales.product_id = products.id").
		Group("customers.gender, products.name").
		Scan(&salesData).Error; err != nil {
		fmt.Println("Error retrieving sales data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sales data"})
		return
	}

	// Create a new Excel file
	file := excelize.NewFile()

	// Define headers for sales data
	headers := []string{"Gender", "Count", "Product", "Customer"}
	for i, header := range headers {
		file.SetCellValue("Sheet1", fmt.Sprintf("%s%d", string(rune(65+i)), 1), header)
	}

	// Populate the sheet with sales data
	for i, data := range salesData {
		row := i + 2 // Start from the second row
		file.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), data.Gender)
		file.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), data.Count)
		file.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), data.Product)
		file.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), data.Customer)
	}

	// Dynamic ranges for the chart
	lastRow := len(salesData) + 1 // +1 for the header row

	// Add the chart to the file
	if err := file.AddChart("Sheet1", "E1", &excelize.Chart{
		Type: excelize.Col3DClustered,
		Series: []excelize.ChartSeries{
			{
				Name:       "Sheet1!$A$1",
				Categories: "Sheet1!$A$2:$A$" + fmt.Sprint(lastRow),
				Values:     "Sheet1!$B$2:$B$" + fmt.Sprint(lastRow),
			},
		},
		Title: []excelize.RichTextRun{
			{
				Text: "Gender Distribution",
			},
		},
	}); err != nil {
		fmt.Println("Error adding chart:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add chart to Excel file"})
		return
	}

	// Define the path to the uploads directory
	uploadsDir := "./uploads"
	// Ensure the uploads directory exists
	if err := os.MkdirAll(uploadsDir, os.ModePerm); err != nil {
		fmt.Println("Error creating uploads directory:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create uploads directory"})
		return
	}

	// Define the path for the file
	fileName := fmt.Sprintf("sales-%d.xlsx", time.Now().Unix())
	filePath := filepath.Join(uploadsDir, fileName)

	// Save the file
	if err := file.SaveAs(filePath); err != nil {
		fmt.Println("Error saving Excel file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save Excel file"})
		return
	}

	// Create a URL for the file
	fileURL := fmt.Sprintf("http://localhost:3000/uploads/%s", fileName)

	// Return the file URL in the response
	c.JSON(http.StatusOK, gin.H{"excelUrl": fileURL})
}
