package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type IdPName string

const (
	GOOGLE IdPName = "GOOGLE"
	LINE   IdPName = "LINE"
)

type APIName string

const (
	OPENAI      APIName = "OPENAI"
	DREAMSTUDIO APIName = "DREAMSTUDIO"
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

func JWTSecret() string {
	return os.Getenv("JWTSECRET")
}

func LogFileName() string {
	return os.Getenv("LOG_FILE_NAME")
}

func StaticFilePath() string {
	return os.Getenv("STATIC_FILE_PATH")
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

func DBSocketType() string {
	return os.Getenv("DB_SOCKET_TYPE")
}

func OAuthRedirectURL() string {
	return os.Getenv("OAUTH_REDIRECT_URL")
}

func ISSURL(idpName IdPName) string {
	return os.Getenv(fmt.Sprintf("%s_ISS_URL", idpName))
}

func ClientId(idpName IdPName) string {
	return os.Getenv(fmt.Sprintf("%s_CLIENT_ID", idpName))
}

func ClientSecret(idpName IdPName) string {
	return os.Getenv(fmt.Sprintf("%s_CLIENT_SECRET", idpName))
}

func APIKey(apiName APIName) string {
	return os.Getenv(fmt.Sprintf("%s_API_KEY", apiName))
}

func SearchEngineId() string {
	return os.Getenv("SEARCH_ENGINE_ID")
}

func GCSServiceKey() string {
	return os.Getenv("GCS_SERVICE_KEY")
}
