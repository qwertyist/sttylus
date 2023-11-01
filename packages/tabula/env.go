package main

import (
	"os"

	"github.com/joho/godotenv"
)

type envVariables struct {
	taskAuthToken string
}

func getEnvVariables() *envVariables {
	godotenv.Load()
	taskAuthToken := os.Getenv("TASK_AUTH_TOKEN")
	return &envVariables{taskAuthToken}
}
