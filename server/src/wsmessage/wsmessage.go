package wsmessage

import (
	"encoding/binary"
	"fmt"
)

const version1 = uint16(1)

const (
	CodeSetName    = 1
	CodeCreateRoom = 2
	CodeJoinRoom   = 3
	CodeLeaveRoom  = 4
	CodeQuit       = 5
	CodeBuzz       = 6
	CodeReset      = 7
)

type WSMessage struct {
	Length     uint32 // 4 bytes unsigned
	Version    uint16 // 2 bytes unsigned
	Code       uint16 // 2 bytes unsigned
	RawPayload []byte // ? bytes
}

func Unmarshal(raw []byte) (*WSMessage, error) {
	if len(raw) < 8 {
		return nil, fmt.Errorf("invalid format")
	}

	msg := WSMessage{
		Length:     binary.LittleEndian.Uint32(raw[0:4]),
		Version:    binary.LittleEndian.Uint16(raw[4:6]),
		Code:       binary.LittleEndian.Uint16(raw[6:8]),
		RawPayload: nil,
	}

	if int(msg.Length) != len(raw) {
		return nil, fmt.Errorf("invalid length")
	}

	if msg.Length > 8 {
		msg.RawPayload = raw[8:]
	}

	return &msg, nil
}

func Marshal(code uint16, rawPayload []byte) ([]byte, error) {
	bufferLen := 8
	if rawPayload != nil {
		bufferLen += len(rawPayload)
	}

	buffer := make([]byte, 0, bufferLen)
	buffer = binary.LittleEndian.AppendUint32(buffer, uint32(bufferLen))
	buffer = binary.LittleEndian.AppendUint16(buffer, version1)
	buffer = binary.LittleEndian.AppendUint16(buffer, code)

	if rawPayload != nil {
		copy(buffer[8:], rawPayload)
	}

	return buffer, nil
}
