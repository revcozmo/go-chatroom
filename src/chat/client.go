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

	In         chan Message
	BufferList chan []byte
}

func NewClient(h *Hub, name string, conn net.Conn) *Client {
	return &Client{h,
		name,
		conn,
		make(chan Message, 256),
		make(chan []byte, 256)}
}

func (c *Client) msgToByte(msg Message) []byte {

	s := make([]byte, 29+len(msg.Content))
	copy(s, []byte(time.Now().String()))
	copy(s[29:], msg.Content)
	return s
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
		c.ParseAndSend(scanner.Bytes())
	}
	if err := scanner.Err(); err != nil {
		log.Printf("%s", err)
		c.Hub.Quit <- c
	}

}

func (c *Client) Listen() {
	for {
		select {
		case msg := <-c.In:
			c.Conn.Write(c.msgToByte(msg))
		}
	}
}
