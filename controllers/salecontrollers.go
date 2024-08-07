package controllers

//national_id is for agent id initiating sale

import (
	// "fmt"
	"net/http"
	"strconv" // Import strconv for type conversion
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muchira007/jambo-green-go/initializers"
	"github.com/muchira007/jambo-green-go/models"
)

// func RecordSale(c *gin.Context) {
// 	var body struct {
// 		Name            string  `json:"name" binding:"required"`
// 		DateOfSale      string  `json:"date_of_sale" binding:"required"`
// 		Gender          string  `json:"gender" binding:"required"`
// 		PhoneNumber     string  `json:"phone_number" binding:"required"`
// 		CustomerID      int     `json:"customer_id" binding:"required"`
// 		Country         string  `json:"country" binding:"required"`
// 		County          string  `json:"county" binding:"required"`
// 		Subcounty       string  `json:"subcounty" binding:"required"`
// 		Village         string  `json:"village" binding:"required"`
// 		Ward            string  `json:"ward" binding:"required"`
// 		Latitude        float64 `json:"latitude" binding:"required"`
// 		Longitude       float64 `json:"longitude" binding:"required"`
// 		ProductID       uint    `json:"product_id" binding:"required"`
// 		SerialNumber    string  `json:"serial_number"`
// 		PaymentOption   string  `json:"payment_option" binding:"required"`
// 		StatusOfAccount string  `json:"status_of_account" binding:"required"`
// 		Quantity        int     `json:"quantity" binding:"required"`
// 		NationalID      uint    `json:"national_id" binding:"required"`
// 	}

// 	fmt.Println(body)

// 	// Rest of the function remains the same
// 	if err := c.Bind(&body); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
// 		return
// 	}

// 	// Check if the agent exists
// 	var agent models.User
// 	if err := initializers.DB.Where("national_id = ?", body.NationalID).First(&agent).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
// 		return
// 	}

// 	var customer models.Customer
// 	// Convert CustomerID to string if necessary
// 	customerNationalID := strconv.Itoa(body.CustomerID) // Convert to string if national_id is text
// 	if err := initializers.DB.Where("national_id = ?", customerNationalID).FirstOrCreate(&customer, models.Customer{
// 		Name:        body.Name,
// 		Gender:      body.Gender,
// 		PhoneNumber: body.PhoneNumber,
// 		CustomerID:  uint(body.CustomerID),
// 		Country:     body.Country,
// 		County:      body.County,
// 		Subcounty:   body.Subcounty,
// 		Village:     body.Village,
// 		Ward:        body.Ward,
// 	}).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find or create customer"})
// 		return
// 	}

// 	var product models.Product
// 	if err := initializers.DB.First(&product, body.ProductID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
// 		return
// 	}

// 	total := product.Price * float64(body.Quantity)

// 	sale := models.Sale{
// 		DateOfSale:      time.Now(),
// 		ProductID:       body.ProductID,
// 		Product:         product,
// 		CustomerID:      customer.ID,
// 		Customer:        customer,
// 		SerialNumber:    body.SerialNumber,
// 		PaymentOption:   body.PaymentOption,
// 		StatusOfAccount: body.StatusOfAccount,
// 		Quantity:        body.Quantity,
// 		Total:           total,
// 		NationalID:      body.NationalID,
// 		Latitude:        body.Latitude,
// 		Longitude:       body.Longitude,
// 	}

// 	if err := initializers.DB.Create(&sale).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record sale"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"sale": sale})
// }

func RecordSale(c *gin.Context) {
	var body struct {
		Name            string  `json:"name" binding:"required"`
		DateOfSale      string  `json:"date_of_sale" binding:"required"`
		Gender          string  `json:"gender" binding:"required"`
		PhoneNumber     string  `json:"phone_number" binding:"required"`
		CustomerID      int     `json:"customer_id" binding:"required"`
		Country         string  `json:"country" binding:"required"`
		County          string  `json:"county" binding:"required"`
		Subcounty       string  `json:"subcounty" binding:"required"`
		Village         string  `json:"village" binding:"required"`
		// Ward            string  `json:"ward" binding:"required"`
		Latitude        float64 `json:"latitude" binding:"required"`
		Longitude       float64 `json:"longitude" binding:"required"`
		ProductName     string  `json:"product_name" binding:"required"`
		SerialNumber    string  `json:"serial_number"`
		PaymentOption   string  `json:"payment_option" binding:"required"`
		StatusOfAccount string  `json:"status_of_account" binding:"required"`
		Quantity        int     `json:"quantity" binding:"required"`
		NationalID      uint    `json:"national_id" binding:"required"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Check if the agent exists
	var agent models.User
	if err := initializers.DB.Where("national_id = ?", body.NationalID).First(&agent).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	var customer models.Customer
	customerNationalID := strconv.Itoa(body.CustomerID)
	if err := initializers.DB.Where("national_id = ?", customerNationalID).FirstOrCreate(&customer, models.Customer{
		Name:        body.Name,
		Gender:      body.Gender,
		PhoneNumber: body.PhoneNumber,
		CustomerID:  uint(body.CustomerID),
		Country:     body.Country,
		County:      body.County,
		Subcounty:   body.Subcounty,
		Village:     body.Village,
		// Ward:        body.Ward,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find or create customer"})
		return
	}

	var product models.Product
	if err := initializers.DB.Where("name = ?", body.ProductName).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	total := product.Price * float64(body.Quantity)

	sale := models.Sale{
		DateOfSale:      time.Now(),
		ProductID:       product.ID,
		Product:         product,
		CustomerID:      customer.ID,
		Customer:        customer,
		SerialNumber:    body.SerialNumber,
		PaymentOption:   body.PaymentOption,
		StatusOfAccount: body.StatusOfAccount,
		Quantity:        body.Quantity,
		Total:           total,
		NationalID:      body.NationalID,
		Latitude:        body.Latitude,
		Longitude:       body.Longitude,
	}

	if err := initializers.DB.Create(&sale).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record sale"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sale": sale})
}

// get gender distribution

func GetSalesByGender(c *gin.Context) {
	var results []struct {
		Gender string        `json:"gender"`
		Count  int64         `json:"count"`
		Sales  []models.Sale `json:"sales"`
	}

	var salesByGender []struct {
		Gender string `json:"gender"`
		Count  int64  `json:"count"`
	}

	// Get count of sales by customer gender
	if err := initializers.DB.Table("sales").
		Select("customers.gender, COUNT(*) as count").
		Joins("JOIN customers ON sales.customer_id = customers.id").
		Group("customers.gender").
		Scan(&salesByGender).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sales count by gender", "details": err.Error()})
		return
	}

	// Get detailed sales data by customer gender
	for _, genderCount := range salesByGender {
		var sales []models.Sale
		if err := initializers.DB.Joins("JOIN customers ON sales.customer_id = customers.id").
			Where("customers.gender = ?", genderCount.Gender).
			Find(&sales).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sales data for gender " + genderCount.Gender, "details": err.Error()})
			return
		}

		results = append(results, struct {
			Gender string        `json:"gender"`
			Count  int64         `json:"count"`
			Sales  []models.Sale `json:"sales"`
		}{
			Gender: genderCount.Gender,
			Count:  genderCount.Count,
			Sales:  sales,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}

// get sales by national id
func GetSalesByNationalID(c *gin.Context) {
	nationalIDStr := c.Param("national_id")
	nationalID, err := strconv.Atoi(nationalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid national_id format"})
		return
	}

	var count int64
	if err := initializers.DB.Model(&models.Sale{}).
		Where("national_id = ?", nationalID).
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sales count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"national_id": nationalID, "sales_count": count})
}

func GetAgentSales(c *gin.Context) {
	var results []struct {
		NationalID uint  `json:"national_id"`
		Count      int64 `json:"sales_count"`
	}

	if err := initializers.DB.Model(&models.Sale{}).
		Select("national_id, COUNT(*) as sales_count").
		Group("national_id").
		Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sales count by national_id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}
