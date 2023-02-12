package migrations

import (
	"database/sql"
	"gin-n-juice/db"
	"gin-n-juice/utils/types"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	OwnerUserId uint       `gorm:"index; not null;"`
	Name        string     `gorm:"not null"`
	Info        types.JSON `gorm:"serializer:json"`
}

func init() {
	goose.AddMigration(upTeams, downTeams)
}

func upTeams(tx *sql.Tx) error {
	return db.DB.AutoMigrate(Team{})
}

func downTeams(tx *sql.Tx) error {
	return db.DB.Migrator().DropTable(Team{})
}
