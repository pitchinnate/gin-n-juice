package main

import (
	"fmt"
	"gin-n-juice/config"
	"gin-n-juice/db"
	"gin-n-juice/routes"
	"github.com/joho/godotenv"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func main() {
	loc, err := time.LoadLocation("UTC")
	if err == nil {
		time.Local = loc
	}

	loadEnv()
	config.SetupEnv()
	db.ConnectDatabase(logger.Silent)
	serve()
}

func loadEnv() {
	directory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = godotenv.Load(fmt.Sprintf("%s/.env", directory))
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
}

func serve() {
	r := routes.SetupRouter()
	r.Run(":8080")
}
