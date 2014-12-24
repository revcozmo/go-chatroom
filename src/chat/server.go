// Package server provides ...
package chat

import (
	"log"
	"math/rand"
	"net"
)

type ChatServer struct {
	Bind_to    string
	Rooms      map[uint32]*Room
	NewRoom    chan uint32
	CancelRoom chan uint32
	JoinRoom   chan JQMessage
	QuitRoom   chan JQMessage
	In         chan Message
}

type JQMessage struct {
	RoomId uint32
	Client *Client
}

func NewChatServer(bind_to string) *ChatServer {

	return &ChatServer{
		bind_to,
		make(map[uint32]*Room),
		make(chan uint32),
		make(chan uint32),
		make(chan JQMessage),
		make(chan JQMessage),
		make(chan Message),
	}
}

func (server *ChatServer) HandleMessage() {

	log.Printf("Start handle Message...")

	for {
		select {
		case rid := <-server.NewRoom:
			server.Rooms[rid] = NewRoom(rid)
		case rid := <-server.CancelRoom:
			delete(server.Rooms, rid)

		case msg := <-server.JoinRoom:
			if _, ok := server.Rooms[msg.RoomId]; ok {
				server.Rooms[msg.RoomId].Register <- msg.Client
			}
		case msg := <-server.QuitRoom:
			if _, ok := server.Rooms[msg.RoomId]; ok {
				server.Rooms[msg.RoomId].UnRegister <- msg.Client
			}
		case msg := <-server.In:
			if _, ok := server.Rooms[msg.RoomId]; ok {
				server.Rooms[msg.RoomId].In <- msg
			}
		}
	}
}

func (server *ChatServer) ListenAndServe() {

	listener, err := net.Listen("tcp", server.Bind_to)

	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	go server.HandleMessage()

	// Main loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go func() {
			c := NewClient(server, rand.Uint32(), conn)
			go c.Listen()
			go c.RecvFromConn()
		}()
	}
}
