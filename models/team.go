package models

import (
	"gin-n-juice/utils/types"
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	OwnerUserId uint       `json:"owner_user_id" binding:"required""`
	Name        string     `json:"name" binding:"required"`
	Info        types.JSON `json:"info" gorm:"serializer:json"`
}
