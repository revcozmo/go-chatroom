// Package main provides ...
package chat

import (
	"bytes"
	"encoding/binary"
	"log"
)

type Room struct {
	Server  *ChatServer
	Name    string
	Clients map[string]*Client
	In      chan []byte
}

func NewRoom(server *ChatServer, name string) *Room {

	room := &Room{server, name, make(map[string]*Client, 0),
		make(chan []byte, 256)}
	go room.Listen()
	return room
}

func (r *Room) Listen() {

	log.Printf("Chatroom: %s opened", r.Name)
	for {
		select {
		case msg := <-r.In:
			cmd, _ := binary.Varint(msg)
			switch cmd {
			default:
				r.broadcast(msg)
			case JOIN:
				data := bytes.SplitN(msg, []byte(" "), 3)
				name := string(data[2])
				if _, ok := r.Server.Clients[name]; ok {
					r.Clients[name] = r.Server.Clients[name]
					r.broadcast(msg)
				}
			}
		}
	}
}

func (r *Room) broadcast(msg []byte) {

	for _, c := range r.Clients {
		c.In <- msg
	}

}
