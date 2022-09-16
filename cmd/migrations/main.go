package main

import (
	"flag"
	"fmt"
	"gin-n-juice/config"
	"gin-n-juice/db"
	_ "gin-n-juice/migrations"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	flags   = flag.NewFlagSet("goose", flag.ExitOnError)
	testing = flags.Bool("testing", false, "running in test")
)

func main() {
	loc, err := time.LoadLocation("UTC")
	if err == nil {
		time.Local = loc
	}

	if len(os.Args) == 1 {
		log.Fatalf("No commands sent to migrate")
	}

	flags.Parse(os.Args[1:])
	args := flags.Args()

	if !*testing {
		loadEnv()
		config.SetupEnv()
	} else {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		config.SetupTestEnv(fmt.Sprintf("%s/test.db", exPath))
	}
	db.ConnectDatabase(logger.Info)

	command := "status"
	if len(args) > 0 {
		command = args[0]
	}

	log.Printf("Set dialect to: %s", config.DB_TYPE)
	if config.DB_TYPE == "sqlite" {
		if err := goose.SetDialect("sqlite3"); err != nil {
			log.Fatalf("goose error setting dialect: %v", err)
		}
	} else if config.DB_TYPE == "mysql" {
		if err := goose.SetDialect("mysql"); err != nil {
			log.Fatalf("goose error setting dialect: %v", err)
		}
	} else if config.DB_TYPE == "postgres" {
		if err := goose.SetDialect("postgres"); err != nil {
			log.Fatalf("goose error setting dialect: %v", err)
		}
	} else {
		log.Fatalf("Using an unsupported database connection type: %s", config.DB_TYPE)
	}

	directory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	migrationDirectory := fmt.Sprintf("%s/migrations", directory)

	log.Println("Running migration script")
	if command != "" {
		log.Printf("Running: %s with values %v", command, args[1:])
		if err := goose.Run(command, db.DbHandle, migrationDirectory, args[1:]...); err != nil {
			log.Fatalf("goose run error: %v", err)
		}
		return
	}
}

func loadEnv() {
	directory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = godotenv.Load(fmt.Sprintf("%s/.env", directory))
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
