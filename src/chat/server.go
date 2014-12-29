// Package server provides ...
package chat

import (
	"fmt"
	"log"
	"net"
	"os"
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
			f, _ := os.Open("/dev/urandom")
			b := make([]byte, 16)
			f.Read(b)
			f.Close()
			c := NewClient(server.Hub, fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), conn)
			go c.RecvFromConn()
			go c.Listen()
		}()
	}
}
