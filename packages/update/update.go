package main

import "net/http"

func handleGetLatestUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=latest.zip")
	http.ServeFile(w, r, "./updates/sttylus_update_0.8.1.zip")
}

func handleGetUpdate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET UPDATE"))
}
