package teams

import (
	"gin-n-juice/db"
	"gin-n-juice/models"
	"gin-n-juice/utils/input"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Delete(c *gin.Context) {
	id, err := input.GetPathInt(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var team models.Team
	if err := db.DB.First(&team, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err = db.DB.Delete(&team).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "User Deleted")
}
