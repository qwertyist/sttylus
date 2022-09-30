package session

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type SessionHandler interface {
	createSession(w http.ResponseWriter, r *http.Request)
	getSessions(w http.ResponseWriter, r *http.Request)
	getSession(w http.ResponseWriter, r *http.Request)
	updateSession(w http.ResponseWriter, r *http.Request)
	deleteSession(w http.ResponseWriter, r *http.Request)
	setSessionData(w http.ResponseWriter, r *http.Request)
}

type sessionHandler struct {
	Service SessionService
}

func NewSessionHandler(s SessionService) SessionHandler {
	return &sessionHandler{s}

}
func AddHandlers(r *mux.Router, h SessionHandler) {
	r.HandleFunc("/session", h.setSessionData).Methods("PUT")
	r.HandleFunc("/session", h.createSession).Methods("POST")
	r.HandleFunc("/sessions", h.getSessions).Methods("GET")
	r.HandleFunc("/sessions/{id}", h.getSessions).Methods("GET")
	r.HandleFunc("/session/{id:[0-9]+}", h.getSession).Methods("GET")
	r.HandleFunc("/session/{id:[0-9]+}", h.updateSession).Methods("PUT")
	r.HandleFunc("/session/{id:[0-9]+}", h.deleteSession).Methods("DELETE")
}

func apiHelper(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to madness"))
}

func respondWithJSON(w http.ResponseWriter, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *sessionHandler) createSession(w http.ResponseWriter, r *http.Request) {
	var s Session
	err := json.NewDecoder(r.Body).Decode(&s)

	if err != nil {
		log.Println("decoder:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := h.Service.CreateSession(s)
	if created == nil {
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}
	respondWithJSON(w, created)
}

func (h *sessionHandler) getSessions(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	log.Println("get sessions for id:", id)
	ss, err := h.Service.GetSessions([]string{id})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if ss == nil {
		w.Write([]byte("{}"))
		return
	}
	respondWithJSON(w, ss)
}

func (h *sessionHandler) getSession(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	fmt.Println(id)
	sess, err := h.Service.GetSession(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}
	respondWithJSON(w, sess)
}

func (h *sessionHandler) updateSession(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	fmt.Fprintf(w, "updateSession #%s", id)
}

func (h *sessionHandler) deleteSession(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sess, err := h.Service.DeleteSession(id)
	if sess == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp, err := json.Marshal(sess)
	if err != nil {
		log.Println(err)
	}
	w.Write(resp)
}

func (h *sessionHandler) setSessionData(w http.ResponseWriter, r *http.Request) {
	var s Session
	err := json.NewDecoder(r.Body).Decode(&s)

	if err != nil {
		log.Println("decoder:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("session id:", s.ID)
	log.Println("password:", s.Password)
	log.Println("api token:", s.Token)
	log.Println("api breakout token:", s.Breakout)
	fmt.Fprintf(w, "ok")
}
