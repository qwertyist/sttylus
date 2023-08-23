package ws

import (
	"fmt"
	"log"

	"github.com/jaevor/go-nanoid"
	"github.com/qwertyist/tabula/collab"
)

var Pools map[string]*Pool

type Pool struct {
	ID         string
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Broadcast
	Tabula     *collab.Tabula
	Password   string
	Started    bool
}


func NewPool(id string) *Pool {
	return &Pool{
		ID:         id,
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Broadcast),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			id, _ := nanoid.CustomASCII("abcdef1234567890", 8)
			client.ID = id()

			fmt.Println("Size of connection pool: ", len(pool.Clients))
			for c := range pool.Clients {
				log.Println("Client:", c.ID)
				msg := Message{Type: JoinSession}
				if client.Interpreter {
					msg.Msg = "interpreter"
				} else {
					msg.Msg = "user"
				}
				msg.ID = client.ID
				if c == client {
					msg.Msg = ""
				}
				c.send(msg)
			}
			break
		case client := <-pool.Unregister:
			interpreter := client.Interpreter
			delete(pool.Clients, client)
			fmt.Println("Size of connection pool:", len(pool.Clients))
			for c := range pool.Clients {
				msg := Message{Type: LeaveSession}
				msg.ID = client.ID
				if interpreter {
					msg.Msg = "interpreter"
				} else {
					msg.Msg = "user"
				}
				c.send(msg)
			}
			break
		case broadcast := <-pool.Broadcast:
			for client := range pool.Clients {
				if broadcast.Conn != client.Conn {

					if err := client.send(*broadcast.Message); err != nil {
						fmt.Println(err)
						return
					}
				}
			}
			break
		}
	}
}
