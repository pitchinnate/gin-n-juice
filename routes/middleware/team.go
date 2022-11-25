package middleware

import (
	"gin-n-juice/db"
	"gin-n-juice/models"
	"gin-n-juice/utils/input"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Team() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("user")
		userObj := user.(models.User)
		id, err := input.GetPathInt(c, "id")
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Team not found"})
			c.Abort()
			return
		}
		team := models.Team{}
		if err := db.DB.First(id).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("team", team)
		if userObj.Admin {
			c.Next()
		} else {
			var members []models.Member
			if err := db.DB.Where("team_id = ? and user_id = ?", team.ID, userObj.ID).First(&members).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
			c.Next()
		}
	}
}
