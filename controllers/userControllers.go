package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/muchira007/jambo-green-go/initializers"
	"github.com/muchira007/jambo-green-go/middleware"
	"github.com/muchira007/jambo-green-go/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	var body struct {
		NationalID int    `json:"national_id" binding:"required"`
		Email      string `json:"email" binding:"required"`
		FirstName  string `json:"first_name" binding:"required"`
		SecondName string `json:"second_name" binding:"required"`
		SurName    string `json:"sur_name" binding:"required"`
		PhoneNum   string `json:"phone_num" binding:"required"`
		Password   string `json:"password" binding:"required"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create the user
	user := models.User{
		NationalID: body.NationalID,
		Email:      body.Email,
		FirstName:  body.FirstName,
		SecondName: body.SecondName,
		SurName:    body.SurName,
		PhoneNum:   body.PhoneNum,
		Password:   string(hash),
	}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func Login(c *gin.Context) {
	// login logic
	var body struct {
		// PhoneNum string `json:"phone_num" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}
	// get user
	var user models.User
	initializers.DB.Where("email = ?", body.Email).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email or password",
		})
		return
	}
	// compare passwords
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid password",
		})
		return
	}

	// Generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the encoded token as a string using the secret key
	secret := os.Getenv("SECRET")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Secret key not set",
		})
		return
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to generate token",
		})
		return
	}
	// respond
	//send using cookie
	c.SetSameSite(http.SameSiteDefaultMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		// "token":   tokenString,
		"message": "Login successful",
	})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}

func ForgotPassword(c *gin.Context) {
	var body struct {
		Email string `json:"email" binding:"required"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Find the user by email
	var user models.User
	result := initializers.DB.Where("email = ?", body.Email).First(&user)

	if result.Error != nil || user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Generate reset token
	resetToken := uuid.New().String()

	// Update the user record with the reset token and expiration time
	user.ResetToken = resetToken
	user.ResetTokenExpiry = time.Now().Add(time.Hour * 1) // Token valid for 1 hour
	initializers.DB.Save(&user)

	// Send reset email
	resetLink := "https://example.com/reset-password?token=" + resetToken
	subject := "Password Reset Request"
	bodyText := "Please use the following link to reset your password: " + resetLink
	middleware.SendGoMail(body.Email, subject, bodyText)

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset email sent",
	})
}
func ResetPassword(c *gin.Context) {
	var body struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Find the user by reset token
	var user models.User
	result := initializers.DB.Where("reset_token = ? AND reset_token_expiry > ?", body.Token, time.Now()).First(&user)

	if result.Error != nil || user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid or expired token",
		})
		return
	}

	// Hash the new password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash new password",
		})
		return
	}

	// Update the user password
	user.Password = string(hash)
	// Clear the reset token and its expiry time
	user.ResetToken = ""
	user.ResetTokenExpiry = time.Time{}
	initializers.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}

// GetAllUsers retrieves all users
func GetAllUsers(c *gin.Context) {
	var users []models.User
	result := initializers.DB.Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}
