package env

import (
	"github.com/joho/godotenv"

	"os"
	"strconv"
)

func LoadDotEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}

func GetString(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return value
}

func GetInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	valAsInt, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return valAsInt
}

func GetBool(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return boolVal
}
