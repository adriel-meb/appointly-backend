package controllers

import (
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"

	"github.com/adriel-meb/appointly-backend/internal/db"
	"github.com/adriel-meb/appointly-backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ---------------------- API RESPONSE ---------------------- //

// ---------------------- AUTH HANDLERS ---------------------- //

// Signup handles POST /signup
// Creates a new user account
// Signup handles POST /signup
func Signup(c *gin.Context) {
	type SignupInput struct {
		Name             string  `json:"name" binding:"required"`
		Email            string  `json:"email" binding:"required,email"`
		Password         string  `json:"password" binding:"required,min=6"`
		Role             string  `json:"role" binding:"omitempty,oneof=patient provider admin"`
		PhoneNumber      *string `json:"phone,omitempty"`
		SpecializationID *uint   `json:"specialization_id,omitempty"` // FK to specialization
		Bio              string  `json:"bio,omitempty"`               // Only for provider
	}

	var input SignupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Status: "error", Error: err.Error()})
		fmt.Println(err.Error())
		return
	}

	// Check if email exists
	var existing models.User
	if err := db.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, APIResponse{Status: "error", Error: "Email already registered"})
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Status: "error", Error: "Failed to hash password"})
		return
	}

	// Default role
	role := input.Role
	if role == "" {
		role = string(models.RolePatient)
	}

	user := models.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hash),
		Role:         models.UserRole(role),
		PhoneNumber:  input.PhoneNumber,
	}

	// Transaction: user + provider (if needed)
	err = db.DB.Transaction(func(tx *gorm.DB) error {
		// Save user
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// If provider, also save specialization + bio
		if role == string(models.RoleProvider) {
			// Ensure specialization is provided
			if input.SpecializationID == nil {
				return fmt.Errorf("specialization_id is required for providers")
			}

			provider := models.Provider{
				UserID:           user.ID,
				SpecializationID: *input.SpecializationID,
				Bio:              input.Bio,
			}
			if err := tx.Create(&provider).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Status: "error", Error: err.Error()})
		return
	}

	fmt.Println("User created successfully")
	// Success response
	c.JSON(http.StatusCreated, APIResponse{
		Status:  "success",
		Message: "User created successfully",
		Data: gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
			"phone": user.PhoneNumber,
		},
	})
}

// Login handles POST /login
// Authenticates a user and returns a JWT
func Login(c *gin.Context) {
	// 1️⃣ Bind request body
	var input struct {
		Email    string `json:"email" binding:"required,email"` // Must be valid email
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Status: "error", Error: err.Error()})
		return
	}

	// 2️⃣ Fetch user from database by email
	var user models.User
	if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{Status: "error", Error: "Invalid email or password"})
		return
	}

	// 3️⃣ Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{Status: "error", Error: "Invalid email or password"})
		return
	}

	// 4️⃣ Create JWT token (signed with secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID, // Subject = User ID
		"email": user.Email,
		"role":  user.Role,
		"iat":   time.Now().Unix(),                          // Issued at
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(), // Expiration 30 days
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Status: "error", Error: "Failed to generate token"})
		return
	}

	// 5️⃣ Set auth cookie for browser clients (optional)
	//c.SetSameSite(http.SameSiteNoneMode) // 👈 Allow cross-site
	c.SetSameSite(http.SameSiteLaxMode) // works for localhost
	c.SetCookie(
		"jwt_token", // name
		tokenString, // value
		3600*24*30,  // maxAge: 30 days
		"/",         // path
		"localhost", // domain (use your frontend domain if cross-site)
		false,
		//gin.Mode() == gin.ReleaseMode, // secure only in prod
		true, // httpOnly
	)

	// 6️⃣ Return JSON with token and user info
	c.Header("Authorization", "Bearer "+tokenString)
	c.JSON(200, APIResponse{
		Status:  "success",
		Message: "Login successful",
		Data: gin.H{
			"token": tokenString,
			"user": gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
				"role":  user.Role,
			},
		},
	})
}

// Validate handles GET /validate
// Checks if the user is logged in via middleware
func Validate(c *gin.Context) {
	// 1️⃣ Retrieve user from context (set by JWT middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{Status: "error", Error: "Unauthorized"})
		return
	}

	// 2️⃣ Type assertion
	u, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, APIResponse{Status: "error", Error: "Invalid user type"})
		return
	}

	// 3️⃣ Return logged-in user info
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "You are logged in",
		Data:    u,
	})
}

func GetProfile(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// Logout handles POST /logout
// Clears the JWT cookie to log the user out
func Logout(c *gin.Context) {
	fmt.Println("----------- logout ------------")

	// Clear the JWT cookie by setting MaxAge to -1
	c.SetCookie(
		"jwt_token", // cookie name
		"",          // empty value
		-1,          // MaxAge -1 deletes the cookie
		"/",         // path
		"localhost", // domain (your frontend domain)
		false,       // secure false for localhost
		true,        // httpOnly
	)

	// Optionally, also clear Authorization header if used
	c.Header("Authorization", "")

	// Return a success response
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Logged out successfully",
	})
}
