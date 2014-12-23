// Package client provides ...
package chat

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

type Client struct {
	Server *ChatServer
	Name   string
	Conn   net.Conn
	Rooms  map[string]*Room
	In     chan []byte
	Out    chan []byte
}

func NewClient(ser *ChatServer, name string, conn net.Conn) *Client {
	return &Client{ser,
		name,
		conn,
		make(map[string]*Room),
		make(chan []byte, 255),
		make(chan []byte, 255)}
}

func (c *Client) Listen() {
	log.Printf("New client: %s", c.Name)
	for {
		select {
		case msg := <-c.In:
			// Client receive message
			c.Write(msg)
		case msg := <-c.Out:
			go c.ParseAndSend(msg)
		}
	}
}

func (c *Client) Write(msg []byte) {
	c.Conn.Write(msg)
}

func (c *Client) ParseAndSend(msg []byte) {

	spliter := []byte(" ")

	data := bytes.SplitN(msg, spliter, 3)
	if len(data) != 3 {
		log.Printf(c.Name, "send invalied msg")
		return
	}
	cmd, _ := binary.Varint(data[0])
	receiver := string(data[1])

	switch cmd {
	case JOIN:
		c.Rooms[receiver] = c.Server.Rooms[receiver]
		c.Rooms[receiver].In <- msg

	default:
		if _, ok := c.Rooms[receiver]; ok {
			c.Rooms[receiver].In <- msg
		}
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
