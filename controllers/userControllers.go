	package controllers

	import (
		"errors"
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

	// Function to generate JWT token
	func generateToken(user models.User) (string, error) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.ID,
			"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		})

		secret := os.Getenv("SECRET")
		if secret == "" {
			return "", errors.New("Secret key not set")
		}

		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			return "", err
		}
		return tokenString, nil
	}

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

		hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to hash password",
			})
			return
		}

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
				"error": result.Error.Error(), // Detailed error message
			})
			return
		}

		tokenString, err := generateToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate token",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User created successfully",
			"token":   tokenString,
			"user":    user,
		})
	}

	func Login(c *gin.Context) {
		var body struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.Bind(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		// Get user
		var user models.User
		initializers.DB.Where("email = ?", body.Email).First(&user)

		if user.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid email or password",
			})
			return
		}

		// Compare passwords
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid password",
			})
			return
		}

		// Generate token
		tokenString, err := generateToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate token",
			})
			return
		}

		// Set cookie
		c.SetSameSite(http.SameSiteDefaultMode)
		c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

		// Respond
		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   tokenString,
			"user":    user,
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

	// GetUserByPhoneNum retrieves a user by their phone number
	func GetUserByPhoneNum(c *gin.Context) {
		// Get the phone number from the URL parameters
		phoneNum := c.Param("phone_num")

		// Find the user by phone number
		var user models.User
		result := initializers.DB.First(&user, "phone_num = ?", phoneNum)

		if result.Error != nil {
			if result.Error.Error() == "record not found" {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "User not found",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to retrieve user",
				})
			}
			return
		}

		// Respond with the user details
		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}

	// CreateUser creates a user without sending back a token
	func CreateUser(c *gin.Context) {
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
		c.JSON(http.StatusOK, gin.H{
			"message": "User created successfully",
			"user":    user,
		})
	}

	// GetUserByID retrieves a user by their ID
	func GetUserByID(c *gin.Context) {
		// Get the user ID from the URL parameters
		ID := c.Param("id")

		// Find the user by ID
		var user models.User
		result := initializers.DB.First(&user, ID)
		if result.Error != nil {
			if result.Error.Error() == "record not found" {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "User not found",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to retrieve user",
				})
			}
			return
		}

		// Respond with the user details
		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}

	// UpdateUser updates a user's information
	func UpdateUser(c *gin.Context) {
		// Get the user ID from the URL parameters
		userID := c.Param("id")

		// Bind the request body to a struct with pointer fields
		var body struct {
			NationalID *int    `json:"national_id"`
			Email      *string `json:"email"`
			FirstName  *string `json:"first_name"`
			SecondName *string `json:"second_name"`
			SurName    *string `json:"sur_name"`
			PhoneNum   *string `json:"phone_num"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		// Find the user by ID
		var user models.User
		result := initializers.DB.First(&user, userID)
		if result.Error != nil {
			if result.Error.Error() == "record not found" {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "User not found",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to retrieve user",
				})
			}
			return
		}

		// Update the user's information only if fields are provided
		if body.NationalID != nil {
			user.NationalID = *body.NationalID
		}
		if body.Email != nil {
			user.Email = *body.Email
		}
		if body.FirstName != nil {
			user.FirstName = *body.FirstName
		}
		if body.SecondName != nil {
			user.SecondName = *body.SecondName
		}
		if body.SurName != nil {
			user.SurName = *body.SurName
		}
		if body.PhoneNum != nil {
			user.PhoneNum = *body.PhoneNum
		}

		// Save the updated user
		if result := initializers.DB.Save(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update user",
			})
			return
		}

		// Respond with the updated user details
		c.JSON(http.StatusOK, gin.H{
			"message": "User updated successfully",
			"user":    user,
		})
	}
