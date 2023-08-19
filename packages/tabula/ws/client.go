package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/fmpwizard/go-quilljs-delta/delta"
	"github.com/gorilla/websocket"
	"github.com/qwertyist/tabula/collab"
)

type Client struct {
	ID          string
	Interpreter bool
	Conn        *websocket.Conn
	Pool        *Pool
	mu          sync.Mutex
}

type Message struct {
	Type     PoolMessage `json:"type"`
	Msg      string      `json:"msg,omitempty"`
	Password string      `json:"password,omitempty"`
	Abb      *SharedAbb  `json:"abb,omitempty"`
	Body     struct {
		Version int         `json:"version"`
		Delta   delta.Delta `json:"delta,omitempty"`
		Index   int         `json:"index"`
	} `json:"body,omitempty"`
	Zoom collab.ZoomCC `json:"zoom,omitempty"`
}

type Broadcast struct {
	Message *Message
	Conn    *websocket.Conn
}

type PoolMessage int

const (
	OK            PoolMessage = iota
	CreateSession             = 1
	JoinSession               = 2
	LeaveSession              = 3
	GetInfo                   = 4
	SetInfo                   = 5
	SetPassword               = 6
	GetPassword               = 7
	NotAuthorized             = 401
	NoSession                 = 404
	TXDelta                   = 20
	RXDelta                   = 21
	TXClear                   = 22
	RXClear                   = 23
	TXAbb                     = 24
	RXAbb                     = 25
	TXManuscript              = 26
	RXManuscript              = 27
	ReadySignal               = 38
	RetrieveDoc               = 30
	ZoomCC                    = 200111
	Ping                      = 200
	Pong                      = 300
	Loss                      = 500
)

func (c *Client) messageHandler(msg Message) (*Message, bool) {
	id := c.Pool.ID
	//log.Println("Message type:", msg.Type)
	switch msg.Type {
	case CreateSession:
		c.Pool.Tabula = collab.NewTabula(id, collab.Delta{Version: msg.Body.Version, Delta: &msg.Body.Delta})
		c.Pool.Password = msg.Password
		log.Println("Session created with password:", msg.Password)
		if msg.Msg == "started" {
			c.Pool.Started = true
		}
		return nil, false
	case JoinSession:
		log.Println("JoinSession:", msg)
		if !c.Interpreter && c.Pool.Password != msg.Password {
			return &Message{Type: NotAuthorized}, false
		} else {
			p := Message{Type: SetPassword, Msg: c.Pool.Password}
			c.send(p)
		}
		if c.Pool.Tabula != nil {
			log.Println("Joining existing Tabula")
			m := Message{Type: RetrieveDoc}
			d := c.Pool.Tabula.RetrieveDoc()
			log.Printf("RetrieveDoc: %+v\n", d)
			m.Body.Delta = *d.Delta
			m.Body.Version = d.Version
			m.Body.Index = d.Index
			if c.Pool.Started {
				m.Msg = "started"
			} else {
				m.Msg = "waiting"
			}
			c.send(m)
		} else {
			log.Println("No session exists")
			m := Message{Type: NoSession}
			c.send(m)
		}
		return &msg, true
	case LeaveSession:
		log.Println("LeaveSession:", msg)
		return &msg, true
	case GetInfo:
		log.Println("Info:", msg)
		return &msg, true
	case SetPassword:
		log.Println("SetPassword:", msg.Password)
		c.Pool.Password = msg.Password
		return nil, false
	case SetInfo:
		log.Println("SetInfo:", msg.Zoom)
		err := c.Pool.Tabula.SetZoomData(msg.Zoom)
		if err != nil {
			msg.Msg = err.Error()
			msg.Zoom.MainStep = -1
			c.send(msg)
			return nil, false
		}
		return &msg, true
	case TXDelta:
		//		log.Printf("TXDelta: (version %d) %v", msg.Body.Version, msg.Body.Delta)
		if c.Pool.Tabula != nil {
			if !c.Pool.Started {
				c.Pool.Started = true
				msg.Msg = "starting"
			}
			d := collab.Delta{
				Version: msg.Body.Version,
				Delta:   &msg.Body.Delta,
			}
			u, _ := c.Pool.Tabula.DeltaHandler(d)
			msg.Body.Version = u.Version
			msg.Body.Index = u.Index
			msg.Type = RXDelta
			//c.Pool.Tabula.ToText()
			return &msg, true
		} else {
			return nil, false
		}
	case TXClear:
		c.Pool.Tabula.ClearHandler()
		msg.Type = RXClear
		return &msg, true
	case RXDelta:
		return nil, false
	case TXAbb:
		msg.Type = RXAbb
		return &msg, true
	case ReadySignal:
		return &msg, true
	case RetrieveDoc:
		return nil, false
	case ZoomCC:
		if c.Pool.Tabula.Zoom.Token != "" {
			c.Pool.Tabula.SendZoomCC(msg.Msg)
		}
		return nil, false
	case Ping:
		msg = Message{Type: Pong}
		return &msg, false
	case Pong:
		return nil, false
	}
	return &msg, true
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if messageType >= 1000 {
			log.Println("Closing somehow...")
			return
		}

		if err != nil {
			log.Println(err)
			return
		}
		var msg Message
		err = json.Unmarshal(p, &msg)
		if err != nil {
			log.Println("first:", err)
			log.Println("failed message is:", string(p))
		}
		/*
			if delta.Version == 0 {
				log.Println("New session")
				tabula := collab.NewTabula(delta.Delta)
				c.Pool.Tabula = tabula
			} else {
				delta, err := c.Pool.Tabula.DeltaHandler(delta)
				body, err = json.Marshal(delta)
				if err != nil {
					log.Println(err)
				}
			}
		*/
		handledMsg, send := c.messageHandler(msg)
		if handledMsg != nil {
			if send {
				broadcast := Broadcast{Conn: c.Conn, Message: handledMsg}
				c.Pool.Broadcast <- broadcast
			} else {
				c.send(*handledMsg)
			}
		} else {
			continue
		}
	}
}
