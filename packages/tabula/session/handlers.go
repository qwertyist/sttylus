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
	createUser(w http.ResponseWriter, r *http.Request)
	getUser(w http.ResponseWriter, r *http.Request)
	getUsers(w http.ResponseWriter, r *http.Request)
}

type sessionHandler struct {
	Service SessionService
}

func NewSessionHandler(s SessionService) SessionHandler {
	return &sessionHandler{s}

}
func AddHandlers(r *mux.Router, h SessionHandler) {
	r.HandleFunc("/user", h.createUser).Methods("POST")
	r.HandleFunc("/user/", h.getUser).Methods("GET")
	r.HandleFunc("/user/{id}", h.getUser).Methods("GET")
	r.HandleFunc("/users", h.getUsers).Methods("GET")
	r.HandleFunc("/session", h.createSession).Methods("POST")
	r.HandleFunc("/sessions", h.getSessions).Methods("GET")
	r.HandleFunc("/sessions/{id}", h.getSessions).Methods("GET")
	r.HandleFunc("/session/{id:[0-9]+}", h.getSession).Methods("GET")
	r.HandleFunc("/session/{id:[0-9]+}/{action}", h.updateSession).Methods("PUT")
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
	printRequestBody(r)
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
	action := params["action"]
	var s Session
	err := json.NewDecoder(r.Body).Decode(&s)

	if err != nil {
		log.Println("decoder:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if id != s.ID {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if action == "bind" {

	}
	sess, err := h.Service.UpdateSession(s)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}
	respondWithJSON(w, sess)

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

func (h *sessionHandler) getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var request User
	if params["id"] != "" {
		request.ID = params["id"]
	}
	phone := r.FormValue("phone")
	if phone != "" {
		request.Phone = phone
	}
	email := r.FormValue("email")
	if email != "" {
		request.Email = email
	}
	u := h.Service.GetUser(request)
	if u != nil {
		response, err := json.Marshal(u)
		if err != nil {
			log.Panic("getUser marshal response", err)
		}
		w.Write(response)
		return
	}
}
func (h *sessionHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	uu := h.Service.GetUsers()
	if uu == nil {
		http.Error(w, "{ \"error\": \"no users found\"}", http.StatusBadRequest)
		return
	}
	response, err := json.Marshal(uu)
	if err != nil {
		log.Println("decoder:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(response)
}
func (h *sessionHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var request User

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println("decoder:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u := h.Service.GetUser(request)
	if u != nil {
		log.Println("user already exists")
		http.Error(w, "{ \"error\": \"user already exists\"}", http.StatusInternalServerError)
		return
	}

	created := h.Service.CreateUser(request)
	if created == nil {
		log.Println("creating user failed")
		http.Error(w, "couldn't create user", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, created)
}

func (h *sessionHandler) deleleUser(w http.ResponseWriter, r *http.Request) {

}
