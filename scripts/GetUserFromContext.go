package scripts

import (
	"github.com/adriel-meb/appointly-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func GetUserFromContext(c *gin.Context) (models.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return models.User{}, false
	}
	u, ok := user.(models.User)
	return u, ok
}
 