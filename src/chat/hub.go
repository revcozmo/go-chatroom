package chat

import (
	"log"
	"sync"
)

type Hub struct {
	RoomNotify   *sync.Cond
	ClientNotify *sync.Cond
	Lock         sync.RWMutex

	Rooms       map[string]*Room
	NewRoom     chan RoomMessage
	DestoryRoom chan RoomMessage
	Register    chan RoomMessage
	UnRegister  chan RoomMessage

	Normal chan Message
	Quit   chan *Client
}

func NewHub() *Hub {
	rl := new(sync.Mutex)
	cl := new(sync.Mutex)

	return &Hub{
		sync.NewCond(rl),
		sync.NewCond(cl),
		sync.RWMutex{},

		make(map[string]*Room),
		make(chan RoomMessage),
		make(chan RoomMessage),
		make(chan RoomMessage),
		make(chan RoomMessage),

		make(chan Message),
		make(chan *Client),
	}
}

//Safe register Room
func (h *Hub) register(name string, c *Client) {

	if _, ok := h.Rooms[name]; ok {
		h.Rooms[name].Clients[c.Name] = c
		log.Printf("%s has %d clients", name, len(h.Rooms[name].Clients))
	}
}

//Safe unregister Room
func (h *Hub) unregister(name string, c *Client) {

	if _, ok := h.Rooms[name].Clients[c.Name]; ok {
		delete(h.Rooms[name].Clients, c.Name)
	}
}

//Safe new Room
func (h *Hub) newRoom(name string) {

	if _, ok := h.Rooms[name]; !ok {
		log.Printf("New Room: %s", name)
		room := NewRoom(name)
		h.Rooms[name] = room
	}
}

//Safe destory Room
func (h *Hub) destoryRoom(name string) {

	delete(h.Rooms, name)
}

func (h *Hub) quit(c *Client) {

	for _, r := range h.Rooms {
		h.unregister(r.Name, c)
	}
}

func (h *Hub) broadcast(msg Message) {

	h.Lock.Lock()
	defer h.Lock.Unlock()

	for _, c := range h.Rooms[msg.RoomName].Clients {
		c.In <- msg
	}
}

func (h *Hub) Serve() {
	// we only run this loop in one hub
	// which lock-free and no data race
	for {
		select {
		case msg := <-h.NewRoom:
			h.newRoom(msg.Name)
		case msg := <-h.DestoryRoom:
			h.destoryRoom(msg.Name)
		case msg := <-h.Register:
			h.register(msg.Name, msg.Client)
		case msg := <-h.UnRegister:
			h.unregister(msg.Name, msg.Client)
		case msg := <-h.Normal:
			h.broadcast(msg)
		case c := <-h.Quit:
			h.quit(c)
		}

	}
}
