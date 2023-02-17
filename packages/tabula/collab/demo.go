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
		log.Panic("LoadAuto bug:", err)
	}
	defer file.Close()

}
