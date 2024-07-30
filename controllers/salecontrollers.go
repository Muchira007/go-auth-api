package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muchira007/jambo-green-go/initializers"
	"github.com/muchira007/jambo-green-go/models"
)

func RecordSale(c *gin.Context) {
	var body struct {
		Name            string `json:"name" binding:"required"`
		DateOfSale      string `json:"date_of_sale" binding:"required"`
		Gender          string `json:"gender" binding:"required"`
		PhoneNumber     string `json:"phone_number" binding:"required"`
		NationalID      string `json:"national_id" binding:"required"`
		Geolocation     string `json:"geolocation" binding:"required"`
		Country         string `json:"country" binding:"required"`
		County          string `json:"county" binding:"required"`
		Subcounty       string `json:"subcounty" binding:"required"`
		Village         string `json:"village" binding:"required"`
		ProductID       uint   `json:"product_id" binding:"required"`
		SerialNumber    string `json:"serial_number" binding:"required"`
		PaymentOption   string `json:"payment_option" binding:"required"`
		StatusOfAccount string `json:"status_of_account" binding:"required"`
		Quantity        int    `json:"quantity" binding:"required"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var customer models.Customer
	if err := initializers.DB.Where("national_id = ?", body.NationalID).FirstOrCreate(&customer, models.Customer{
		Name:        body.Name,
		Gender:      body.Gender,
		PhoneNumber: body.PhoneNumber,
		NationalID:  body.NationalID,
		Geolocation: body.Geolocation,
		Country:     body.Country,
		County:      body.County,
		Subcounty:   body.Subcounty,
		Village:     body.Village,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find or create customer"})
		return
	}

	var product models.Product
	if err := initializers.DB.First(&product, body.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	total := product.Price * float64(body.Quantity)

	sale := models.Sale{
		DateOfSale:      time.Now(),
		ProductID:       body.ProductID,
		Product:         product,
		CustomerID:      customer.ID,
		Customer:        customer,
		SerialNumber:    body.SerialNumber,
		PaymentOption:   body.PaymentOption,
		StatusOfAccount: body.StatusOfAccount,
		Quantity:        body.Quantity,
		Total:           total,
	}

	if err := initializers.DB.Create(&sale).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record sale"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sale": sale})
}
