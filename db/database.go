package db

import (
	"database/sql"
	"gin-n-juice/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

var DB *gorm.DB
var DbHandle *sql.DB

func ConnectDatabase(logLevel logger.LogLevel) {
	var database *gorm.DB
	var err error

	gormConfig := gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	if config.DB_TYPE == "sqlite" {
		database, err = gorm.Open(sqlite.Open(config.DB_CONNECTION_STRING), &gormConfig)
	} else if config.DB_TYPE == "mysql" {
		database, err = gorm.Open(mysql.Open(config.DB_CONNECTION_STRING), &gormConfig)
	} else if config.DB_TYPE == "postgres" {
		database, err = gorm.Open(postgres.Open(config.DB_CONNECTION_STRING), &gormConfig)
	} else {
		log.Fatalf("Using an unsupported database connection type: %s", config.DB_TYPE)
	}

	if err != nil {
		log.Print("Connection Error: ", err)
		panic("Failed to connect to database!")
	}

	DB = database
	DbHandle, err = database.DB()
	if err != nil {
		panic("Error getting db handle")
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	DbHandle.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	DbHandle.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	DbHandle.SetConnMaxLifetime(time.Hour)
}

func Disconnect() {
	DbHandle.Close()
}
