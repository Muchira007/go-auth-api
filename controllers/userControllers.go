package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/muchira007/jambo-green-go/initializers"
	"github.com/muchira007/jambo-green-go/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	// signup logic
	var body struct {
		NationalID int    `json:"national_id" binding:"required"`
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
	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})
}

func Login(c *gin.Context) {
	// login logic
	var body struct {
		PhoneNum string `json:"phone_num" binding:"required"`
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
	initializers.DB.Where("phone_num = ?", body.PhoneNum).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid phone or password",
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

	user.(models.User).
		c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
