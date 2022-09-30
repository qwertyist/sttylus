package collab

import (
	"fmt"

	"github.com/fmpwizard/go-quilljs-delta/delta"
)

type Tabula struct {
	Version int
	Doc     *delta.Delta
	Ops     map[int]delta.Op
}

type Delta struct {
	Version int
	Delta   *delta.Delta
	Index   int
}

func NewTabula(d Delta) *Tabula {
	return &Tabula{
		Version: d.Version,
		Doc:     d.Delta,
		Ops:     make(map[int]delta.Op),
	}
}

func (t *Tabula) ClearHandler() error {
	t.Version = 0
	t.Doc = delta.New(nil)
	t.Ops = make(map[int]delta.Op)
	return nil
}

func (t *Tabula) DeltaHandler(d Delta) (Delta, error) {
	index := 0
	for i, op := range d.Delta.Ops {
		t.Ops[d.Version+i] = op
		index += op.Length()
	}
	//	log.Println("Index of delta:", index)
	//	log.Println("Tabula version:", t.Version)
	t.Version += len(d.Delta.Ops)
	t.Doc = t.Doc.Compose(*d.Delta)
	return Delta{t.Version, nil, index}, nil
}

func (t *Tabula) RetrieveDoc() Delta {
	fmt.Println(t.Version)
	return Delta{
		Version: t.Version,
		Delta:   t.Doc,
		Index:   t.Doc.Length(),
	}
}

func (t *Tabula) RetrieveOp(version int) (delta.Op, error) {
	op, ok := t.Ops[version]
	if ok {
		return op, nil
	}
	return delta.Op{}, fmt.Errorf("Op not found")
}
