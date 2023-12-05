package collab

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ZoomCC struct {
	Token        string
	Breakout     string
	MainStep     int
	BreakoutStep int
	LastKeycode  byte
	Block        string
	BlockTimer   *time.Timer
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
		return -1
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1
	}
	//Convert the body to type string
	step, err := strconv.Atoi(string(body))
	if err != nil {
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

	if t.Zoom.BlockTimer != nil {
		return fmt.Errorf("Zoom data already set")
	}
	t.Zoom.BlockTimer = time.NewTimer(5 * time.Second)
	go t.ZoomPOST()

	return nil
}
func (t *Tabula) GetLastAppend() string {
	text := t.ToText()
	last_index := strings.LastIndex(strings.TrimRight(text, " \n"), " ")
	return text[last_index+1:]
}

func (t *Tabula) ZoomPOST() {
	<-t.Zoom.BlockTimer.C
	if t.Zoom.Block == "" {
		log.Println("Block empty, re run")
		t.Zoom.BlockTimer = time.NewTimer(5 * time.Second)
		go t.ZoomPOST()
		return
	}
	t.Zoom.MainStep++
	log.Printf("step: %d\ttarget: %s\n", t.Zoom.MainStep, t.Zoom.Block)
	step := strconv.Itoa(t.Zoom.MainStep)
	target := t.Zoom.Token + "&lang=sv-SE" + "&seq=" + step
	msg := bytes.NewBuffer([]byte(t.Zoom.Block))
	resp, err := http.Post(target, "text/plain", msg)
	log.Println("POST sent")
	if err != nil {
		log.Println("ZoomPOST err", err)
	}
	t.Zoom.Block = ""
	if resp.StatusCode != 200 {
		err = fmt.Errorf("POST failed with status '%s'", resp.Status)
		fmt.Printf("error: %s\n", err.Error())
		if resp.StatusCode == 403 {
			t.Zoom.MainStep++
		}
	}

	t.Zoom.BlockTimer = time.NewTimer(5 * time.Second)
	go t.ZoomPOST()
}

func (t *Tabula) SendZoomCC(text string) {
	re := regexp.MustCompile(`[\n]{2,}`)
	text = re.ReplaceAllString(text, " ")
	t.Zoom.Block += text
}
