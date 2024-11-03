package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
  "sync"

	"github.com/gorilla/mux"
	"github.com/qwertyist/tabula/repo"
	"github.com/qwertyist/tabula/session"
	"github.com/qwertyist/tabula/ws"
)

type App struct {
	Router *mux.Router
	db     repo.Repository
	pools  map[string]*ws.Pool
  mu     sync.Mutex
}

func (a *App) Initialize(config *envVariables) {
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/conn/{id}", func(w http.ResponseWriter, r *http.Request) {
		a.serveWebsocket(w, r)
	})
	a.initializeRoutes(config)
}

func (a *App) CreatePool(id string) {
  a.mu.Lock()
	a.pools[id] = ws.NewPool(id)
  a.mu.Unlock()
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
	var interpreter bool

	if string(p) == "interpreter" {
		interpreter = true
	}

	if a.pools[id] == nil {
		a.CreatePool(id)
		conn.WriteJSON(ws.Message{Type: ws.CreateSession})
	}

	client := &ws.Client{
		ID:          string(p),
		Interpreter: interpreter,
		Conn:        conn,
		Pool:        a.pools[id],
	}
	a.pools[id].Register <- client
	client.Read()
}

func (a *App) Run(addr string) {
	http.Handle("/", accessControl(a.Router))
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

func (a *App) initializeRoutes(config *envVariables) {
	a.db = repo.NewRepository(repo.OpenBoltDB("bolt.db"))
	a.Router.HandleFunc("/api", apiHelper).Methods("GET", "POST")

	sessionService := session.NewSessionService(a.db, a.pools)
	sessionHandler := session.NewSessionHandler(sessionService)
	sessionService.SetAuthToken(config.taskAuthToken)
	session.AddHandlers(a.Router, sessionHandler)
}

func apiHelper(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to madness"))
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		///origin := req.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Id-Token")
		//w.Header().Set("Access-Control-Allow-Origin", "*")

		if req.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, req)
	})
}
