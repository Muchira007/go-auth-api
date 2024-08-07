package middleware

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muchira007/jambo-green-go/initializers"
	"github.com/muchira007/jambo-green-go/models"
	"github.com/xuri/excelize/v2"
)

func DownloadExcelFile(c *gin.Context) {
	// Get products from the database
	var products []models.Product
	if err := initializers.DB.Find(&products).Error; err != nil {
		fmt.Println("Error retrieving products:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}

	// Create a new Excel file
	file := excelize.NewFile()

	// Define headers
	headers := []string{"ID", "Name", "Price", "Image"}
	for i, header := range headers {
		file.SetCellValue("Sheet1", fmt.Sprintf("%s%d", string(rune(65+i)), 1), header)
	}

	// Populate the sheet with product data
	for i, product := range products {
		row := i + 2 // Start from the second row

		// Encode image data to Base64
		var imageData string
		if product.ImageData != nil {
			imageData = base64.StdEncoding.EncodeToString(product.ImageData)
		}

		file.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), product.ID)
		file.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), product.Name)
		file.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), product.Price)
		file.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), imageData)
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
	fileName := fmt.Sprintf("products-%d.xlsx", time.Now().Unix())
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
