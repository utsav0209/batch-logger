package utils

import (
	"fmt"
	"os"
	"strconv"
)

func GetEnvAsInt(name string) int {
	val := GetEnvAsString(name)
	env, err := strconv.Atoi(val)

	if err != nil {
		panic(fmt.Sprintf("Environment variable %s with value %s could not be parsed as int", name, val))
	}

	return env
}

func GetEnvAsString(name string) string {
	env, available := os.LookupEnv(name)

	if !available {
		panic(fmt.Sprintf("Environment variable %s not fount", name))
	}

	return env
}
