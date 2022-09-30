package main

import (
	"net/http"

	"github.com/qwertyist/tabula/ws"
)

func main() {
	a := App{pools: make(map[string]*ws.Pool)}
	a.Initialize()
	a.Run(":8080")
	http.ListenAndServe(":8080", nil)
}
