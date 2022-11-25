package migrations

import (
	"database/sql"
	"gin-n-juice/db"
	"gin-n-juice/utils/types"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email         string `gorm:"index;not null;unique"`
	Password      string
	Admin         bool `gorm:"index"`
	EmailVerified *time.Time
	Info          types.JSON `gorm:"serializer:json"`
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
