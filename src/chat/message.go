// Package message provides ...
package chat

import (
	"bytes"
	"encoding/binary"
	"log"
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

type MessageHeader struct {
	Command uint32
	RoomId  uint32
}

type Message struct {
	Command uint32
	RoomId  uint32
	Content []byte
}

func NewMessage(line []byte) Message {
	var header MessageHeader

	err := binary.Read(bytes.NewBuffer(line[:8]),
		binary.LittleEndian,
		&header)

	if err != nil {
		log.Printf("send invalied msg:%v", err)
	}
	msg := Message{header.Command, header.RoomId, line[8:]}
	return msg
}
