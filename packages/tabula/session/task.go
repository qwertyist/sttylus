package session

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *sessionService) SetAuthToken(token string) {
	s.AuthToken = token
	log.Println("Set TaskAuthToken to:", token)
}

func (s *sessionService) CheckAuthToken(token string) bool {
	return s.AuthToken == token
}

func (h *sessionHandler) taskResetSessions(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	token := params["token"]
	if h.Service.CheckAuthToken(token) {
		errs := h.Service.ResetSessions()
		if errs > 0 {
			log.Println("[API] CheckAuthToken failed:", token)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte("OK"))
	} else {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("ACCESS DENIED"))
	}
}

func (s *sessionService) ResetSessions() int {
	var errs []error
	log.Println("[API] Reset textcontents for", len(s.pools), "sessions")
	for _, sess := range s.pools {
		if len(sess.Clients) > 0 {
			continue
		}
		err := sess.Tabula.ClearHandler()
		if err != nil {
			errs = append(errs, err)
			log.Println("[API] ResetSessions:", err)
		}
	}
	return len(errs)
}
