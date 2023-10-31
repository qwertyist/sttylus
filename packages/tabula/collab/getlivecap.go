package collab

import (
	"encoding/json"
	"strings"
	"unicode"
	"unicode/utf8"
)

type GETlivecapXML struct {
}

type JSONcaption struct {
	Line1 string
	Line2 string
}

type GETlivecapOptions struct {
}

func wordWrap(text string, lineWidth int) string {
	wrap := make([]byte, 0, len(text)+2*len(text)/lineWidth)
	eoLine := lineWidth
	inWord := false
	for i, j := 0, 0; ; {
		r, size := utf8.DecodeRuneInString(text[i:])
		if size == 0 && r == utf8.RuneError {
			r = ' '
		}
		if unicode.IsSpace(r) {
			if inWord {
				if i >= eoLine {
					wrap = append(wrap, '\n')
					eoLine = len(wrap) + lineWidth
				} else if len(wrap) > 0 {
					wrap = append(wrap, ' ')
				}
				wrap = append(wrap, text[j:i]...)
			}
			inWord = false
		} else if !inWord {
			inWord = true
			j = i
		}
		if size == 0 && r == ' ' {
			break
		}
		i += size
	}
	return string(wrap)
}

func getCaptionBlock(text string, c, n int, rolling bool) []string {
	rolling = true
	c = 37
	n = 2
	text = wordWrap(text, c)
	lines := strings.Split(text, "\n")
	switch len(lines) {
	case 0:
		return []string{"", ""}
	case 1:
		return []string{lines[0], ""}
	case 2:
		return []string{lines[0], lines[1]}
	}
	return []string{lines[len(lines)-2], lines[len(lines)-1]}
}

func createJSONcaption(lines []string) ([]byte, error) {
	caption := JSONcaption{Line1: lines[0], Line2: lines[1]}
	bytes, err := json.Marshal(caption)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (t *Tabula) GETlivecap(options *GETlivecapOptions) ([]byte, error) {
	if options == nil {
		text := strings.TrimSuffix(t.ToText(), "\n")
		block := []string{}
		lines := strings.Split(text, "\n")
		block = getCaptionBlock(lines[len(lines)-1], 37, 2, true)
		bytes, err := createJSONcaption(block)
		if err != nil {
			return nil, err
		}
		return bytes, nil
	}
	return nil, nil
}
