package main

import (
	"github.com/gin-gonic/gin"
	"github.com/muchira007/jambo-green-go/controllers"
	"github.com/muchira007/jambo-green-go/initializers"
	"github.com/muchira007/jambo-green-go/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	//Authentication Details
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.POST("/validate", middleware.RequireAuth, controllers.Validate)

	//Product Details
	r.POST("/products", controllers.PostsCreate)
	r.POST("/products/get-product", controllers.PostsIndex)
	r.POST("/products/:id", controllers.PostsShow)
	r.POST("/products/Update/:id", controllers.ProductUpdate)
	r.POST("/products/delete/:id", controllers.ProductDelete)

	r.Run()
}
