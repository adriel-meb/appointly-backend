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
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// RequireAuthMiddleware ensures that requests include a valid JWT
func RequireAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 1️⃣ Check token in cookie
		if cookie, err := c.Cookie("Authorization"); err == nil {
			tokenString = cookie
		} else {
			// 2️⃣ If no cookie, check "Authorization" header with Bearer scheme
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

		// 3️⃣ Parse and validate the JWT signature
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			// Ensure the signing method is HMAC
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

		// 4️⃣ Extract claims (payload) from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: Could not parse claims",
			})
			return
		}

		// 5️⃣ Check expiration time
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

		// 6️⃣ Fetch the user from the database using the subject ("sub") claim
		var user models.User
		if err := db.DB.First(&user, "id = ?", claims["sub"]).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: User not found",
			})
			return
		}

		// 7️⃣ Store the authenticated user in Gin's context for downstream handlers
		c.Set("user", user)

		// 8️⃣ Log the request for monitoring/debugging
		log.Printf("✅ Authenticated request - UserID: %d, Path: %s", user.ID, c.Request.URL.Path)

		// 9️⃣ Continue to the next handler
		c.Next()
	}
}
