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

var hostname string

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
	var ok bool
	r := routes.SetupRouter()
	hostname, ok = os.LookupEnv("SERVER_HOSTNAME")
	if !ok {
		hostname = "localhost"
	}

	connectionString := fmt.Sprintf("%s:%s", hostname, config.PORT)
	r.Run(connectionString)
}
