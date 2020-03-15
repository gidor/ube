package infra

import (
	"os"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func Getenv(key string, defvalue string) string {
	val, found := os.LookupEnv(key)
	if found {
		return val
	} else {
		return defvalue
	}
}

func Setenv(key string, value string) {
	os.Setenv(key, value)
}
