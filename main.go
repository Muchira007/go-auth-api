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
	// middleware.SendMailSimpleHTML("Another subject", "./templates/mail.html", []string{"stivmicah@gmail.com"})
	// middleware.SendGoMail("./templates/mail.html")
}

func main() {
	r := gin.Default()

	// Authentication Details
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.POST("/validate", middleware.RequireAuth, controllers.Validate)
	r.POST("/forgot-password", middleware.RequireAuth, controllers.ForgotPassword)
	r.POST("/reset-password", middleware.RequireAuth, controllers.ResetPassword)

	// Product Details (with RequireAuth middleware)
	productRoutes := r.Group("/products")
	productRoutes.Use(middleware.RequireAuth)
	{
		productRoutes.POST("/", controllers.PostsCreate)
		productRoutes.POST("/get-product", controllers.PostsIndex)
		productRoutes.POST("/:id", controllers.PostsShow)
		productRoutes.POST("/Update/:id", controllers.ProductUpdate)
		productRoutes.POST("/delete/:id", controllers.ProductDelete)
	}

	// Sale Details (with RequireAuth middleware)
	r.POST("/sales", middleware.RequireAuth, controllers.RecordSale)

	// Route to download file
	r.POST("/download-excel", middleware.RequireAuth, middleware.DownloadExcelFile)
	r.POST("/download-pdf", middleware.RequireAuth, middleware.DownloadPDF)

	r.Run()
}
