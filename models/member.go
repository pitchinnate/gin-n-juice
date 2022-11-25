package models

import (
	"gin-n-juice/utils/types"
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	UserId uint       `json:"user_id" binding:"required"`
	TeamId uint       `json:"team_id" binding:"required"`
	Role   string     `json:"role" binding:"required"`
	Info   types.JSON `json:"info" gorm:"serializer:json"`
}
