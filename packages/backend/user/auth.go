package user

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type authRequest struct {
	Register   string `json:"register"`
	Inform     string `json:"inform"`
	MachineID  string `json:"machine_id"`
	LicenseKey string `json:"license_key"`
	Email      string `json:"email"`
}

type loginRequest struct {
	Email       string    `json:"email"`
	ID          string    `json:"id"`
	Password    string    `json:"password"`
	OldPassword string    `json:"oldpassword"`
	Local       bool      `json:"local"`
	LastLogin   time.Time `json:"last_login"`
	LastSync    time.Time `json:"last_sync"`
}

func (h *userHandler) IsRegistered(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]
	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, err := h.userService.GetUserByEmail(email)
	if err != nil {
		log.Printf("handlers|IsRegistered failed getting user by email: %q\n", err)
		if h.mode == "desktop" {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if u == nil {
		log.Printf("handlers|IsRegistered couldn't find user with provided email: %s\n", email)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if h.userService.Offline() {
		if u.LastSync.IsZero() {
			w.Write([]byte("sync"))
			return
		}
		w.Write([]byte("sync"))
		return
	}
	if u.PasswordHash == nil {
		if u.ResetPassword == true {
			log.Println(u)
			w.Write([]byte("reset"))
			return
		}
		w.Write([]byte("register"))
	} else {
		w.Write([]byte("login"))
	}

}

func (s *userService) RegisterWeb(l loginRequest) (*User, error) {
	u, err := s.repo.GetUserByEmail(l.Email)
	if err != nil {
		return nil, fmt.Errorf("service|Register failed retrieving user by email: %q", err)
	}
	if u == nil {
		return &User{}, fmt.Errorf("service|Register couldn't get user with email %s", l.Email)
	}
	//TODO: Initialize base settings and subscribe to standard global lists
	s.UpdatePassword(u.ID, l.Password)
	return u, nil
}

func (s *userService) Auth(id string, role Role) error {
	if id == "" {
		return fmt.Errorf("No ID token provided")
	}

	u, err := s.GetUser(id)
	if err != nil {
		return fmt.Errorf("service|Auth failed: %s", err.Error())
	}
	if u.Role == "admin" {
		return nil
	}
	if u.Role == role {

		return nil
	}
	return fmt.Errorf("Not authorized")
}

func (s *userService) GetMachineID() (string, error) {
	return "", nil
}

func (s *userService) Login(login loginRequest) (*User, error) {
	log.Println("what email are we looking for?", login.Email)
	user, err := s.repo.Login(login.Email, "")
	if user == nil || err != nil {
		log.Println("repo.Login found no user")
		return nil, fmt.Errorf("User not found")
	}

	if login.Local && s.mode == "desktop" {
		log.Println("Found local user:", user)
		return user, nil
	}
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(login.Password))
	if err != nil {
		return nil, fmt.Errorf("Wrong password provided")
	}

	return user, nil
}
