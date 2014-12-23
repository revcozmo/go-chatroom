// Package server provides ...
package chat

import (
	"fmt"
	"log"
	"net"
	"time"
)

type ChatServer struct {
	Bind_to string
	Rooms   map[string]*Room
	Clients map[string]*Client
}

func (server *ChatServer) reportStatus() {

	for {
		time.Sleep(10 * time.Second)
		for _, room := range server.Rooms {
			log.Printf("%s:%d", room.Name, len(room.Clients))
		}
	}

}

func (server *ChatServer) ListenAndServe() {

	listener, err := net.Listen("tcp", server.Bind_to)

	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	go server.reportStatus()
	// Main loop
	server.Rooms["WORLD"] = NewRoom(server, "WORLD")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go func(conn net.Conn, server *ChatServer) {
			c := NewClient(server, fmt.Sprintf("%s", conn.RemoteAddr()),
				conn)
			server.Clients[c.Name] = c
			go c.Listen()
			go c.RecvFromConn()
		}(conn, server)
	}

}
