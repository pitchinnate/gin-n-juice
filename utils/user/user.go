package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"tagdeploy/models"
)

func GetUserFromContext(c *gin.Context) (*models.User, error) {
	userObj, ok := c.Get("user")
	if !ok {
		return nil, fmt.Errorf("no user in context")
	}
	user := userObj.(models.User)
	return &user, nil
}
