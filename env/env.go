package env

import (
	"log"
	"os"
)

func Get(key string, nvl string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return nvl
	}
	return value
}

func MustGet(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Fatalf("required ENV variable [%s] is missing", key)
	}
	return value
}
