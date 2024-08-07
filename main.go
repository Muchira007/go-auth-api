package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/muchira007/jambo-green-go/controllers"
	"github.com/muchira007/jambo-green-go/initializers"
	"github.com/muchira007/jambo-green-go/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
	// middleware.SendMailSimpleHTML("Another subject", "/templates/mail.html", []string{"stivmicah@gmail.com"})
	// middleware.SendGoMail("./templates/mail.html")
}

func main() {
	r := gin.Default()

	// CORS middleware configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Allow your frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.Static("/uploads", "./uploads")

	// Authentication Details
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.POST("/validate", middleware.RequireAuth, controllers.Validate)
	r.POST("/forgot-password", controllers.ForgotPassword)
	r.POST("/reset-password", controllers.ResetPassword)
	r.POST("/users/phone/:phone_num", controllers.GetUserByPhoneNum)

	// Get all users
	r.POST("/users", middleware.RequireAuth, controllers.GetAllUsers)
	r.POST("/create-users", middleware.RequireAuth, controllers.CreateUser)
	r.POST("/get-users/:id", middleware.RequireAuth, controllers.GetUserByID)
	r.POST("/update-users/:id", middleware.RequireAuth, controllers.UpdateUser)

	// r.POST("/products", controllers.PostsCreate)

	// Product Details (with RequireAuth middleware)
	productRoutes := r.Group("/products")
	productRoutes.Use(middleware.RequireAuth) // Uncomment if authentication is needed for product routes
	{
		productRoutes.POST("", controllers.PostsCreate)
		productRoutes.POST("/get-product", controllers.PostsIndex)
		productRoutes.POST("/:id", controllers.PostsShow)
		productRoutes.POST("/Update/:id", controllers.ProductUpdate)
		productRoutes.POST("/delete/:id", controllers.ProductDelete)
		productRoutes.POST("/product-types", controllers.GetAllProductTypes)
		productRoutes.POST("/total-products", controllers.GetTotalProducts)
		productRoutes.POST("/get-all-products", controllers.GetAllProducts)
	}

	// Sale Details (with RequireAuth middleware)
	r.POST("/sales", middleware.RequireAuth, controllers.RecordSale)
	r.POST("/sales-by-gender", middleware.RequireAuth, controllers.GetSalesByGender)
	r.POST("/sales-by-NatId", middleware.RequireAuth, controllers.GetSalesByNationalID)
	r.POST("/agent-sales", middleware.RequireAuth, controllers.GetAgentSales)

	// Start the server
	r.POST("/download-pdf", middleware.DownloadPDF)
	r.POST("/download-excel", middleware.DownloadExcelFile)
	r.Run()

}
