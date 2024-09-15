package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Env struct {
	URL          string
	NUM_REQUESTS int
	TIMEOUT      int
}

func getProjectPath() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	for {
		if _, err := os.Stat(filepath.Join(basepath, "go.mod")); err == nil {
			return basepath
		}

		newBasePath := filepath.Dir(basepath)
		if newBasePath == basepath {
			break
		}
		basepath = newBasePath
	}
	return ""
}

func LoadEnv() *Env {
	envPath := getProjectPath() + "/.env"
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
		os.Exit(1)
	}

	env := Env{
		URL:          getEnvString("URL"),
		NUM_REQUESTS: getEnvAsInt("NUM_REQUESTS"),
		TIMEOUT:      getEnvAsInt("TIMEOUT"),
	}

	fmt.Printf("Environment variables loaded successfully")
	return &env
}

func getEnvString(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	panic("Environment variable " + key + " is not set")
}

// Helper to read an environment variable into integer
func getEnvAsInt(name string) int {
	valueStr := getEnvString(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	panic("Environment variable " + name + " is not an integer")
}

// Helper to read an environment variable into a bool
func getEnvAsBool(name string) bool {
	valStr := getEnvString(name)
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	panic("Environment variable " + name + " is not a boolean")
}

// Helper to read an environment variable into a string slice
func getEnvAsSlice(name string, sep string) []string {
	valStr := getEnvString(name)

	if valStr == "" {
		panic("Environment variable " + name + " is empty")
	}

	val := strings.Split(valStr, sep)

	return val
}
