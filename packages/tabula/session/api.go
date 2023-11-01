package session

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *sessionService) GetCaption(id string) string {
	pool, ok := s.pools[id]
	if ok {
		text, err := pool.Tabula.GETlivecap(nil)
		if err != nil {
			log.Println("Couldn't get cap:", err)
			return ""
		}
		return string(text)
	}
	return ""
}

func (h *sessionHandler) getCaption(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	captions := h.Service.GetCaption(id)
	if captions != "" {
		w.Write([]byte(captions))
		return
	}
	w.Write([]byte(`{ "error": "no caps"}`))
}
