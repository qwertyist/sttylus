package abbreviation

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type baseList struct {
	ID string `json:"id"`
}

var src = rand.New(rand.NewSource(time.Now().UnixNano()))

func (h *abbHandler) InitSharedList(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Id-Token")
	if userID == "" {
		log.Println("No X-Id-Token provided")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method == "GET" {
		listid := h.abbService.InitSharedList(userID, "")
		w.Write([]byte(listid))
		return
	} else if r.Method == "POST" {
		var baselist baseList
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&baselist)
		if err != nil {
			log.Println(err)
		}
		listid := h.abbService.InitSharedList(userID, baselist.ID)

		w.Write([]byte(listid))
	}
}

func (h *abbHandler) JoinSharedList(w http.ResponseWriter, r *http.Request) {
	userid := r.Header.Get("X-Id-Token")
	if userid == "" {
		log.Println("No X-Id-Token provided")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	vars := mux.Vars(r)
	listid := vars["id"]
	err := h.abbService.JoinSharedList(userid, listid)
	if err != nil {
		log.Println("JoinSharedList failed", err)
		listid = h.abbService.InitSharedList(userid, "")
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(listid))
}

func (h *abbHandler) CreateSharedAbb(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listid := vars["id"]

	var abb Abbreviation

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&abb)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.abbService.CreateSharedAbb(listid, abb)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *abbHandler) RemoveSharedAbb(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listid := vars["id"]
	abb := vars["abb"]

	err := h.abbService.RemoveSharedAbb(listid, abb)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *abbHandler) GetSharedAbbs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listid := vars["id"]

	abbs, err := h.abbService.GetSharedAbbs(listid)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(abbs)
	if err != nil {
		log.Println("GetSharedAbbs handler failed marshalling:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

func (s *abbService) InitSharedList(userID string, baseListID string) string {

	b := make([]byte, (4))
	if _, err := src.Read(b); err != nil {
		panic(err)
	}
	listid := hex.EncodeToString(b)[:4]

	s.cache.UserSharedList[userID] = listid
	s.cache.SharedAbbs[listid] = make(map[string]string)

	if baseListID != "" {
		abbs, err := s.GetAbbs(baseListID)
		if err != nil {
			return listid
		}

		for _, abb := range abbs {
			s.cache.SharedAbbs[listid][abb.Abb] = abb.Word
		}
	}
	return listid
}

func (s *abbService) JoinSharedList(userid, listid string) error {
	//log.Printf("Try join shared list %s ...", listid)
	if listid == "leave" {
		delete(s.cache.UserSharedList, userid)
		return nil
	}
	if s.cache.SharedAbbs[listid] != nil {
		//	log.Println("OK")
		s.cache.UserSharedList[userid] = listid
		return nil
	}
	//log.Println("Failed")
	return fmt.Errorf("Join: Shared list doesn't exist")
}

func (s *abbService) CreateSharedAbb(listid string, abb Abbreviation) error {
	if s.cache.SharedAbbs[listid] != nil {
		s.cache.SharedAbbs[listid][abb.Abb] = abb.Word
		return nil
	}

	return fmt.Errorf("Create abb: Shared list doesn't exist")
}

func (s *abbService) RemoveSharedAbb(listid, abb string) error {
	if s.cache.SharedAbbs[listid] != nil {
		delete(s.cache.SharedAbbs[listid], abb)
		return nil
	}
	return fmt.Errorf("Remove abb: Shared list doesn't exist")
}

func (s *abbService) GetSharedAbbs(listid string) ([]Abbreviation, error) {
	var aa []Abbreviation
	if s.cache.SharedAbbs[listid] != nil {
		for abb, word := range s.cache.SharedAbbs[listid] {
			a := Abbreviation{Abb: abb, Word: word}
			aa = append(aa, a)
		}
		return aa, nil
	}
	return nil, fmt.Errorf("Get abbs: Shared list doesn't exist")
}
