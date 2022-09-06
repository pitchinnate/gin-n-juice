package migrations

import (
	"database/sql"
	"gin-n-juice/db"
	"github.com/pressly/goose/v3"
	"time"
)

func init() {
	goose.AddMigration(upPasswordResetCreate, downPasswordResetCreate)
}

type PasswordReset struct {
	Email     string `json:"email" gorm:"index"`
	Token     string `json:"token"`
	CreatedAt time.Time
}

func upPasswordResetCreate(tx *sql.Tx) error {
	return db.DB.AutoMigrate(PasswordReset{})
}

func downPasswordResetCreate(tx *sql.Tx) error {
	return db.DB.Migrator().DropTable(PasswordReset{})
}
