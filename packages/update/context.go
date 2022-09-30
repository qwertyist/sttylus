package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func ManifestCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		manifest := &Manifest{}
		var err error
		version := chi.URLParam(r, "version")
		log.Println("url param version:", version)
		if version := chi.URLParam(r, "version"); version != "" {
			log.Println("Look for version...")
			semver, err := StringToSemVer(version)
			if err != nil {
				render.Render(w, r, ErrInvalidRequest(err))
				return
			}

			manifest, err = dbGetManifest(semver)
		}
		/*
			else if major := chi.URLParam(r, "major"); major != "" {
				log.Println("urlParam 'major':", major)
				var minor, patch string
				if minor := chi.URLParam(r, "minor"); minor == "" {
					minor = "0"
				}
				log.Println("urlParam 'minor':", minor)
				if patch := chi.URLParam(r, "patch"); patch == "" {
					patch = "0"
				}
				log.Println("urlParam 'patch':", patch)
				semver, err := StringToSemVer(major + "." + minor + "." + patch)
				if err != nil {
					render.Render(w, r, ErrInvalidRequest(err))
					return
				}
				manifest, err = dbGetManifest(semver)
			}
		*/

		log.Println("manifest from db:", manifest)
		log.Println("db error:", err)
		if manifest == nil || err != nil {
			log.Println("not found @", chi.RouteContext(r.Context()))
			render.Render(w, r, ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), "manifest", manifest)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
