package ws

func (c *Client) send(msg Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Conn.WriteJSON(msg)
}
