package main

import (
	"log"
	"net/http"

	"github.com/minio/selfupdate"
)

func doUpdate(url string) error {
	log.Println("Updating from:", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()

	err = selfupdate.Apply(resp.Body, selfupdate.Options{})
	if err != nil {
		log.Println("Failed update")
	}
	log.Println("update successful")
	return err
}
