package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	if os.Getenv("GO_ENV") == "dev" {
		if err := godotenv.Load("/usr/src/app/.env"); err != nil {
			log.Fatalf("Failed to load env file: %s", err)
		}
	}
}

func GoEnv() string {
	return os.Getenv("GO_ENV")
}

func Port() uint16 {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Failed to parse string to int: %s", err)
	}

	return uint16(port)
}

func APIDomain() string {
	return os.Getenv("API_DOMAIN")
}

func FrontEndURL() string {
	return os.Getenv("FE_URL")
}

func Secret() string {
	return os.Getenv("SECRET")
}

func LogFileName() string {
	return os.Getenv("LOG_FILE")
}

func MySQLDBName() string {
	return os.Getenv("MYSQL_DB")
}

func MySQLUser() string {
	return os.Getenv("MYSQL_USER")
}

func MySQLPass() string {
	return os.Getenv("MYSQL_PASS")
}

func MySQLHost() string {
	return os.Getenv("MYSQL_HOST")
}
