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
	Name        string
	Interpreter bool
	Conn        *websocket.Conn
	Pool        *Pool
	mu          sync.Mutex
}

type Message struct {
	ID       string       `json:"id,omitempty"`
	Type     PoolMessage  `json:"type"`
	Chat     *ChatMessage `json:"chat,omitempty"`
	Msg      string       `json:"msg,omitempty"`
	Data     int          `json:"data,omitempty"`
	Password string       `json:"password,omitempty"`
	Abb      *SharedAbb   `json:"abb,omitempty"`
	Body     struct {
		Version int         `json:"version"`
		Delta   delta.Delta `json:"delta,omitempty"`
		Index   int         `json:"index"`
	} `json:"body,omitempty"`
	Zoom    collab.ZoomCC `json:"zoom,omitempty"`
	Clients []Client      `json:"clients,omitempty"`
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
	GetClients                = 8
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
	TXChat                    = 40
	RXChat                    = 41
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
    log.Printf("%s: Session created with password: %s\n", id, msg.Password)
		if msg.Msg == "started" {
			c.Pool.Started = true
		}
		return nil, false
	case JoinSession:
		log.Println("JoinSession", msg.ID, msg)
		if msg.ID != "" {
			c.ID = msg.ID
		}
		if !c.Interpreter && c.Pool.Password != msg.Password {
			return &Message{Type: NotAuthorized}, false
		} else {
			p := Message{Type: SetPassword, Msg: c.Pool.Password}
			c.send(p)
		}
		if c.Pool.Tabula != nil {
			m := Message{Type: RetrieveDoc}
			d := c.Pool.Tabula.RetrieveDoc()
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
			log.Println("No session exists", id)
			m := Message{Type: NoSession}
			c.send(m)
			return nil, false
		}
		return &msg, true
	case LeaveSession:
    log.Printf("%s: LeaveSession: %s\n", id, msg)
		return &msg, true
	case GetInfo:
    log.Printf("%s: Info: %s\n", id, msg)
		return &msg, true
	case SetPassword:
    log.Printf("%s: SetPassword: %s\n", id, msg)
		c.Pool.Password = msg.Password
		return nil, false
	case SetInfo:
		err := c.Pool.Tabula.SetZoomData(msg.Zoom)
		if err != nil {
      log.Printf("%s: SetZoomData Err: %s\n", id, err.Error())
			msg.Msg = err.Error()
			msg.Zoom.MainStep = -1
			c.send(msg)
			return nil, false
		}
    log.Printf("%s: SetZoomData OK\n", id)
		c.send(msg)
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
	case GetClients:
		msg.Data = len(c.Pool.Clients)
		c.send(msg)
		return &msg, true
	case TXChat:
		c.Name = msg.Chat.Name
		c.sendChat(msg.Chat)
		return nil, false
	case ZoomCC:
		if c.Pool.Tabula.Zoom.Token != "" {
			c.Pool.Tabula.SendZoomCC(msg.Msg)
		}
		return nil, false
	case Ping:
		pong := TXMessage{Type: Pong}
		c.sendTX(pong)
		return nil, false
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
			log.Println("Closing somehow...", c.Pool.ID)
			return
		}

		if err != nil {
			log.Println(c.Pool.ID, err)
			return
		}
		var msg Message
		err = json.Unmarshal(p, &msg)
		if err != nil {
      log.Println(c.Pool.ID)
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
