package document

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type DocHandler interface {
	CreateDoc(w http.ResponseWriter, req *http.Request)
	ImportDoc(w http.ResponseWriter, req *http.Request)
	GetDoc(w http.ResponseWriter, req *http.Request)
	GetDocs(w http.ResponseWriter, req *http.Request)
	UpdateDoc(w http.ResponseWriter, req *http.Request)
	DeleteDoc(w http.ResponseWriter, req *http.Request)
}

type docHandler struct {
	docService DocService
}

type query struct {
	IDs    []string `json:"ids"`
	UserID string   `json:"user_id"`
}

//NewDocHandler returns a DocHandler paired with provided DocService
func NewDocHandler(docService DocService) DocHandler {
	return &docHandler{
		docService,
	}
}

//Endpoints sets handlers to the provided router
func Endpoints(r *mux.Router, h DocHandler) {
	r.HandleFunc("/api/docs", h.GetDocs).Methods("GET")
	r.HandleFunc("/api/docs/{id}", h.GetDoc).Methods("GET")
	r.HandleFunc("/api/docs", h.CreateDoc).Methods("POST")
	r.HandleFunc("/api/docs", h.UpdateDoc).Methods("PUT")
	r.HandleFunc("/api/docs/{id}", h.DeleteDoc).Methods("DELETE")
}

func (h *docHandler) CreateDoc(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Id-Token")

	var d *Document
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&d)
	d.Creator = userID
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	created, err := h.docService.CreateDoc(d)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(created)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(response)
}

func (h *docHandler) ImportDoc(w http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(32 << 20)
	file, handler, err := req.FormFile("fileupload")
	if err != nil {
		fmt.Println("FormFile error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()
	//	fmt.Fprintf(w, "%v", handler.Header)

	f, err := os.OpenFile("./uploads/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		w.WriteHeader(http.StatusNotModified)
		w.Write([]byte(err.Error()))
		log.Fatal(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	doc, err := h.docService.ImportDoc(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("couldn't parse file:" + err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(doc))
}

func (h *docHandler) GetDoc(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	doc, err := h.docService.GetDoc(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(doc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (h *docHandler) GetDocs(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Id-Token")
	docs, err := h.docService.GetDocs(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	response, err := json.Marshal(docs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (h *docHandler) UpdateDoc(w http.ResponseWriter, req *http.Request) {
	var d *Document
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&d)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	updated, err := h.docService.UpdateDoc(d)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	response, err := json.Marshal(updated)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (h *docHandler) DeleteDoc(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := h.docService.DeleteDoc(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("deleted [" + id + "]"))
}
