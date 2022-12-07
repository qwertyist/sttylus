package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type UserHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	GetUsers(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	UpdatePassword(w http.ResponseWriter, r *http.Request)
	ResetPassword(w http.ResponseWriter, r *http.Request)

	Identify(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	RegisterWeb(w http.ResponseWriter, r *http.Request)
	IsRegistered(w http.ResponseWriter, r *http.Request)
	Auth(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)

	GetSettings(w http.ResponseWriter, r *http.Request)
	UpdateSettings(w http.ResponseWriter, r *http.Request)
	SyncUser(w http.ResponseWriter, r *http.Request)
	SyncUserByID(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	userService UserService
	mode        string
}

//NewUserHandler is...
func NewUserHandler(userService UserService, mode string) UserHandler {
	return &userHandler{
		userService,
		mode,
	}
}

//Endpoints adds User methods to the API router
func Endpoints(r *mux.Router, h UserHandler) {
	r.HandleFunc("/user", h.CreateUser).Methods("POST")
	r.HandleFunc("/users", h.GetUsers).Methods("GET")
	r.HandleFunc("/user/{id}", h.GetUser).Methods("GET")
	r.HandleFunc("/user/{id}", h.UpdateUser).Methods("PUT")
	r.HandleFunc("/user/{id}", h.DeleteUser).Methods("DELETE")

	r.HandleFunc("/settings", h.UpdateSettings).Methods("POST")
	r.HandleFunc("/settings", h.GetSettings).Methods("GET")

	r.HandleFunc("/auth", h.Auth).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")
	r.HandleFunc("/register", h.RegisterWeb).Methods("POST")
	r.HandleFunc("/registered/{email}", h.IsRegistered).Methods("GET")
	r.HandleFunc("/password", h.UpdatePassword).Methods("POST")
	r.HandleFunc("/password", h.ResetPassword).Methods("PUT")

	r.HandleFunc("/sync", h.SyncUser).Methods("POST")
}

func LogUnauthorizedRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Unauthorized request from: %s\n", r.RemoteAddr)
	w.WriteHeader(http.StatusUnauthorized)
	return
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	id := r.Header.Get("X-Id-Token")
	if id == "" {
		log.Printf("Created by unknown user, aborting\n")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	auth := h.userService.Auth(id, UserAdmin)
	if auth != nil {
		log.Printf("Someone tried to create a user %s\n", id)
		LogUnauthorizedRequest(w, r)
		return
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		msg := fmt.Sprintf("%s - /api/user - CreateUser couldn't decode request: %s\n", r.Method, err)
		log.Print(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("POST /api/user - CreateUser couldn't decode request:" + err.Error()))
		return
	}
	if user.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No name provided"))
		return
	}

	u, err := h.userService.CreateUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
	json, err := json.MarshalIndent(u, "", " ")
	w.WriteHeader(http.StatusOK)
	w.Write(json)

}

func (h *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id_token := r.Header.Get("X-Id-Token")
	if id_token == "" {
		log.Printf("handler|GetUser no x-id-token provided\n")
		LogUnauthorizedRequest(w, r)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	u, err := h.userService.GetUser(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request failed: " + err.Error()))
		return
	}
	u.PasswordHash = nil
	u.Salt = nil
	json, err := json.Marshal(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(json)

}

func (h *userHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if !h.userService.Offline() {
		id_token := r.Header.Get("X-Id-Token")
		if id_token == "" {
			log.Printf("handler|GetUsers no x-id-token provided\n")
			LogUnauthorizedRequest(w, r)
			return
		}
		err := h.userService.Auth(id_token, UserAdmin)
		if err != nil {
			log.Println(err)
			LogUnauthorizedRequest(w, r)
			return
		}
	}
	uu, err := h.userService.GetUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("GET /api/users - Something went wrong: " + err.Error()))
		return
	}

	for i := 0; i < len(uu); i++ {
		uu[i].PasswordHash = nil
		uu[i].Salt = nil
	}

	json, err := json.MarshalIndent(uu, "", " ")
	if err != nil {

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("GET /api/users - Something went wrong:" + err.Error()))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id_token := r.Header.Get("X-Id-Token")
	if id_token == "" {
		log.Printf("UpdateUser no X-Id-token provided\n")
		LogUnauthorizedRequest(w, r)
		return
	}
	vars := mux.Vars(r)
	userID := vars["id"]
	if userID == "" {
		log.Println("No user id provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		msg := fmt.Sprintf("%s - /api/user - UpdateUser couldn't decode request: %s\n", r.Method, err)
		log.Print(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("POST /api/user - UpdateUser couldn't decode request:" + err.Error()))
		return
	}
	auth := h.userService.Auth(id_token, "user")
	if auth != nil {
		log.Printf("handler|UpdateUser failed: %s\n", auth.Error())
		LogUnauthorizedRequest(w, r)
		return
	}
	if user.ID != id_token {
		u, err := h.userService.GetUser(id_token)
		if err != nil {
			log.Printf("Couldn't get User from id_token\n")
			LogUnauthorizedRequest(w, r)
			return
		}
		if u.Role != "admin" {
			log.Printf("Someone is trying to fiddle with the backend\n")
			LogUnauthorizedRequest(w, r)
			return
		}
	}
	err = h.userService.UpdateUser(&user)
	if err != nil {
		log.Printf("%s - /api/user - Couldn't Update user: %s\n", r.Method, err)
	}
	json, err := json.MarshalIndent(user, "", " ")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id_token := r.Header.Get("X-Id-Token")
	if id_token == "" {
		log.Printf("DeleteUser no X-Id-token provided\n")
		LogUnauthorizedRequest(w, r)
		return
	}

	auth := h.userService.Auth(id_token, "admin")
	if auth != nil {
		log.Printf("handler|UpdateUser failed: %s\n", auth.Error())
		LogUnauthorizedRequest(w, r)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	u, _ := h.userService.GetUser(id)
	if u != nil {
		err := h.userService.DeleteUser(id)
		if err != nil {
			log.Printf("handler|DeleteUser failed: %q\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Write([]byte("Deleted user"))
		return
	}
	w.WriteHeader(http.StatusNotFound)

}

func (h *userHandler) Auth(w http.ResponseWriter, r *http.Request) {
	var auth authRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&auth)
	if err != nil {
		log.Printf("handlers|Auth couldn't decode: %q", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR"))
	}

	log.Println("handlers|Auth:", auth)
	w.Write([]byte("OKAY"))
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	var login loginRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&login)
	if err != nil {
		log.Printf("handlers|Login couldn't decode: %q", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR"))
	}

	user, err := h.userService.Login(login)
	if err != nil {
		log.Println(user)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized:" + err.Error()))
		return
	}
	json, err := json.Marshal(user)
	if err != nil {
		log.Printf("handlers|Login couldn't marshal: %q\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(json)
}

func (h *userHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var credentials loginRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&credentials)

	if err != nil {
		log.Println("UpdatePassword:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.userService.UpdatePassword(credentials.ID, credentials.Password)
	if err != nil {
		log.Printf("handler|UpdatePassword - failed: %q\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password updated!"))
}

func (h *userHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var credentials loginRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&credentials)

	if err != nil {
		log.Println("ResetPassword failed:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.userService.ResetPassword(credentials.Email, credentials.Password)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User got new password"))
}

func (h *userHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Id-Token")
	s, err := h.userService.GetSettings(userID)
	if err != nil {
		log.Printf("handler|GetSettings failed: %q\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	json, err := json.Marshal(s)
	if err != nil {
		log.Printf("handler|GetSettings failed: %q\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//log.Println("handlers|GetSettings sending: ", s)
	w.Write(json)
}

func (h *userHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Id-Token")

	var settings *Settings
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&settings)
	log.Println(settings)
	if err != nil {
		log.Printf("handlers|UpdateSettings got a faulty body: %q\n", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.userService.UpdateSettings(userID, settings)
	if err != nil {
		log.Printf("handlers|UpdateSettings failed updating: %q\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Ok"))
}

func (h *userHandler) ResetDefaults(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s, err := h.userService.ResetDefaults(id)
	if err != nil {
		log.Printf("handler|ResetDefaults failed: %q\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json, err := json.Marshal(s)
	if err != nil {
		log.Printf("handler|ResetDefaults failed marshalling: %q", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func (h *userHandler) Identify(w http.ResponseWriter, r *http.Request) {

}

func (h *userHandler) Register(w http.ResponseWriter, r *http.Request) {
}

func (h *userHandler) RegisterWeb(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		log.Printf("handlers|RegisterWeb got faulty body: %q\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, err := h.userService.RegisterWeb(req)
	if err != nil {
		if err != nil {
			log.Printf("handlers|RegisterWeb failed: %q\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	json, err := json.Marshal(u)
	if err != nil {
		log.Printf("handlers|RegisterWeb failed marshalling: %q\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(json)
}
