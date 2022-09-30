package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/docgen"
	"github.com/go-chi/render"
)

var routes = flag.Bool("routes", true, "Generate route documentation")

func registerMiddlewares(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	//Timeout request context after 60 seconds
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(render.SetContentType(render.ContentTypeJSON))
}

func registerHandlers(r *chi.Mux) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("sttylus update server 0.0.0"))
	})
	r.Route("/versions", func(r chi.Router) {
		r.Get("/", getManifestsHandler)
		r.Get("/latest", getLatestManifestHandler)
		r.Route("/{version}", func(r chi.Router) {
			r.Use(ManifestCtx)
			r.Get("/", getManifestHandler)
		})
	})
	r.Route("/get", func(r chi.Router) {
		r.Get("/latest", handleGetLatestUpdate)
		r.Get("/{version}", handleGetUpdate)
	})
	r.Mount("/admin", adminRouter())
}

func init() {
	createFixtures()
	findLatestUpdate()
	render.Respond = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		if err, ok := v.(error); ok {

			if _, ok := r.Context().Value(render.StatusCtxKey).(int); !ok {
				w.WriteHeader(400)
			}

			fmt.Printf("Respond err: %s\n", err.Error())
			render.DefaultResponder(w, r, render.M{"status": "errror"})
			return
		}
		render.DefaultResponder(w, r, v)
	}
}

func main() {
	port := ":3000"
	r := chi.NewRouter()
	registerMiddlewares(r)
	registerHandlers(r)
	if *routes {

		docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
			Intro: "STTYLUS Update server endpoints",
		})
	}

	log.Println("STTylus update server listens at " + port)
	http.ListenAndServe(port, r)
}
