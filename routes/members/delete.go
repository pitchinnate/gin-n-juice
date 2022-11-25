package members

import (
	"gin-n-juice/db"
	"gin-n-juice/models"
	"gin-n-juice/utils/input"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Delete(c *gin.Context) {
	id, err := input.GetPathInt(c, "mid")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var member models.Member
	if err := db.DB.First(&member, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err = db.DB.Delete(&member).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "User Deleted")
}
