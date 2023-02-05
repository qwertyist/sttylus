package collab

import "fmt"

func (t *Tabula) ToText() string {
	if t.Doc.Ops != nil {
		return fmt.Sprintf("%+v", string(t.Doc.Ops[0].Insert))
	}
	return ""
}
