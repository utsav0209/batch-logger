package utils

import (
	"fmt"
	"os"
	"strconv"
)

func GetEnvAsInt(name string) int {

	env, err := strconv.Atoi(GetEnvAsString(name))

	if err != nil {
		panic("Environment variable could not be parsed as int")
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
