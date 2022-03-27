package config

import (
	"himatro-api/internal/util"
	"log"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

func LoginTokenExpSec() int {
	tokenExpSec, err := strconv.Atoi(os.Getenv("LOGIN_TOKEN_EXP_SEC"))

	if err != nil {
		util.LogErr("WARN", "LOGIN_TOKEN_EXP_SEC is not found in the env", err.Error())
		log.Println("Unable to locate token exp sec from .env file, using default value...")
		return 604800 // 7 days
	}

	return tokenExpSec
}

func UpdateAbsentListTokenExpSec() int {
	exp, err := strconv.Atoi(os.Getenv("UPDATE_ABSENT_LIST_TOKEN_EXP_SEC"))

	if err != nil {
		util.LogErr("WARN", "UPDATE_ABSENT_LIST_TOKEN_EXP_SEC is not found in the env", err.Error())
		log.Println("unable to locate absent list token expired sec, using default value...")
		exp = 3600 // 1 hour expiry
	}

	return exp
}

func UpdateAbsentListCookieName() string {
	name := os.Getenv("UPDATE_ABSENT_LIST_COOKIE_NAME")

	if name == "" {
		util.LogErr("WARN", "UPDATE_ABSENT_LIST_COOKIE_NAME is not found in the env", "")
		log.Print("unable to locate update absent list cookie name, using default value...")

		return "UPDATE_ABSENT_LIST_COOKIE"
	}

	return name
}

func JWTSigningKey() string {
	key := os.Getenv("JWT_SECRET_SIGNING_KEY")

	if key == "" {
		util.LogErr("WARN", "JWT_SECRET_SIGNING_KEY is not found in the env", "")
		log.Println("WARNING! No JWT_SECRET_SIGNING_KEY provided. Using the default one is not save at all.")

		return defaultSecretKey
	}

	return key
}
