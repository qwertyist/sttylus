package abbreviation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

//AbbHandler interface for all rest api methods
type AbbHandler interface {
	Abbreviate(w http.ResponseWriter, r *http.Request)
	Lookup(w http.ResponseWriter, r *http.Request)
	FilterAbbs(w http.ResponseWriter, r *http.Request)

	GetAbb(w http.ResponseWriter, r *http.Request)
	GetAbbs(w http.ResponseWriter, r *http.Request)
	CreateAbb(w http.ResponseWriter, r *http.Request)
	UpdateAbb(w http.ResponseWriter, r *http.Request)
	DeleteAbb(w http.ResponseWriter, r *http.Request)
	DeleteAbbByID(w http.ResponseWriter, r *http.Request)

	GetList(w http.ResponseWriter, r *http.Request)
	GetLists(w http.ResponseWriter, r *http.Request)
	GetUserLists(w http.ResponseWriter, r *http.Request)
	CreateList(w http.ResponseWriter, r *http.Request)
	UpdateList(w http.ResponseWriter, r *http.Request)
	DeleteList(w http.ResponseWriter, r *http.Request)

	CopyStandardList(w http.ResponseWriter, r *http.Request)
	Cache(w http.ResponseWriter, r *http.Request)
	InitSharedList(w http.ResponseWriter, r *http.Request)
	JoinSharedList(w http.ResponseWriter, r *http.Request)
	CreateSharedAbb(w http.ResponseWriter, r *http.Request)
	RemoveSharedAbb(w http.ResponseWriter, r *http.Request)
	GetSharedAbbs(w http.ResponseWriter, r *http.Request)

	DontRemindAbb(w http.ResponseWriter, r *http.Request)
	GetSuggestions(w http.ResponseWriter, r *http.Request)
	IgnoreSuggestion(w http.ResponseWriter, r *http.Request)
	IgnoreAllSuggestions(w http.ResponseWriter, r *http.Request)

	UploadProType(w http.ResponseWriter, r *http.Request)
	UploadTextOnTop(w http.ResponseWriter, r *http.Request)
	UploadIllumiType(w http.ResponseWriter, r *http.Request)
	UploadTxt(w http.ResponseWriter, r *http.Request)
	Import(w http.ResponseWriter, r *http.Request)
	ImportTo(w http.ResponseWriter, r *http.Request)
	CheckForConflicts(w http.ResponseWriter, r *http.Request)
	ImportProType(w http.ResponseWriter, r *http.Request)

	ExportLists(w http.ResponseWriter, r *http.Request)
	GetPublicList(w http.ResponseWriter, r *http.Request)
	CreatePublicList(w http.ResponseWriter, r *http.Request)
}

type abbHandler struct {
	abbService AbbService
}

//NewAbbHandler returns handlers connected to the given abbservice
func NewAbbHandler(abbService AbbService) AbbHandler {
	return &abbHandler{
		abbService,
	}
}

//Endpoints returns the gorilla mux router with handlers
func Endpoints(r *mux.Router, h AbbHandler) {
	r.HandleFunc("/abbs/abbreviate/", h.Abbreviate).Methods("GET")
	r.HandleFunc("/abbs/abbreviate/{abb}", h.Abbreviate).Methods("GET")
	r.HandleFunc("/abbs/lookup/{phrase}", h.Lookup).Methods("GET")

	r.HandleFunc("/abbs/abbreviations/{listID}", h.GetAbbs).Methods("POST")
	r.HandleFunc("/abbs/abbreviation/{listID}/{abb}", h.GetAbb).Methods("GET")
	r.HandleFunc("/abbs/abbreviation/{list}", h.CreateAbb).Methods("POST")
	r.HandleFunc("/abbs/abbreviation/{list}", h.UpdateAbb).Methods("PUT")
	r.HandleFunc("/abbs/abbreviation/{list}/{id}", h.DeleteAbb).Methods("DELETE")
	r.HandleFunc("/abbs/abbreviation/{list}/id/{id}", h.DeleteAbbByID).Methods("DELETE")

	r.HandleFunc("/abbs/lists", h.GetLists).Methods("POST")
	r.HandleFunc("/abbs/lists", h.GetUserLists).Methods("GET")
	r.HandleFunc("/abbs/filter", h.FilterAbbs).Methods("POST")
	r.HandleFunc("/abbs/standardlist/{id}", h.CopyStandardList).Methods("GET")
	r.HandleFunc("/abbs/list/{id}", h.GetList).Methods("GET")
	r.HandleFunc("/abbs/list", h.CreateList).Methods("POST")
	r.HandleFunc("/abbs/list", h.UpdateList).Methods("PUT")
	r.HandleFunc("/abbs/list/{id}", h.DeleteList).Methods("DELETE")

	r.HandleFunc("/abbs/cache", h.Cache).Methods("POST")
	r.HandleFunc("/abbs/shared", h.InitSharedList).Methods("GET")
	r.HandleFunc("/abbs/shared", h.InitSharedList).Methods("POST")
	r.HandleFunc("/abbs/shared/{id}", h.GetSharedAbbs).Methods("GET")
	r.HandleFunc("/abbs/shared/{id}", h.JoinSharedList).Methods("PUT")
	r.HandleFunc("/abbs/shared/{id}", h.CreateSharedAbb).Methods("POST")
	r.HandleFunc("/abbs/shared/{id}/{abb}", h.RemoveSharedAbb).Methods("DELETE")

	r.HandleFunc("/abbs/suggestions", h.GetSuggestions).Methods("GET")
	r.HandleFunc("/abbs/suggestions/{word}", h.IgnoreSuggestion).Methods("DELETE")
	r.HandleFunc("/abbs/suggestions", h.IgnoreAllSuggestions).Methods("DELETE")
	r.HandleFunc("/abbs/learned/{abb}", h.DontRemindAbb).Methods("GET")

	r.HandleFunc("/abbs/upload/protype", h.UploadProType).Methods("POST")
	r.HandleFunc("/abbs/upload/textontop", h.UploadTextOnTop).Methods("POST")
	r.HandleFunc("/abbs/upload/txt", h.UploadTxt).Methods("POST")
	r.HandleFunc("/abbs/upload/illumitype", h.UploadIllumiType).Methods("POST")

	r.HandleFunc("/abbs/import", h.Import).Methods("POST")
	r.HandleFunc("/abbs/import/{listID}", h.ImportTo).Methods("POST")
	r.HandleFunc("/abbs/conflicts/{listID}", h.CheckForConflicts).Methods("POST")
	r.HandleFunc("/abbs/export/{target}", h.ExportLists).Methods("POST")
	r.HandleFunc("/abbs/public/{id}", h.CreatePublicList).Methods("POST")
	r.HandleFunc("/abbs/public/{short_id}", h.GetPublicList).Methods("GET")

}

func (h *abbHandler) GetAbb(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	request := vars["abb"]
	listID := vars["listID"]
	if listID == "" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	abb, err := h.abbService.GetAbb(listID, request)
	if err != nil {
		log.Printf("handler couldn't Get Abb (%s: %s): %s\n", listID, request, err)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	response, err := json.Marshal(abb)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *abbHandler) GetAbbs(w http.ResponseWriter, r *http.Request) {
	/*	var query query
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&query)

		if err != nil {
			log.Printf("handler|GetAbbs got bad request: %q\n", err)
			return
		}
	*/
	//id := r.Header.Get("X-Id-Token")
	vars := mux.Vars(r)

	listID := vars["listID"]
	abbs, err := h.abbService.GetAbbs(listID)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(err)
		_, _ = w.Write(response)
		return
	}

	response, err := json.Marshal(abbs)
	if err != nil {
		log.Printf("handlers|GetAbbs failed marshalling: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (h *abbHandler) CreateAbb(w http.ResponseWriter, req *http.Request) {
	var abb Abbreviation
	vars := mux.Vars(req)
	listID := vars["list"]
	userID := req.Header.Get("X-Id-Token")
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&abb)
	if err != nil {
		log.Printf("Bad request: %q\n", err)
		w.WriteHeader(http.StatusBadRequest)
		json, _ := json.Marshal(abb)
		w.Write(json)
		return
	}

	if strings.ToLower(abb.Abb) == abb.Word {
		deleted, err := h.abbService.DeleteAbb(listID, abb.Abb)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("[create abb] Couldn't delete abb=abb: " + err.Error()))
			return
		}
		h.abbService.UpdateCache(userID)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("deleted" + deleted.Abb))
		return
	}
	if abb.Abb == "" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	err = h.abbService.CreateAbb(listID, &abb)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("[createAbb] couldn't create abb" + err.Error()))
		log.Printf("handler|createAbb couldn't create abb: %s\n", err.Error())
		return
	}
	h.abbService.UpdateCache(userID)

	response, _ := json.Marshal(abb)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(response)
}

func (h *abbHandler) UpdateAbb(w http.ResponseWriter, req *http.Request) {
	var abb Abbreviation
	vars := mux.Vars(req)
	listID := vars["list"]

	decoder := json.NewDecoder(req.Body)
	_ = decoder.Decode(&abb)
	err := h.abbService.UpdateAbb(listID, &abb)
	if err != nil {
		log.Println("couldnt update abb", err)
	}

	response, _ := json.Marshal(abb)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)

}

func (h *abbHandler) DeleteAbb(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	listID := vars["list"]
	abb := vars["id"]

	deleted, err := h.abbService.DeleteAbb(listID, abb)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("[delete abb] Couldn't delete abb: " + err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("deleted" + deleted.Abb))
}

func (h *abbHandler) DeleteAbbByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	listID := vars["list"]
	id := vars["id"]

	deleted, err := h.abbService.DeleteAbbByID(listID, id)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("[delete abb] Couldn't delete abb: " + err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("deleted" + deleted.Abb))
}

func (h *abbHandler) GetList(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	list, err := h.abbService.GetList(id)
	if err != nil {
		log.Printf("handler|GetList failed: %s\n", err.Error())
	}
	response, _ := json.Marshal(list)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (h *abbHandler) GetLists(w http.ResponseWriter, req *http.Request) {
	var query query
	buf, bodyErr := ioutil.ReadAll(req.Body)
	if bodyErr != nil {
		log.Print("bodyErr ", bodyErr.Error())
		http.Error(w, bodyErr.Error(), http.StatusInternalServerError)
		return
	}
	//	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
	//log.Printf("handler|GetLists: %q", rdr1)
	req.Body = rdr2
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&query)

	if err != nil {
		log.Println("handler GetLists failed decoding:", err)
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(err)
		_, _ = w.Write(response)
		return
	}
	lists, err := h.abbService.GetLists(query.ListIDs)
	response, err := json.Marshal(lists)

	if err != nil {
		log.Println("handlerGetLists failed marshalling", err)
		response, _ := json.Marshal(err)
		_, _ = w.Write(response)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (h *abbHandler) GetUserLists(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Id-Token")
	if userID == "" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	lists, err := h.abbService.GetUserLists(userID)
	response, err := json.Marshal(lists)

	if err != nil {
		log.Println("handlerGetLists failed marshalling", err)
		response, _ := json.Marshal(err)
		_, _ = w.Write(response)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (h *abbHandler) CopyStandardList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	l, err := h.abbService.CopyStandardList(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(l))

}

func (h *abbHandler) CreateList(w http.ResponseWriter, r *http.Request) {
	var list List
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&list)
	if err != nil {
		log.Printf("handler|CreateListh got faulty request: %q", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID := r.Header.Get("X-Id-Token")
	list.Creator = userID
	if userID == "" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id, err := h.abbService.CreateList(&list)

	if err != nil {
		log.Printf("handler|CreateList failed: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	l := List{ID: id}
	response, err := json.Marshal(l)
	if err != nil {
		log.Printf("handler|CreateList couldn't marshal: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (h *abbHandler) UpdateList(w http.ResponseWriter, req *http.Request) {
	var list List
	decoder := json.NewDecoder(req.Body)
	_ = decoder.Decode(&list)
	h.abbService.UpdateList(&list)

	response, _ := json.Marshal(list)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)

}

func (h *abbHandler) DeleteList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listID := vars["id"]

	userID := r.Header.Get("X-Id-Token")

	list, err := h.abbService.DeleteList(userID, listID)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Delete list failed: " + err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("List (" + list.Name + ") deleted"))
}

func (h *abbHandler) Cache(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Id-Token")
	var query query
	buf, bodyErr := ioutil.ReadAll(r.Body)
	if bodyErr != nil {
		log.Print("bodyErr ", bodyErr.Error())
		http.Error(w, bodyErr.Error(), http.StatusInternalServerError)
		return
	}

	//	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
	//	log.Printf("handler|Cache: %q", rdr1)
	r.Body = rdr2
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&query)

	if err != nil {
		log.Printf("handler|Cache: Decode: %q\n", err)
		return
	}
	query.UserID = userID
	err = h.abbService.Cache(query)
	if err != nil {
		log.Printf("handler|Cache %q:", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *abbHandler) Abbreviate(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Id-Token")
	vars := mux.Vars(r)
	abb := vars["abb"]
	if abb == "" {
		log.Println("Empty abb request")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	word, found := h.abbService.Abbreviate(userID, abb)

	response := response{
		word,
		word,
		found,
	}
	json, err := json.Marshal(response)
	if err != nil {
		log.Printf("handler|Abbreviate couldn't marshal response: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func (h *abbHandler) Lookup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	phrase := vars["phrase"]
	userID := r.Header.Get("X-Id-Token")

	abbs := h.abbService.Lookup(userID, phrase)
	if len(abbs) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	json, err := json.Marshal(abbs)
	if err != nil {
		log.Printf("handler|Lookup couldn't marshal: %q\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(json)

}

func (h *abbHandler) GetSuggestions(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Id-Token")
	suggs, _ := h.abbService.GetSuggestions(userID)
	json, err := json.Marshal(suggestions{Suggestions: suggs})
	if err != nil {
		log.Printf("handler|GetSuggestions couldn't marshal suggestions: %q\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func (h *abbHandler) IgnoreSuggestion(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Id-Token")
	vars := mux.Vars(r)
	word := vars["word"]

	ok := h.abbService.IgnoreSuggestion(userID, word)
	if !ok {
		w.Write([]byte("NO"))
		return
	}
	w.Write([]byte("OK"))
}

func (h *abbHandler) IgnoreAllSuggestions(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Id-Token")
	h.abbService.IgnoreAllSuggestions(userID)

}

func (h *abbHandler) DontRemindAbb(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Id-Token")
	if userID != "" {
		vars := mux.Vars(r)
		abb := vars["abb"]
		err := h.abbService.DontRemindAbb(userID, abb)

		if err != nil {
			log.Println(err)
			w.Write([]byte("No abbreviation to change"))
			return
		}
		w.Write([]byte(abb))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("No user ID provided"))

}

func (h *abbHandler) UploadProType(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Id-Token")
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	abbs, err := h.abbService.ImportProtype(userID, fileBytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json, err := json.Marshal(abbs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(json)
}

func (h *abbHandler) UploadTextOnTop(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Id-Token")

	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	abbs, err := h.abbService.ImportTextOnTop(userID, fileBytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json, err := json.Marshal(abbs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(json)
}

func (h *abbHandler) UploadIllumiType(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Id-Token")

	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	abbs, err := h.abbService.ImportIllumiType(userID, fileBytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json, err := json.Marshal(abbs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(json)
}

func (h *abbHandler) UploadTxt(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Id-Token")
	var abbs []*Abbreviation

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&abbs)

	if err != nil {
		log.Printf("UploadTxt failed: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	abbs = h.abbService.Import(userID, abbs)
	response, err := json.Marshal(abbs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (h *abbHandler) Import(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-Id-Token")
	var abbs []*Abbreviation

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&abbs)
	if err != nil {
		log.Printf("handler|Import failed: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	abbs = h.abbService.Import(userID, abbs)
	response, err := json.Marshal(abbs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (h *abbHandler) ImportTo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listID := vars["listID"]
	//fmt.Printf("Import abbreviations to listID %s\n", listID)
	list, err := h.abbService.GetList(listID)
	if list == nil {
		log.Printf("handler|ImportTo couldn't find list %s\n", listID)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("handler|ImportTo failed: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var abbs []*Abbreviation

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&abbs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.abbService.ImportAbbsToList(listID, abbs)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte("OK"))
}

func (h *abbHandler) CheckForConflicts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listID := vars["listID"]

	var abbs []*Abbreviation

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&abbs)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conflicts, err := h.abbService.CheckForConflicts(listID, abbs)
	if err != nil {
		log.Printf("handler|CheckForConflicts failed: %q", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json, err := json.Marshal(conflicts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func (h *abbHandler) ImportProType(w http.ResponseWriter, r *http.Request) {

}

func (h *abbHandler) GetPublicList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	short_id := vars["short_id"]
	log.Println("short_id:", short_id)
	publicList, err := h.abbService.GetPublicList(short_id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("no public list with ID corresponding to short_id"))
		return
	}

	if publicList.Name == "" {
		publicList, err = h.abbService.GetList(publicList.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("get public list metadata failed"))
			return
		}
	}

	abbs, err := h.abbService.GetAbbs(publicList.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("get public list abbs failed"))
		return
	}

	if len(abbs) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	publicListObject := map[string]interface{}{
		"metadata": publicList,
		"abbs":     abbs,
	}

	resp, err := json.Marshal(publicListObject)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("couldn't marshal public list "))
		return
	}

	w.Write(resp)
}

func (h *abbHandler) CreatePublicList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var list List
	log.Println(id)
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&list)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("couldn't decode body", err)
		return
	}

	h.abbService.CreatePublicList(id, &list)
}
