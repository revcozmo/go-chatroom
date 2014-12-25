// Package main provides ...
package chat

type Room struct {
	Name    string
	Clients map[string]*Client
}

func NewRoom(name string) *Room {

	room := &Room{name, make(map[string]*Client, 0)}
	return room
}
