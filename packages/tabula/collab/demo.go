package collab

import (
	"log"
	"os"
	"time"
)

type DemoText struct {
	time time.Duration
}

func LoadAuto() {
	file, err := os.Open("auto.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

}
