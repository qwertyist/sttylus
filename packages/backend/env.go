package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type envVariables struct {
	port           string
	db             string
	dbFileName     string
	serverURL      string
	globalStandard string
	updateURL      string // ...
	mode           string // These should be set by build script
	version        string // ... and not in .env
}

func getEnvVariables(fileName string) *envVariables {
	if fileName == "" {
		godotenv.Load()
	} else {
		godotenv.Load(fileName)
	}

	log.Printf("Running build #%s\n", os.Getenv("STTYLUS_BUILD"))
	port := ":" + os.Getenv("STTYLUS_PORT")
	db := os.Getenv("STTYLUS_DB")
	dbFileName := os.Getenv("STTYLUS_DB_FILENAME")
	serverURL := os.Getenv("STTYLUS_SERVER_URL")
	globalStandard := os.Getenv("STTYLUS_GLOBAL_STANDARD")

	mode := os.Getenv("STTYLUS_MODE")
	version := os.Getenv("STTYLUS_VERSION")
	updateURL := os.Getenv("STTYLUS_UPDATE_URL")

	if port == ":" {
		port = ":80"
	}

	if db == "" {
		db = "boltdb"
	}

	if dbFileName == "" {
		dbFileName = "./sttylus.db"
	}

	if serverURL == "" {
		serverURL = "http://95.179.170.78"
	}

	if globalStandard == "" {
		globalStandard = "ed03f3fc-f552-437b-8103-cdbd7f61a59e"
	}

	mode = "desktop"
	if version == "" {
		version = "0.8.1"
	}

	return &envVariables{
		port:           port,
		db:             db,
		dbFileName:     dbFileName,
		serverURL:      serverURL,
		globalStandard: globalStandard,
		mode:           mode,
		version:        version,
		updateURL:      updateURL,
	}
}
