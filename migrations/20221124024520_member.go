package migrations

import (
	"database/sql"
	"gin-n-juice/db"
	"gin-n-juice/utils/types"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

func init() {
	goose.AddMigration(upMember, downMember)
}

type Member struct {
	gorm.Model
	UserId uint       `gorm:"index; not null;"`
	TeamId uint       `gorm:"index; not null;"`
	Role   string     `gorm:"index; not null;"`
	Info   types.JSON `gorm:"serializer:json"`
}

func upMember(tx *sql.Tx) error {
	return db.DB.AutoMigrate(Member{})
}

func downMember(tx *sql.Tx) error {
	return db.DB.Migrator().DropTable(Member{})
}
