package main

import (
	"os"

	"github.com/joho/godotenv"
)

type envVariables struct {
	taskAuthToken string
  logFile string
}

func getEnvVariables() *envVariables {
	godotenv.Load()
	taskAuthToken := os.Getenv("TASK_AUTH_TOKEN")
  logFile := os.Getenv("TABULA_LOG_FILE")
  if logFile == "" {
    logFile = "/tmp/tabula.log"
  }
	return &envVariables{taskAuthToken,logFile}
}
