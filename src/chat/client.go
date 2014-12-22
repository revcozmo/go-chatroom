// Package client provides ...
package chat

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

type Client struct {
	Server *ChatServer
	Name   string
	Conn   net.Conn
	Rooms  map[string]*Room
	In     chan *Message
	Out    chan *Message
}

func NewClient(ser *ChatServer, name string, conn net.Conn) *Client {
	return &Client{ser,
		name,
		conn,
		make(map[string]*Room),
		make(chan *Message),
		make(chan *Message)}
}

func (c *Client) Listen() {
	log.Printf("New client: %s", c.Name)
	for {
		select {
		case msg := <-c.In:
			// Client receive message
			go c.Write(msg)
		case msg := <-c.Out:
			switch msg.Command {
			case DISCONNECT:
				// broadcast to all rooms
				for _, r := range c.Rooms {
					r.In <- msg
				}
			case JOIN:
				name := msg.Receiver
			default:
				c.Rooms[msg.Receiver].In <- msg
			}
		}
	}
}

func (c *Client) Write(msg *Message) {

	fmt.Fprintf(c.Conn,
		"%s %s:%s\n",
		msg.Time.Format(time.RFC3339),
		msg.Sender.Name,
		msg.Content)
}

func (c *Client) RecvFromConn() {

	var msg *Message
	spliter := []byte(" ")

	scanner := bufio.NewScanner(c.Conn)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {

		data := bytes.SplitN(scanner.Bytes(), spliter, 3)
		if len(data) != 2 {
			log.Printf("Not validated message %s from %s", data, c.Name)
			continue
		}

		c.In <- NewMessage(c, receiver, command, content)
	}
	if err := scanner.Err(); err != nil {
		log.Printf("%s", err)
	}

}
