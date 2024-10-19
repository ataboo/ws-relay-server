package wsmessage

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

const version1 = uint16(1)

const CurrentMsgVersion = version1

const (
	CodeWelcome = 1
	CodeJoin    = 2
	CodeGame    = 3
)

type WSMessage struct {
	Length     uint32 // 4 bytes unsigned
	Version    uint16 // 2 bytes unsigned
	Code       uint16 // 2 bytes unsigned
	Sender     uint16 // 2 bytes unsigned
	RawPayload []byte // ? bytes
}

func Unmarshal(raw []byte) (*WSMessage, error) {
	if len(raw) < 10 {
		return nil, fmt.Errorf("invalid format")
	}

	msg := WSMessage{
		Length:     binary.LittleEndian.Uint32(raw[0:4]),
		Version:    binary.LittleEndian.Uint16(raw[4:6]),
		Code:       binary.LittleEndian.Uint16(raw[6:8]),
		Sender:     binary.LittleEndian.Uint16(raw[8:10]),
		RawPayload: nil,
	}

	if int(msg.Length) != len(raw) {
		return nil, fmt.Errorf("invalid length")
	}

	if msg.Length > 8 {
		msg.RawPayload = raw[10:]
	}

	return &msg, nil
}

func Marshal(code uint16, sender uint16, rawPayload []byte) ([]byte, error) {
	bufferLen := 10
	if rawPayload != nil {
		bufferLen += len(rawPayload)
	}

	buffer := make([]byte, 0, bufferLen)
	buffer = binary.LittleEndian.AppendUint32(buffer, uint32(bufferLen))
	buffer = binary.LittleEndian.AppendUint16(buffer, CurrentMsgVersion)
	buffer = binary.LittleEndian.AppendUint16(buffer, code)
	buffer = binary.LittleEndian.AppendUint16(buffer, sender)

	if rawPayload != nil {
		buffer = append(buffer, rawPayload...)
	}

	return buffer, nil
}

func ParseMessageWithPayload(mType int, p []byte, expectedCode uint16, payloadStruct interface{}) error {
	if mType != websocket.BinaryMessage {
		return fmt.Errorf("unexpected message type")
	}

	msg, err := Unmarshal(p)
	if err != nil {
		return fmt.Errorf("malformed message")
	}

	if msg.Version != CurrentMsgVersion {
		return fmt.Errorf("unexpected message version")
	}

	if msg.Code != expectedCode {
		return fmt.Errorf("unexpected message type code")
	}

	payloadErr := json.Unmarshal(msg.RawPayload, payloadStruct)
	if payloadErr != nil {
		return fmt.Errorf("failed to parse payload")
	}

	return nil
}
