package collab

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"
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
	if text == "" {
		return
	}
	if len(text) < 3 {
		return
	}
	if len(text) > 35 {
		lines := strings.Split(text, "\n")
		log.Println("Number of lines:", len(lines))
		if len(lines) > 3 {
			text = strings.Join(lines[len(lines)-2:], "\n")
			log.Println("Last two lines?\n#############\n", text)
		}
	}
	last := string(text[len(text)-1])
	log.Println("last:", last)
	if last == " " || unicode.IsPunct([]rune(last)[0]) {
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
