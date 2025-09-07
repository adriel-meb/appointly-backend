package controllers

import (
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
func Signup(c *gin.Context) {
	// 1️⃣ Define input structure and bind request body
	type SignupInput struct {
		Name        string  `json:"name" binding:"required"`                               // Required: User's full name
		Email       string  `json:"email" binding:"required,email"`                        // Required: Must be valid email
		Password    string  `json:"password" binding:"required,min=6"`                     // Required: Min 6 chars
		Role        string  `json:"role" binding:"omitempty,oneof=patient provider admin"` // Optional role
		PhoneNumber *string `json:"phone,omitempty"`                                       // Optional phone number
	}

	var input SignupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		// Return validation error
		c.JSON(http.StatusBadRequest, APIResponse{Status: "error", Error: err.Error()})
		return
	}

	// 2️⃣ Hash password for security
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Status: "error", Error: "Failed to hash password"})
		return
	}

	// 3️⃣ Set default role if not provided
	role := input.Role
	if role == "" {
		role = string(models.RolePatient)
	}

	// 4️⃣ Map input to User model
	user := models.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hash),
		Role:         models.UserRole(role),
		PhoneNumber:  input.PhoneNumber,
	}

	// 5️⃣ Save user to database
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Status: "error", Error: err.Error()})
		return
	}

	// 6️⃣ Respond with user info (without password)
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
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*2, "", "", false, true)

	// 6️⃣ Return JSON with token and user info
	c.Header("Authorization", "Bearer "+tokenString)
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Code:    http.StatusOK,
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
