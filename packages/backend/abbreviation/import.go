package abbreviation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type textOnTopJSON struct {
	Autocorrect struct {
		Default struct {
			List map[string]string
		} `json:"<default>"`
	} `json:"autocorrect"`
	Shortform map[string]struct {
		Shortforms map[string]string
	}
}

type illumiTypeJSON struct {
	Lists []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Type int    `json:"type"`
	} `json:"lists"`
	Abbreviations []struct {
		ListID       int    `json:"listId"`
		Abbreviation string `json:"abbreviation"`
		Word         string `json:"word"`
	} `json:"abbreviations"`
}

func (s *abbService) ImportTextOnTop(userID string, dat []byte) (map[string][]*Abbreviation, error) {
	var totlist textOnTopJSON
	var ac []*Abbreviation
	var imported = make(map[string][]*Abbreviation)

	json.Unmarshal(dat, &totlist)

	for abb, word := range totlist.Autocorrect.Default.List {
		ac = append(ac, &Abbreviation{
			ID:      uuid.New().String(),
			Word:    word,
			Abb:     abb,
			Creator: userID,
			Updated: time.Now(),
		})
	}
	for name, list := range totlist.Shortform {
		var abbs []*Abbreviation
		if name == "<default>" {
			name = "# Förkortningslista (Standard)"
		}
		for abb, word := range list.Shortforms {
			abb := Abbreviation{ID: uuid.New().String(), Word: word, Abb: abb, Creator: userID, Updated: time.Now()}
			abbs = append(abbs, &abb)
		}
		imported[name] = abbs
	}
	imported["# Autokorrigering (Standard)"] = ac
	return imported, nil
}

func (s *abbService) ImportIllumiType(userID string, dat []byte) (map[string][]*Abbreviation, error) {
	var illumiList illumiTypeJSON
	var listNames = make(map[int]string)
	var imported = make(map[string][]*Abbreviation)
	err := json.Unmarshal(dat, &illumiList)
	if err != nil {
		return nil, fmt.Errorf("ImportIllumiType|Couldn't unmarshal:\n%s", err.Error())
	}
	for _, list := range illumiList.Lists {
		listNames[list.ID] = list.Name
	}
	for _, abb := range illumiList.Abbreviations {
		a := Abbreviation{ID: uuid.New().String(), Word: abb.Word, Abb: abb.Abbreviation, Creator: userID, Updated: time.Now(), Remind: true}
		imported[listNames[abb.ListID]] = append(imported[listNames[abb.ListID]], &a)
	}
	return imported, nil
}

func (s *abbService) ImportProtype(userID string, dat []byte) ([]*Abbreviation, error) {
	rs := bytes.Runes(dat)
	var first bool
	first = true
	var abbs []*Abbreviation
	var rawAbb, rawWord []rune
	var length rune
	length = rs[2]
	//log.Printf("First length: %d\n", length)
	//log.Printf("File length: %d\n", len(rs[2:])+1)
	for i := range rs[2:] {
		if i+1 == len(rs[2:]) {
			abb := Abbreviation{ID: uuid.New().String(), Word: string(rawWord), Abb: string(rawAbb), Creator: userID, Updated: time.Now()}
			abbs = append(abbs, &abb)
			rawWord = nil
			rawAbb = nil
			break

		}

		if length == 0 {
			length = rune(dat[i+3])
			if length == 255 {
				length += rune(dat[i+4]) + 3
			}
			//	log.Printf("Length is 0, restarting with %d\n", length)
			if first {
				first = false
			} else if first == false {
				first = true
				abb := Abbreviation{ID: uuid.New().String(), Word: string(rawWord), Abb: string(rawAbb), Creator: userID, Updated: time.Now()}
				abbs = append(abbs, &abb)
				rawWord = nil
				rawAbb = nil
			}
			continue
		}
		if first == true {
			rawAbb = append(rawAbb, rune(dat[i+3]))
			length--
		} else {
			rawWord = append(rawWord, rune(dat[i+3]))
			length--
		}
	}

	return abbs, nil
}

func (s *abbService) CheckForConflicts(listID string, abbs []*Abbreviation) ([]conflict, error) {
	target, err := s.GetAbbs(listID)
	if err != nil {
		return nil, fmt.Errorf("service|CheckForConflicts couldn't get list: %s", err.Error())
	}
	conflicts := []conflict{}
	for _, a := range target {
		for _, newAbb := range abbs {
			if a.Abb == newAbb.Abb {
				conflicts = append(conflicts, conflict{Abb: a.Abb, Old: a.Word, New: newAbb.Word})

			}
		}
	}
	return conflicts, nil
}

func (s *abbService) Import(userID string, abbs []*Abbreviation) []*Abbreviation {
	for _, a := range abbs {
		a.Abb = strings.ToLower(a.Abb)
		a.Creator = userID
		a.Updated = time.Now()
		a.ID = uuid.New().String()
	}
	return abbs
}

func (s *abbService) ImportAbbsToList(listID string, abbs []*Abbreviation) error {
	err := s.repo.ImportAbbsToList(listID, abbs)
	if err != nil {
		return fmt.Errorf("service|ImportAbbsToList failed: %q", err)
	}
	return s.TouchList(listID)
}
