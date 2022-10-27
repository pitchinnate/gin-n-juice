package config

import (
	"log"
	"os"
)

var MAILGUN_DOMAIN string
var MAILGUN_PRIVATE_KEY string
var MAILGUN_VALIDATION_KEY string
var ENCRYPT_KEY string
var DB_TYPE string
var DB_CONNECTION_STRING string
var IS_TESTING bool
var PORT string
var DEBUG = true
var EMAIL_FROM string
var PACKAGE_NAME string

func SetupEnv() {
	var ok bool
	MAILGUN_DOMAIN, ok = os.LookupEnv("MAILGUN_DOMAIN")
	if !ok {
		log.Fatalf("missing required env var MAILGUN_DOMAIN")
	}
	MAILGUN_PRIVATE_KEY, ok = os.LookupEnv("MAILGUN_PRIVATE_KEY")
	if !ok {
		log.Fatalf("missing required env var MAILGUN_PRIVATE_KEY")
	}
	MAILGUN_VALIDATION_KEY, ok = os.LookupEnv("MAILGUN_VALIDATION_KEY")
	if !ok {
		log.Fatalf("missing required env var MAILGUN_VALIDATION_KEY")
	}
	ENCRYPT_KEY, ok = os.LookupEnv("ENCRYPT_KEY")
	if !ok {
		log.Fatalf("missing required env var ENCRYPT_KEY")
	}
	DB_TYPE, ok = os.LookupEnv("DB_TYPE")
	if !ok {
		log.Fatalf("missing required env var DB_TYPE")
	}
	DB_CONNECTION_STRING, ok = os.LookupEnv("DB_CONNECTION_STRING")
	if !ok {
		log.Fatalf("missing required env var DB_CONNECTION_STRING")
	}
	EMAIL_FROM, ok = os.LookupEnv("EMAIL_FROM")
	if !ok {
		log.Fatalf("missing required env var EMAIL_FROM")
	}
	PORT, ok = os.LookupEnv("PORT")
	if !ok {
		PORT = "8080"
	}
	PACKAGE_NAME, ok = os.LookupEnv("PACKAGE_NAME")
	if !ok {
		PACKAGE_NAME = "gin-n-juice"
	}

	debugVal, ok := os.LookupEnv("MODE")
	if ok && debugVal == "production" {
		DEBUG = false
	}

	IS_TESTING = false
}

func SetupTestEnv(dbFile string) {
	ENCRYPT_KEY = "testing"
	DB_TYPE = "sqlite"
	DB_CONNECTION_STRING = dbFile
	IS_TESTING = true
	PORT = "8080"
}
