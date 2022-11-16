package ws

import (
	"fmt"

	"github.com/qwertyist/tabula/collab"
)

var Pools map[string]*Pool

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Broadcast
	Tabula     *collab.Tabula
	Started    bool
}

func NewPool() *Pool {
	return &Pool{
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
			fmt.Println("Size of connection pool: ", len(pool.Clients))
			for c := range pool.Clients {
				if c == client {
					continue
				}
				msg := Message{Type: JoinSession}
				msg.Msg = "Client connected"
				c.Conn.WriteJSON(msg)
			}
			break
		case client := <-pool.Unregister:
			interpreter := client.Interpreter
			delete(pool.Clients, client)
			fmt.Println("Size of connection pool:", len(pool.Clients))
			for client := range pool.Clients {
				msg := Message{Type: LeaveSession}
				if interpreter {
					msg.Msg = "interpreter"
				} else {
					msg.Msg = "user"
				}
				client.Conn.WriteJSON(msg)
			}
			break
		case broadcast := <-pool.Broadcast:
			for client := range pool.Clients {
				if broadcast.Conn != client.Conn {

					client.mu.Lock()

					if err := client.Conn.WriteJSON(broadcast.Message); err != nil {
						fmt.Println(err)
						client.mu.Unlock()
						return
					}
					client.mu.Unlock()
				}
			}
			break
		}
	}
}
