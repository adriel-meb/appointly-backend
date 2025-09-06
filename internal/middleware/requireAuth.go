package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/adriel-meb/appointly-backend/internal/db"
	"github.com/adriel-meb/appointly-backend/internal/models"
	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)

// RequireAuthMiddleware checks for a valid JWT token in the Authorization header or cookie
func RequireAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 1. Try cookie first
		if cookie, err := c.Cookie("Authorization"); err == nil {
			tokenString = cookie
		} else {
			// 2. If no cookie, try Authorization header
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Unauthorized: No token provided",
				})
				return
			}
		}

		// 3. Parse and validate JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: Invalid token",
			})
			return
		}

		// 4. Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: Could not parse claims",
			})
			return
		}

		// 5. Check expiration
		if exp, ok := claims["exp"].(float64); ok {
			if float64(time.Now().Unix()) > exp {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Unauthorized: Token expired",
				})
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: Missing expiration claim",
			})
			return
		}

		// 6. Fetch the user
		var user models.User
		if err := db.DB.First(&user, "id = ?", claims["sub"]).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: User not found",
			})
			return
		}

		// 7. Attach user to context
		c.Set("user", user)

		// 8. Log and continue
		log.Printf("Authenticated request from user ID: %d, Path: %s", user.ID, c.Request.URL.Path)
		c.Next()
	}
}
