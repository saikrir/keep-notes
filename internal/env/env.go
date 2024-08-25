package env

import (
	"os"
	"strconv"
)

func GetEnvValAsString(key string) string {
	val := os.Getenv(key)

	if len(val) == 0 {
		panic(key + " does not exist")
	}

	return val
}

func GetEnvValAsNumber(key string) int {
	val := GetEnvValAsString(key)

	retVal, err := strconv.Atoi(val)
	if err != nil {
		panic("Failed to convert to number " + key)
	}
	return retVal
}
