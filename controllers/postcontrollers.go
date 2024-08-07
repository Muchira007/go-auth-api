package controllers

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/muchira007/jambo-green-go/initializers"
	"github.com/muchira007/jambo-green-go/models"
	"io"
	"net/http"
	"strconv" // Add the import statement for strconv
)

// PostsCreate handles the creation of a new product, including image upload.
func PostsCreate(c *gin.Context) {
	// Initialize imageData as nil
	var imageData []byte

	// Check if a file was uploaded
	file, _, err := c.Request.FormFile("image")
	if err == nil {
		// Read the file data into a byte slice
		defer file.Close()
		imageData, err = io.ReadAll(file)
		if err != nil {
			c.JSON(500, gin.H{"error": "Unable to read file"})
			return
		}
	} else if err != http.ErrMissingFile {
		// If the error is not ErrMissingFile, handle it as a file processing error
		c.JSON(400, gin.H{"error": "Error processing file"})
		return
	}

	// Get other form data off req body
	name := c.PostForm("name")
	description := c.PostForm("description")
	price := c.PostForm("price")
	quantity := c.PostForm("quantity")
	color := c.PostForm("color")

	// Convert price and quantity to appropriate types
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid price"})
		return
	}

	quantityInt, err := strconv.Atoi(quantity)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid quantity"})
		return
	}

	// Create a product
	post := models.Product{
		Name:        name,
		Description: description,
		Price:       priceFloat,
		Quantity:    quantityInt,
		Color:       color,
		ImageData:   imageData, // Store the image data
	}

	result := initializers.DB.Create(&post)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": result.Error})
		return
	}

	// Encode image data to Base64 for response
	var imageDataBase64 string
	if post.ImageData != nil {
		imageDataBase64 = base64.StdEncoding.EncodeToString(post.ImageData)
	}

	// Return the product
	c.JSON(200, gin.H{
		"message": "Product created successfully",
		"product": gin.H{
			"id":          post.ID,
			"name":        post.Name,
			"description": post.Description,
			"price":       post.Price,
			"quantity":    post.Quantity,
			"color":       post.Color,
			"imageData":   imageDataBase64, // Return image data as Base64 string
		},
	})
}

// PostsIndex handles fetching all products.
func PostsIndex(c *gin.Context) {
	// Get products
	var posts []models.Product
	if err := initializers.DB.Find(&posts).Error; err != nil {
		// Return 500 if there is an error retrieving products
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}

	// Encode image data to Base64 for each product
	for i := range posts {
		if posts[i].ImageData != nil {
			posts[i].ImageData = []byte(base64.StdEncoding.EncodeToString(posts[i].ImageData))
		}
	}

	// Respond with all products
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// PostsShow handles fetching a single product by ID.
func PostsShow(c *gin.Context) {
	// Get ID from URL
	id := c.Param("id")

	// Get product
	var post models.Product
	if err := initializers.DB.First(&post, id).Error; err != nil {
		// Return 404 if the product is not found
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		// Return 500 if there is a database error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
		return
	}

	// Encode image data to Base64 if it exists
	var imageDataBase64 string
	if post.ImageData != nil {
		imageDataBase64 = base64.StdEncoding.EncodeToString(post.ImageData)
	}

	// Respond with product data
	c.JSON(http.StatusOK, gin.H{
		"post": gin.H{
			"id":          post.ID,
			"name":        post.Name,
			"description": post.Description,
			"price":       post.Price,
			"quantity":    post.Quantity,
			"color":       post.Color,
			"imageData":   imageDataBase64, // Return image data as Base64 string
		},
	})
}

// ProductUpdate handles updating a product, including its image.
func ProductUpdate(c *gin.Context) {
	// Get ID off the URL
	id := c.Param("id")

	// Get data off req body
	var body struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Quantity    int     `json:"quantity"`
		Color       string  `json:"color"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	var imageData []byte
	// Check if a file was uploaded
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			// No file uploaded, set imageData to nil
			imageData = nil
		} else {
			// File upload error
			c.JSON(400, gin.H{"error": "Error processing file"})
			return
		}
	} else {
		// Read the file data into a byte slice
		defer file.Close()
		imageData, err = io.ReadAll(file)
		if err != nil {
			c.JSON(500, gin.H{"error": "Unable to read file"})
			return
		}
	}

	// Find the product we are updating
	var post models.Product
	initializers.DB.First(&post, id)

	// Update the product
	updates := models.Product{
		Name:        body.Name,
		Description: body.Description,
		Price:       body.Price,
		Quantity:    body.Quantity,
		Color:       body.Color,
	}
	if len(imageData) > 0 {
		updates.ImageData = imageData
	}
	initializers.DB.Model(&post).Updates(updates)

	// Respond
	c.JSON(200, gin.H{"message": "Product updated successfully"})
}

// ProductDelete handles deleting a product by ID.
func ProductDelete(c *gin.Context) {
	// Get ID off the URL
	id := c.Param("id")

	// Delete the product
	initializers.DB.Delete(&models.Product{}, id)

	// Respond
	c.JSON(200, gin.H{"message": "Product deleted successfully"})
}

// GetTotalProducts handles fetching the total number of products.
func GetTotalProducts(c *gin.Context) {
	var count int64
	if err := initializers.DB.Model(&models.Product{}).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_products": count})
}

// GetAllProductTypes handles fetching all unique product types.
func GetAllProductTypes(c *gin.Context) {
	var types []struct {
		Type string `json:"type"`
	}

	if err := initializers.DB.Model(&models.Product{}).
		Select("DISTINCT color as type"). // Replace 'color' with the actual column that represents product types
		Scan(&types).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product types"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product_types": types})
}

//getallproducts:

// PostsIndex handles fetching all products.
func GetAllProducts(c *gin.Context) {
	// Get products
	var products []models.Product
	if err := initializers.DB.Find(&products).Error; err != nil {
		// Return 500 if there is an error retrieving products
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}

	// Encode image data to Base64 for each product
	for i := range products {
		if products[i].ImageData != nil {
			products[i].ImageData = []byte(base64.StdEncoding.EncodeToString(products[i].ImageData))
		}
	}

	// Respond with all products
	c.JSON(http.StatusOK, gin.H{"products": products})
}
