package main

import (
	"log"
	"time"
)

var manifests map[string]*Manifest
var latest *Manifest

var users = []*User{
	{ID: 000, Email: "info@sttylus.se"},
	{ID: 001, Email: "tester@sttylus.se"},
}

type User struct {
	ID    int64
	Email string
}

type Manifest struct {
	Version string      `json:"version,omitempty"`
	Date    time.Time   `json:"date,omitempty"`
	Safe    bool        `json:"safe"`
	Major   int         `json:"major"`
	Minor   int         `json:"minor"`
	Patch   int         `json:"patch"`
	Comment string      `json:"comment"`
	Log     []*LogEntry `json:"change_log,omitempty"`
	Hash    string      `json:"hash"`
	ID      int         `json:"id"`
}

type Release struct {
	Description string
	Version     map[string]Manifest
}

func createFixtures() {
	preTime, err := time.Parse("2006-01-02", "2022-11-01")
	if err != nil {
		log.Fatal(err)
	}
	release_day, err := time.Parse("2006-01-02", "2022-05-09")
	if err != nil {
		log.Fatal(err)
	}
	manifests = make(map[string]*Manifest)
	m := &Manifest{Version: "0.8.0", Major: 0, Minor: 8, Patch: 0, Date: preTime, Comment: "Tabula rasa", Safe: true}
	initChangeLog(m)
	manifests["0.8.0"] = m
	m = &Manifest{Version: "0.7.4", Major: 0, Minor: 7, Patch: 4, Date: release_day, Comment: "Let's get this bread!", Safe: true}
	initChangeLog(m)
	manifests["0.7.4"] = m

}
