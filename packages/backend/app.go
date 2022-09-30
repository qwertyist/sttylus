package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/botvid/webapp/repository"
	"github.com/botvid/webapp/session"
	"github.com/botvid/webapp/ws"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	db     repository.Repository
	pools  map[string]*ws.Pool
}

func (a *App) Initialize(r *mux.Router) {
	a.Router = r
	a.Router.HandleFunc("/conn/{id}", func(w http.ResponseWriter, r *http.Request) {
		a.serveWebsocket(w, r)
	})
	a.initializeRoutes()
}

func (a *App) CreatePool(id string) {
	a.pools[id] = ws.NewPool()
	go a.pools[id].Start()
	log.Println("Creating pool for session id:", id)
}

func (a *App) serveWebsocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket endpoint hit")
	params := mux.Vars(r)
	id := params["id"]
	conn, err := ws.Upgrade(w, r)

	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
	}
	_, p, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
	}
	log.Println("ID:", string(p))
	if a.pools[id] == nil {
		a.CreatePool(id)
		conn.WriteJSON(ws.Message{Type: ws.CreateSession})
	}

	client := &ws.Client{
		ID:   string(p),
		Conn: conn,
		Pool: a.pools[id],
	}
	a.pools[id].Register <- client
	client.Read()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, nil))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) initializeRoutes() {
	sessionService := session.NewSessionService(repo)
	sessionHandler := session.NewSessionHandler(sessionService)
	session.AddHandlers(a.Router, sessionHandler)
}
