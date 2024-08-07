package middleware

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	// "time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gin-gonic/gin"
	"github.com/muchira007/jambo-green-go/initializers"
	"github.com/muchira007/jambo-green-go/models"
)

// TemplateData holds the data for the PDF template
type TemplateData struct {
	Content  string
	Products []ProductData
}

type ProductData struct {
	Name      string
	Price     float64
	ImageData string
}

func DownloadPDF(c *gin.Context) {
	// Get products
	var products []models.Product
	if err := initializers.DB.Find(&products).Error; err != nil {
		fmt.Println("Error retrieving products:", err)
		c.String(http.StatusInternalServerError, "Failed to retrieve products")
		return
	}

	// Encode image data to Base64 for each product
	var productData []ProductData
	for _, product := range products {
		pd := ProductData{
			Name:  product.Name,
			Price: product.Price,
		}
		if product.ImageData != nil {
			pd.ImageData = base64.StdEncoding.EncodeToString(product.ImageData)
		}
		productData = append(productData, pd)
	}

	// Load and parse the HTML template
	tmpl, err := template.ParseFiles("templates/pdf.html")
	if err != nil {
		fmt.Println("Error loading template:", err)
		c.String(http.StatusInternalServerError, "Unable to load template file")
		return
	}

	// Prepare data for the template
	data := TemplateData{
		Content:  "This is a list of all products.",
		Products: productData,
	}

	// Execute the template with data
	var htmlBuffer bytes.Buffer
	if err := tmpl.Execute(&htmlBuffer, data); err != nil {
		fmt.Println("Error executing template:", err)
		c.String(http.StatusInternalServerError, "Unable to execute template")
		return
	}

	// Create a new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		fmt.Println("Error creating PDF generator:", err)
		c.String(http.StatusInternalServerError, "Unable to create PDF generator")
		return
	}

	// Set PDF options
	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)

	// Add a page to the PDF
	page := wkhtmltopdf.NewPageReader(bytes.NewReader(htmlBuffer.Bytes()))
	page.EnableLocalFileAccess.Set(true)

	pdfg.AddPage(page)

	// Generate the PDF
	err = pdfg.Create()
	if err != nil {
		fmt.Println("Error generating PDF:", err)
		c.String(http.StatusInternalServerError, "Unable to generate PDF file")
		return
	}

	// Define the path to the uploads directory
	uploadsDir := "./uploads"
	// Ensure the uploads directory exists
	if err := os.MkdirAll(uploadsDir, os.ModePerm); err != nil {
		fmt.Println("Error creating uploads directory:", err)
		c.String(http.StatusInternalServerError, "Unable to create uploads directory")
		return
	}

	// Save PDF to a temporary file
	// Save PDF to the uploads folder
	fileName := fmt.Sprintf("report-%d.pdf", time.Now().Unix())
	filePath := filepath.Join(uploadsDir, fileName)
	tmpFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file in uploads directory:", err)
		c.String(http.StatusInternalServerError, "Unable to create file in uploads directory")
		return
	}
	defer os.Remove(filePath) // Clean up the file after serving

	if _, err := tmpFile.Write(pdfg.Bytes()); err != nil {
		fmt.Println("Error writing to file:", err)
		c.String(http.StatusInternalServerError, "Unable to write to file")
		return
	}

	// Create a URL for the temporary file
	// fileURL := fmt.Sprintf("http://localhost:3000/uploads/%s", filepath.Base(tmpFile.Name()))
	// fmt.Print(tmpFile.Name())
	fileURL := fmt.Sprintf("http://localhost:3000/uploads/%s", fileName)
	fmt.Print(fileURL)

	c.JSON(http.StatusOK, gin.H{"pdfUrl": fileURL})
}
