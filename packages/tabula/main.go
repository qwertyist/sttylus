package main

import (
  "log"
	"net/http"
  "os"

	"github.com/qwertyist/tabula/ws"
)

const debug = false

var config *envVariables

func init() {
	config = getEnvVariables()
}

func main() {
  logFile, err := os.OpenFile(config.logFile, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
  if err != nil {
    log.Panic(err)
  }
  defer logFile.Close()

  log.SetOutput(logFile)
  log.SetFlags(log.Lshortfile | log.LstdFlags)

  log.Println("Logging to custom file", config.logFile)

	a := App{pools: make(map[string]*ws.Pool)}
	a.Initialize(config)
	a.Run(":8888")
	http.ListenAndServe(":8888", nil)
}
