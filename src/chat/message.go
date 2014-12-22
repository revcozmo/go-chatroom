// Package message provides ...
package chat

import (
	"time"
)

const (
	_ = iota
	NORMAL
	DISCONNECT
	SETUP
	QUIT
	JOIN
	DISMISS
	PAUSE
	KICK
)

type Message struct {
	Sender   *Client
	Receiver string
	Command  int
	Content  []byte
	Time     time.Time
}

func NewMessage(sender *Client, receiver string, command int,
	content []byte) *Message {
	return &Message{sender, receiver, command, content, time.Now()}
}
