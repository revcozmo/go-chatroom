// Package client provides ...
package chat

import (
	"bufio"
	"log"
	"net"
	"time"
)

type Client struct {
	Hub  *Hub
	Name string
	Conn net.Conn

	In chan Message
}

func NewClient(h *Hub, name string, conn net.Conn) *Client {
	return &Client{h,
		name,
		conn,
		make(chan Message)}
}

func (c *Client) Write(msg Message) {

	s := make([]byte, 29+len(msg.Content))
	copy(s, []byte(time.Now().String()))
	copy(s[29:], msg.Content)
	c.Conn.Write(s)
}

func (c *Client) ParseAndSend(line []byte) {

	msg := NewMessage(line)
	switch msg.Command {
	case NORMAL:
		c.Hub.Normal <- msg
	case SETUP:
		c.Hub.NewRoom <- *&RoomMessage{msg.RoomName, c}
	case JOIN:
		c.Hub.Register <- *&RoomMessage{msg.RoomName, c}
	}

}

func (c *Client) RecvFromConn() {

	scanner := bufio.NewScanner(c.Conn)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		go c.ParseAndSend(scanner.Bytes())
	}
	if err := scanner.Err(); err != nil {
		log.Printf("%s", err)
		c.Hub.Quit <- c
	}

}

func (c *Client) Listen() {
	for msg := range c.In {
		c.Write(msg)
	}
}
