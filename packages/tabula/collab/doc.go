package collab

import "log"

func (t *Tabula) ToText() error {
	if t.Doc.Ops != nil {
		log.Printf("%+v\n", string(t.Doc.Ops[0].Insert))
	}
	return nil
}
