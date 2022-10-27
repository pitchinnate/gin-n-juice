package user

import (
	"fmt"
	"gin-n-juice/models"
	"github.com/gin-gonic/gin"
)

func GetUserFromContext(c *gin.Context) (*models.User, error) {
	userObj, ok := c.Get("user")
	if !ok {
		return nil, fmt.Errorf("no user in context")
	}
	user := userObj.(models.User)
	return &user, nil
}
