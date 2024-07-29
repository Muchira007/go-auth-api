package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/muchira007/jambo-green-go/initializers"
	"github.com/muchira007/jambo-green-go/models"
)

func PostsCreate(c *gin.Context) {

	//Get data off req body
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	//create a product
	post := models.Product{Name: body.Description, Description: body.Description}

	result := initializers.DB.Create(&post)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": result.Error,
		})
		return
	}
	//return the product
	c.JSON(200, gin.H{
		"message": "PostsCreate",
	})
}

func PostsIndex(c *gin.Context) {
	//get products
	var posts []models.Product
	initializers.DB.Find(&posts)
	//repond
	c.JSON(200, gin.H{
		"posts": posts,
	})
}

func PostsShow(c *gin.Context) {
	//GET ID OFF URL
	id := c.Param("id")

	//get product
	var post models.Product
	initializers.DB.First(&post, id)
	//respond
	c.JSON(200, gin.H{
		"post": post,
	})
}

func ProductUpdate(c *gin.Context) {
	//get id off the URL
	id := c.Param("id")
	//get data off req body
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	c.Bind(&body)

	//Find the product were updating
	var post models.Product
	initializers.DB.First(&post, id)

	//update the product
	initializers.DB.Model(&post).Updates(models.Product{
		Name:        body.Name,
		Description: body.Description,
	})
	//respond
	c.JSON(200, gin.H{
		"message": "ProductUpdate",
	})

}

func ProductDelete(c *gin.Context) {
	//get id off the URL
	id := c.Param("id")

	//Delete the posts
	initializers.DB.Delete(&models.Product{}, id)

	//respond
	c.JSON(200, gin.H{
		"message": "ProductDeleted",
	})
}
