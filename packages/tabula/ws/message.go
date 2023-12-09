package ws

import (
	"github.com/fmpwizard/go-quilljs-delta/delta"
	"github.com/qwertyist/tabula/collab"
)

type TXMessage struct {
	ID       string       `json:"id,omitempty"`
	Type     PoolMessage  `json:"type"`
	Chat     *ChatMessage `json:"chat,omitempty"`
	Msg      string       `json:"msg,omitempty"`
	Data     int          `json:"msg,omitempy`
	Password string       `json:"password,omitempty"`
	Abb      *SharedAbb   `json:"abb,omitempty"`
	Body     *struct {
		Version int         `json:"version"`
		Delta   delta.Delta `json:"delta,omitempty"`
		Index   int         `json:"index"`
	} `json:"body,omitempty"`
	Zoom    *collab.ZoomCC `json:"zoom,omitempty"`
	Clients []Client       `json:"clients,omitempty"`
}
