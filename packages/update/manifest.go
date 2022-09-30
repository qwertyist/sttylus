package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/render"
)

func getManifestHandler(w http.ResponseWriter, r *http.Request) {
	manifest, ok := r.Context().Value("manifest").(*Manifest)
	if !ok {
		log.Println("something went wrong!")
	}
	if manifest == nil {
		render.Render(w, r, ErrNotFound)
		return
	}
	err := render.Render(w, r, NewManifestResponse(manifest))
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func getManifestsHandler(w http.ResponseWriter, r *http.Request) {
	mm := make(map[int]*Manifest)
	for _, m := range manifests {
		if mm[m.Major*100+m.Minor] == nil {
			log.Println("Overwrite")
			mm[m.Major*100+m.Minor] = m

		}

		if mm[m.Major*100+m.Minor].Log == nil {
			initChangeLog(mm[m.Major*100+m.Minor])
		} else {
			mm[m.Major*100+m.Minor].Log = append(mm[m.Major*100+m.Minor].Log, m.Log...)
		}
		/*
				mm[m.Major].Log = append(mm[m.Major].Log, m.Log...)
			} else {
				mm[m.Major].Log = []*LogEntry{}
			}
		*/

	}
	release := Release{}
	release.Version = make(map[string]Manifest)
	for _, m := range mm {
		version, err := SemVerToString(&SemVer{Major: m.Major, Minor: m.Minor, Patch: 0})
		if err != nil {
			log.Fatal(err)
		}
		m.Version = version
	}

}

func findLatestUpdate() {
	var major, minor, patch int
	for _, m := range manifests {
		if m.Major > major {
			major = m.Major
			minor = m.Minor
			patch = m.Patch
			latest = m
		} else if m.Major < major {
			continue
		}
		if m.Minor > minor {
			minor = m.Minor
			patch = m.Patch
			latest = m
		} else if m.Minor < minor {
			continue
		}
		if m.Patch > patch {
			patch = m.Patch
			latest = m
		}
	}
}

func getLatestManifestHandler(w http.ResponseWriter, r *http.Request) {
	err := render.Render(w, r, NewManifestResponse(latest))
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

type ManifestResponse struct {
	*Manifest
}

func NewManifestResponse(manifest *Manifest) *ManifestResponse {

	return &ManifestResponse{Manifest: manifest}

}

func (rd *ManifestResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ManifestRequest struct {
	Manifest *Manifest    `json:"manifest,omitempty"`
	User     *UserPayload `json:"user,omitempty"`
}

func (m *ManifestRequest) Bind(r *http.Request) error {
	log.Printf("%+v\n", m)
	if m.Manifest == nil {
		return errors.New("Missing required manifest fields")
	}
	return nil
}

func createManifest(w http.ResponseWriter, r *http.Request) {
	data := &ManifestRequest{}
	if err := render.Bind(r, data); err != nil {
		log.Println("Couldn't bind!", err)
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	manifest := data.Manifest
	saved, found := manifests[manifest.Version]
	if found {
		log.Printf("update [\"%s\"] exists", saved.Version)
		render.Status(r, http.StatusNotModified)
		return
	}
	_, err := dbNewManifest(manifest)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
	}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewManifestResponse(manifest))
}

func deleteManifest(w http.ResponseWriter, r *http.Request) {

}
