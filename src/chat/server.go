// Package server provides ...
package chat

import (
	"log"
	"math/rand"
	"net"
)

type ChatServer struct {
	Bind_to string
	Hub     *Hub
}

func NewChatServer(bind_to string) *ChatServer {

	return &ChatServer{
		bind_to,
		NewHub(),
	}
}

func (server *ChatServer) ListenAndServe() {

	listener, err := net.Listen("tcp", server.Bind_to)

	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	go server.Hub.Serve()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go func() {
			c := NewClient(server.Hub, string(rand.Uint32()), conn)
			go c.RecvFromConn()
			go c.Listen()
		}()
	}
}
