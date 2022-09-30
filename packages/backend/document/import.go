package document

import (
	"log"
	"os"
	"regexp"

	"github.com/lu4p/cat"
)

func (s *docService) ImportDoc(file *os.File) (string, error) {
	doc, err := cat.File(file.Name())
	if err != nil {
		return doc, err
	}
	re := regexp.MustCompile(`\r?\n`)
	parsed := re.ReplaceAllString(doc, "<br />")
	log.Println(parsed)
	return parsed, err
}
