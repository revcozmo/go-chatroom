// Package main provides ...
package chat

import (
	"log"
)

type Room struct {
	Id         uint32
	Clients    map[uint32]*Client
	In         chan Message
	Register   chan *Client
	UnRegister chan *Client
}

func NewRoom(id uint32) *Room {

	room := &Room{id, make(map[uint32]*Client, 0),
		make(chan Message),
		make(chan *Client),
		make(chan *Client),
	}
	go room.Listen()
	return room
}

func (r *Room) Listen() {

	log.Printf("Chatroom: %d opened", r.Id)
	for {
		select {
		case msg := <-r.In:
			r.broadcast(msg)
		case c := <-r.Register:
			r.Clients[c.Id] = c
		case c := <-r.UnRegister:
			delete(r.Clients, c.Id)
		}
	}
}

func (r *Room) broadcast(msg Message) {

	for _, c := range r.Clients {
		c.In <- msg
	}

}
