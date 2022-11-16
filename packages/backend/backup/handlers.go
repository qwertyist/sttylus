package backup

import (
	"net/http"

	"github.com/gorilla/mux"
)

type BackupHandler interface {
	DumpUsers(w http.ResponseWriter, r *http.Request)
	DumpAbbs(w http.ResponseWriter, r *http.Request)
	DumpManuscripts(w http.ResponseWriter, r *http.Request)
	DumpAll(w http.ResponseWriter, r *http.Request)
	RestoreUsers(w http.ResponseWriter, r *http.Request)
	RestoreAbbs(w http.ResponseWriter, r *http.Request)
	RestoreManuscripts(w http.ResponseWriter, r *http.Request)
	RestoreAll(w http.ResponseWriter, r *http.Request)
}

type backupHandler struct {
	backupService BackupService
}

func NewBackupHandler(backupService BackupService) BackupHandler {
	return &backupHandler{
		backupService,
	}
}

func Endpoints(r *mux.Router, h BackupHandler) {
	r.HandleFunc("/export", h.DumpAll).Methods("GET")
	r.HandleFunc("/export/users", h.DumpUsers).Methods("GET")
	r.HandleFunc("/export/abbs", h.DumpAbbs).Methods("GET")
}

func (h *backupHandler) DumpUsers(w http.ResponseWriter, r *http.Request) {
	err := h.backupService.DumpUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := "failed dumping Users:" + err.Error()
		w.Write([]byte(errMsg))
		return
	}
	w.Write([]byte("Exported users"))
}

func (h *backupHandler) DumpAbbs(w http.ResponseWriter, r *http.Request) {
	err := h.backupService.DumpAbbs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Exported abbs"))
}

func (h *backupHandler) DumpManuscripts(w http.ResponseWriter, r *http.Request) {}

func (h *backupHandler) DumpAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Export all not implemented"))
}

func (h *backupHandler) RestoreUsers(w http.ResponseWriter, r *http.Request) {}

func (h *backupHandler) RestoreAbbs(w http.ResponseWriter, r *http.Request) {}

func (h *backupHandler) RestoreManuscripts(w http.ResponseWriter, r *http.Request) {}

func (h *backupHandler) RestoreAll(w http.ResponseWriter, r *http.Request) {}
