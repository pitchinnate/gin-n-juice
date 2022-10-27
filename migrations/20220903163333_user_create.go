package migrations

import (
	"database/sql"
	"gin-n-juice/db"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email         string     `json:"email"`
	Password      string     `json:"-"`
	Admin         bool       `json:"admin"`
	EmailVerified *time.Time `json:"email_verified"`
}

func init() {
	goose.AddMigration(upUserCreate, downUserCreate)
}

func upUserCreate(tx *sql.Tx) error {
	return db.DB.AutoMigrate(User{})
}

func downUserCreate(tx *sql.Tx) error {
	return db.DB.Migrator().DropTable(User{})
}
