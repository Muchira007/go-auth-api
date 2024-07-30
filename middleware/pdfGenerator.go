package middleware

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func DownloadPDF(c *gin.Context) {
	// Create a new PDF
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Add a page
	pdf.AddPage()

	// Set font
	pdf.SetFont("Arial", "B", 16)

	// Add a cell
	pdf.Cell(40, 10, "Hello, this is a PDF report!")

	// Buffer to store PDF
	var buffer bytes.Buffer
	if err := pdf.Output(&buffer); err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Unable to generate PDF file")
		return
	}

	// Set the buffer as the response data
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=report.pdf")
	c.Data(http.StatusOK, "application/pdf", buffer.Bytes())
}
