package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func adminRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(AdminOnly)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("admin: index"))
	})
	r.Get("/accounts", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("admin: list accounts..."))
	})
	r.Get("/users/{userId}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("admin: view user id %v", chi.URLParam(r, "userId"))))
	})

	r.Route("/versions", func(r chi.Router) {
		r.Use(ManifestCtx)
		r.Post("/", createManifest)
		r.Delete("/{major}-{minor}-{patch}", deleteManifest)
	})
	return r
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		return
		/*isAdmin, ok := r.Context().Value("user").(bool)
		log.Println("isAdmin:", isAdmin)
		if !ok || !isAdmin {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		*/
	})

}
