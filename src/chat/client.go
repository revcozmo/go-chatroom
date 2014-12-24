// Package client provides ...
package chat

import (
	"bufio"
	"log"
	"net"
)

type Client struct {
	Server *ChatServer
	Id     uint32
	Conn   net.Conn
	Rooms  map[uint32]*Room
	In     chan Message
	Out    chan []byte
}

func NewClient(server *ChatServer, id uint32, conn net.Conn) *Client {
	return &Client{server, id,
		conn,
		make(map[uint32]*Room),
		make(chan Message, 255),
		make(chan []byte, 255)}
}

func (c *Client) Listen() {
	log.Printf("New client: %d", c.Id)
	for {
		select {
		case msg := <-c.In:
			// Client receive message
			c.Write(msg)
		case msg := <-c.Out:
			c.ParseAndSend(msg)
		}
	}
}

func (c *Client) Write(msg Message) {

	c.Conn.Write(msg.Content)
}

func (c *Client) ParseAndSend(line []byte) {

	msg := NewMessage(line)

	switch msg.Command {
	case SETUP:
		c.Server.NewRoom <- msg.RoomId
	case CANCEL:
		c.Server.CancelRoom <- msg.RoomId
	case JOIN:
		c.Server.JoinRoom <- *&JQMessage{msg.RoomId, c}
	case QUIT:
		c.Server.QuitRoom <- *&JQMessage{msg.RoomId, c}
	case NORMAL:
		c.Server.In <- msg
	}

}

func (c *Client) RecvFromConn() {

	scanner := bufio.NewScanner(c.Conn)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		c.Out <- scanner.Bytes()
	}
	if err := scanner.Err(); err != nil {
		log.Printf("%s", err)
	}

}
