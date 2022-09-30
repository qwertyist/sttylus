package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/botvid/webapp/abbreviation"
	"github.com/botvid/webapp/document"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func authedRequest(method, endpoint, id string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, "https://sttylus.se/api/"+endpoint, body)
	if err != nil {
		log.Fatal("couldn't create authed request:", err)
		return nil
	}
	req.Header.Add("X-Id-Token", id)
	return req
}

func (h *userHandler) SyncUser(w http.ResponseWriter, r *http.Request) {
	var login loginRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&login)
	if err != nil {
		log.Printf("handlers|SyncUser couldn't decode: %q", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR:" + err.Error()))
		return
	}
	pre, err := http.Get("https://sttylus.se/api/")
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	log.Println(pre.Status)
	user, err := h.userService.SyncUser(login)
	if err != nil {
		if err.Error() == "wrong password provided" {
			log.Println("Wrong password provided actually")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		log.Println("SyncUserHandler|", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println(login)
	w.Write(resp)
}

func (h *userHandler) SyncUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Println("Sync user by ID:", id)
}

func (s *userService) Offline() bool {
	if s.mode == "desktop" {
		log.Println("OFFLINE")
		return true
	}
	return false
}

func (s *userService) LocalLogin(l loginRequest) (*User, error) {
	user, err := s.repo.Login(l.Email, "")
	if err != nil {
		log.Fatal(err)
	}
	if user != nil {
		log.Println("User already stored locally:", user)
		return user, nil
	}
	credentials, err := json.Marshal(l)
	if err != nil {
		log.Fatal(err)
	}
	data := bytes.NewBuffer(credentials)
	resp, err := http.Post("https://sttylus.se/api/login", "application/json", data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	err = s.repo.CreateUser(user)
	if err != nil {
		log.Fatal(err)
	}
	return user, nil
}

func (s *userService) SyncUser(login loginRequest) (*User, error) {
	log.Println("Look for user locally")
	var user *User
	var err error
	if login.Email != "" {
		user, err = s.repo.GetUserByEmail(login.Email)
	} else if login.ID != "" {
		user, err = s.repo.GetUser(login.ID)
	}
	if err != nil {
		return nil, fmt.Errorf("SyncUser|Failed reading botldb: %s", err.Error())
	}
	if user == nil {
		log.Println("No user found, look for user online")
		credentials, err := json.Marshal(login)
		if err != nil {
			log.Fatal("Couldn't marshall login request", err)
		}
		data := bytes.NewBuffer(credentials)
		resp, err := http.Post("https://sttylus.se/api/login", "application/json", data)
		if err != nil {
			log.Fatal("login post request failed:", err)
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&user)
		if err != nil {
			return nil, fmt.Errorf("wrong password provided")
		}
		err = s.repo.CreateUser(user)
		if err != nil {
			log.Fatal(err)
		}
		err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(login.Password))
		if err != nil {
			return nil, fmt.Errorf("Wrong password provided")
		}
	}

	lists, err := s.SyncUserLists(user)
	if err != nil {
		return user, fmt.Errorf("Couldn't sync user lists: %s", err)
	}
	log.Printf("synced %d lists\n", len(lists))
	abbs, err := s.SyncUserAbbs(user.ID, lists)
	if err != nil {
		return user, fmt.Errorf("Couldn't sync user abbs: %s", err)
	}
	log.Printf("synced %d abbs\n", len(abbs))
	docs, err := s.SyncUserDocs(user)
	if err != nil {
		return user, fmt.Errorf("Couldn't sync user docs: %s", err)
	}
	log.Printf("synced %d docs\n", len(docs))
	log.Println("last sync:", login.LastSync)
	if login.LastSync.Before(time.Unix(1, 0)) {
		log.Println("First sync")
		user.LastSync = time.Now()
	} else {
		log.Println("Later sync")
		user.LastSync = login.LastSync
	}
	err = s.UpdateUser(user)
	if err != nil {
		return nil, fmt.Errorf("Couldn't update user: %s", err)
	}
	return user, nil
}

func (s *userService) SyncUserLists(user *User) ([]*abbreviation.List, error) {
	var lists []*abbreviation.List
	var newerLists []*abbreviation.List
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req := authedRequest("GET", "abbs/lists", user.ID, nil)
	if req == nil {
		return nil, fmt.Errorf("couldn't create GetUserLists post request")
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("getUserLists post request failed %s:", err)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&lists)
	if err != nil {
		return nil, fmt.Errorf("couldn't decode lists response: %s", err)
	}

	for _, list := range lists {
		if list.Updated.After(user.LastSync) {
			newerLists = append(newerLists, list)
		}
	}

	err = s.abbService.CreateManyLists(newerLists)
	if err != nil {
		return nil, fmt.Errorf("SyncUserLists|couldn't create many lists: %s", err.Error())
	}
	return newerLists, nil
}

func (s *userService) SyncUserAbbs(userID string, lists []*abbreviation.List) ([]*abbreviation.Abbreviation, error) {
	var newerAbbs []*abbreviation.Abbreviation
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	//r.HandleFunc("/api/abbs/abbreviations/{listID}", h.GetAbbs).Methods("POST")
	for _, list := range lists {
		req := authedRequest("POST", "abbs/abbreviations/"+list.ID, userID, nil)
		if req == nil {
			return nil, fmt.Errorf("Couldn't create GetAbbs post request")
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("GetAbbs post request failed %s:", err)
		}
		defer resp.Body.Close()
		var abbs []*abbreviation.Abbreviation
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&abbs)
		if err != nil {
			return nil, fmt.Errorf("couldn't decode abbs response: %s", err)
		}
		err = s.abbService.CreateManyAbbs(list.ID, abbs)
		if err != nil {
			return nil, fmt.Errorf("couldn't batch create abbs: %s", err)
		}
		newerAbbs = append(newerAbbs, abbs...)
	}
	return newerAbbs, nil
}

func (s *userService) SyncUserDocs(user *User) ([]*document.Document, error) {
	var docs []*document.Document
	var newerDocs []*document.Document
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req := authedRequest("GET", "docs", user.ID, nil)
	if req == nil {
		return nil, fmt.Errorf("Couldn't create GetUserLists post request")
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetUserLists post request failed %s:", err)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&docs)
	if err != nil {
		return nil, fmt.Errorf("couldn't decode lists response: %s", err)
	}
	preRelease := time.Date(2022, 06, 22, 0, 0, 0, 0, time.UTC)
	for _, doc := range docs {
		if doc.Updated.After(preRelease) {
			if doc.Updated.After(user.LastSync) {
				newerDocs = append(newerDocs, doc)
			}
		}
	}

	err = s.docService.CreateManyDocs(newerDocs)
	if err != nil {
		return nil, fmt.Errorf("SyncUserDocs|couldn't create many docs: %s", err.Error())
	}
	return newerDocs, nil
}
