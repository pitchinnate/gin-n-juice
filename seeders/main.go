package seeders

import (
	"gin-n-juice/db"
	"gin-n-juice/models"
	"time"
)

func RunAllSeeders() {
	now := time.Now()

	adminUser := models.User{}
	adminUser.Email = "admin@test.com"
	adminUser.Password = "testing1"
	adminUser.EmailVerified = &now
	adminUser.Admin = true
	db.DB.Create(&adminUser)

	normalUser := models.User{}
	normalUser.Email = "user@test.com"
	normalUser.Password = "testing1"
	normalUser.EmailVerified = &now
	normalUser.Admin = false
	db.DB.Create(&normalUser)

	team := models.Team{}
	team.OwnerUserId = normalUser.ID
	team.Name = "Testing"
	db.DB.Create(&team)

	member := models.Member{}
	member.UserId = adminUser.ID
	member.TeamId = team.ID
	member.Role = "Admin"
	db.DB.Create(&member)

	member2 := models.Member{}
	member2.UserId = normalUser.ID
	member2.TeamId = team.ID
	member2.Role = "Viewer"
	db.DB.Create(&member2)
}
