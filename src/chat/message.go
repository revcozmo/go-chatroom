// Package message provides ...
package chat

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
	Command uint32
	RoomId  uint32
	Content []byte
}

func toLittleUint32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3]<<24)
}

func NewMessage(line []byte) Message {

	return Message{toLittleUint32(line[:4]),
		toLittleUint32(line[4:8]),
		line[8:]}
}
