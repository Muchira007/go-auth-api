package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/muchira007/jambo-green-go/initializers"
    "github.com/muchira007/jambo-green-go/models"
    "strconv"
)

func CheckProductQuantity(c *gin.Context) {
	// Check if the product quantity is greater than 0
		productID := c.Param("product_id")
		requestQuantityStr := c.Query("quantity")

		//convert the request quantity to an integer
		requestQuantity, err := strconv.Atoi(requestQuantityStr)
		if err != nil || requestQuantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity"})
			c.Abort()
			return
		}


		// Get the product from the database
		var product models.Product
		if err := initializers.DB.First(&product, productID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			c.Abort()
			return
		}

		// Check if the product quantity is greater than the requested quantity
		if product.Quantity < requestQuantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient product quantity"})
			c.Abort()
			return
		}

		c.Next()


}