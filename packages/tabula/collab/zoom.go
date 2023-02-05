package collab

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ZoomCC struct {
	Token        string
	Breakout     string
	MainStep     int
	BreakoutStep int
}

func getZoomStep(token string) int {
	if token == "" {
		return -1
	}
	d := strings.Index(token, "closedcaption") + len("closedcaption")
	seq := token[:d] + "/seq/" + token[d:]
	log.Println(seq)
	resp, err := http.Get(seq)
	if err != nil {
		log.Fatalln(err)
		return -1
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return -1
	}
	//Convert the body to type string
	step, err := strconv.Atoi(string(body))
	if err != nil {
		log.Fatalln(err)
		return -1
	}

	return step
}

func (t *Tabula) SetZoomData(data ZoomCC) error {
	log.Println("tabula received zoom data:", data)
	t.Zoom = data
	step := getZoomStep(data.Token)
	if step < 0 {
		return fmt.Errorf("couldn't retrieve cc step for zoom meeting")
	}
	log.Println("connect to zoom cc at step:", step)
	t.Zoom.MainStep = step
	t.Zoom.Token = data.Token

	return nil
}

func (t *Tabula) SendZoomCC() {
	step := strconv.Itoa(t.Zoom.MainStep)
	target := t.Zoom.Token + "&lang=sv-SE" + "&seq=" + step
	text := t.ToText()
	log.Println(text)
	if text == "" {
		log.Println("Nothing to send")
		return
	}
	if len(text) < 3 {
		log.Println("Not sending less than 3 characters")
		return
	}
	last := text[len(text)-1]
	log.Println(last)
	// . 46
	// ! 33
	// ? 63
	// , 44
	// \n 10
	if last == 32 || last == 46 || last == 33 || last == 63 || last == 44 {
		log.Println("should send something after punctuation or whitespace")
		msg := bytes.NewBuffer([]byte(text))
		resp, err := http.Post(target, "text/plain", msg)
		if err != nil {
			log.Fatal(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		t.Zoom.MainStep++
		log.Println(string(body))
	}

}
