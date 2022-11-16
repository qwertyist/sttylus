package main

import (
	"net/http"

	"github.com/qwertyist/tabula/ws"
)

const debug = false

func main() {
	a := App{pools: make(map[string]*ws.Pool)}
	a.Initialize()
	a.Run(":8888")
	http.ListenAndServe(":8888", nil)
}
