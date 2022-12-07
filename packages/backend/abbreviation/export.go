package abbreviation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"
)

type protypeExport map[string]*bytes.Buffer

type textOnTopExport struct {
	Shortform   map[string]map[string]abbMap `json:"shortform"`
	Autocorrect struct {
		Default struct {
			List abbMap `json:"list,omitempty"`
		} `json:"<default>"`
	} `json:"autocorrect,omitempty"`
}

type abbMap map[string]string

func (h *abbHandler) ExportLists(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	target := vars["target"]
	var q query
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&q)
	if err != nil {
		log.Println("exportLists handler|Couldn't decode request:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var resp []byte
	if target == "textontop" {
		tot, err := h.abbService.CreateToTExport(q.Standard, q.Addon)
		enc := json.NewEncoder(w)
		enc.Encode(tot)
		if err != nil {
			log.Println("exportLists handler|Couldn't marshal to text-on-top:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(resp)
	} else if target == "protype" {
		first := strings.Fields(q.User)[0]
		err := h.abbService.CreateProtypeExport(q.User, q.Standard, q.Addon)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Disposition", "attachment")
		w.Header().Set("Content-Type", "application/zip; filename="+first+".zip")
		w.Header().Set("Content-Transfer-Encoding", "binary")
		w.Header().Set("Expires", "0")
		http.ServeFile(w, r, "tmp/"+q.User+".zip")
		defer os.Remove("tmp/" + q.User + ".zip")
		return
	}

}

func (s *abbService) CreateProtypeExport(User, sID string, aIDs []string) error {
	args := []string{"-zv", "--password", "SkrivTolk", "../" + User + ".zip"}
	standard, err := s.GetAbbs(sID)
	if err != nil {
		return fmt.Errorf("CreateProtypeExport|Couldn't get abbs from standard list:\n%s", err)
	}
	bytes, err := json.Marshal(standard)
	if err != nil {
		return fmt.Errorf("CreateProtypeExport|Couldn't marshall standard list:\n%s", err)
	}
	tmp, err := ioutil.TempDir("tmp", "protype")
	if err != nil {
		return fmt.Errorf("CreateProtypeExport|Couldn't create tempdir:\n%s", err)
	}
	defer os.RemoveAll(tmp)
	err = ioutil.WriteFile(tmp+"/list.json", bytes, 0644)
	if err != nil {
		return fmt.Errorf("CreateProtypeExport|Couldn't create list.json\n%s", err)
	}
	cmd := exec.Command("./script/protype.py", tmp+"/list.json", "Standard")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("CreateProtypeExport|Couldn't run python script\n%s\n\n%s", output, err)
	}
	args = append(args, "Standard/wordlist.dat")
	for _, addonID := range aIDs {
		abbs, err := s.GetAbbs(addonID)
		if err != nil {
			return fmt.Errorf("CreateProtypeExport|Couldn't get addon abbs:\n%s", err)
		}
		list, err := s.GetList(addonID)
		if err != nil {
			return fmt.Errorf("CreateProtypeExport|Couldn't get addon list:\n%s", err)
		}
		bytes, err := json.Marshal(abbs)
		if err != nil {
			return fmt.Errorf("CreateProtypeExport|Couldn't marshall addon list:\n%s", err)
		}
		err = ioutil.WriteFile(tmp+"/list.json", bytes, 0644)
		if err != nil {
			return fmt.Errorf("CreateProtypeExport|Couldn't create list.json:\n%s", err)
		}
		cmd := exec.Command("./script/protype.py", tmp+"/list.json", list.Name)
		_, err = cmd.Output()

		if err != nil {
			return fmt.Errorf("CreateProtypeExport|Couldn't run python script:\n%s", err)
		}
		args = append(args, list.Name+"/wordlist.dat")
	}
	args = append(args, "settings.dat")
	cmd = exec.Command("zip", args...)
	cmd.Dir = tmp
	stdin, err := cmd.StdinPipe()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	io.WriteString(stdin, "SkrivTolk\n.\n")
	if err != nil {
		return fmt.Errorf("createProtypeExport|Couldn't use stdin pipe:\n%s", err)
	}
	defer stdin.Close()
	if err = cmd.Start(); err != nil {
		return fmt.Errorf("CreateProtypeExport|Couldn't zip files:\n%s", err)
	}
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("CreateProtypeExport|Couldn't quit zip:\n%s", err)
	}
	return nil
}

func (s *abbService) CreateToTExport(sID string, aIDs []string) (textOnTopExport, error) {
	var tot textOnTopExport
	if sID != "" {
		list := make(abbMap)
		standard, err := s.GetAbbs(sID)
		if err != nil {
			return tot, fmt.Errorf("CreateToTExport couldn't get standard abbs:\n%s", err.Error())
		}
		for _, abb := range standard {
			list[abb.Abb] = abb.Word
		}
		tot.Autocorrect.Default.List = list
	}
	if aIDs != nil {
		tot.Shortform = make(map[string]map[string]abbMap)
		log.Printf("\t")
		for _, addon := range aIDs {
			meta, err := s.GetList(addon)
			tot.Shortform[meta.Name] = make(map[string]abbMap)
			if err != nil {
				return tot, fmt.Errorf("CreateToTExport couldn't get lists:\n%s", err.Error())
			}
			tot.Shortform[meta.Name]["shortforms"] = make(abbMap)
			abbs, err := s.GetAbbs(addon)
			if err != nil {
				return tot, fmt.Errorf("CreateToTExport couldn't get addon abbs:\n%s", err.Error())
			}
			for _, abb := range abbs {
				tot.Shortform[meta.Name]["shortforms"][abb.Abb] = abb.Word

			}
		}
	}
	return tot, nil

}
