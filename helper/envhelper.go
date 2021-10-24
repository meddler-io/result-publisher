package helper

import "os"

func GetenvStr(key string, defaultValue string) string {
	return getenvStr(key, defaultValue)
}

func getenvStr(key string, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}
