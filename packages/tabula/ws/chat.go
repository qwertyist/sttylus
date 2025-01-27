package ws

import "fmt"

type ChatMessage struct {
	To          string `json:"to,omitempty"`
	Msg         string `json:"message,omitempty"`
	Name        string `json:"name,omitempty"`
	Interpreter bool   `json:"interpreter"`
}

func (c *Client) sendChat(msg *ChatMessage) {
	fmt.Println("sending chat msg:", msg.To, msg.Msg)
	m := Message{ID: c.ID, Type: RXChat, Chat: msg}
	for to := range c.Pool.Clients {
		if to == c {
			continue
		}
		if msg.To == to.ID {
			//m.Chat.Msg // += "to someone" + msg.To

			fmt.Println("chat - dm")
			to.send(m)
			continue
		}
		if msg.To == "users" {
			if to.Interpreter {
				continue
			}

			fmt.Println("chat - to all users")
			// m.Chat.Msg // += "to users"
			to.send(m)
			continue
		}
		if msg.To == "interpreters" {
			if !to.Interpreter {
				continue
			}
			// m.Chat.Msg // += "to interpreters"
			fmt.Println("chat - to all interpreters")
			to.send(m)
			continue
		}

		if msg.To == "" {
			// m.Chat.Msg // += "to everybody"
			to.send(m)
		}
	}
	c.send(m)
}
