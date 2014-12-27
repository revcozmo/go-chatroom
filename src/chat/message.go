// Package message provides ...
package chat

import (
	"bytes"
)

const (
	_ = iota
	NORMAL
	DISCONNECT
	SETUP
	QUIT
	JOIN
	PAUSE
	KICK
	CANCEL
)

type Message struct {
	Command  uint32
	RoomName string
	Content  []byte
}

type RoomMessage struct {
	Name   string
	Client *Client
}

func toLittleUint32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3]<<24)
}

func NewMessage(line []byte) Message {

	spliter := []byte(" ")
	data := bytes.SplitN(line[4:], spliter, 2)
	return Message{toLittleUint32(line[:4]),
		string(data[0]),
		data[1]}
}
