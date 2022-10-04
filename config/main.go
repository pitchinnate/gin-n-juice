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

	PORT, ok = os.LookupEnv("PORT")
	if !ok {
		log.Printf("PORT environment variable was not found, setting port to 8080...")
		PORT = "8080"
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
