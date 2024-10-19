package wsmessage

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

const (
	Version1          uint16 = 1
	CurrentMsgVersion        = Version1
	ServerSenderId    uint16 = 0
	CodeWelcome       uint16 = 1
	CodeJoin          uint16 = 2
	CodeGame          uint16 = 3
)

type WSMessage struct {
	//Length     uint32 // 4 bytes unsigned
	Version    uint16 // 2 bytes unsigned
	Code       uint16 // 2 bytes unsigned
	Sender     uint16 // 2 bytes unsigned
	RawPayload []byte // ? bytes
}

const MsgHeaderLen uint32 = 4 + 2 + 2 + 2

func NewWsMessage(code uint16, sender uint16, payload interface{}) (*WSMessage, error) {
	rawPayload := []byte{}
	if payload != nil {
		pldBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		rawPayload = pldBytes
	}

	totalLen := int(MsgHeaderLen) + len(rawPayload)
	if totalLen >= 1<<32 {
		return nil, fmt.Errorf("total length too long")
	}

	return &WSMessage{
		Version:    CurrentMsgVersion,
		Code:       code,
		Sender:     sender,
		RawPayload: rawPayload,
	}, nil
}

func Unmarshal(raw []byte) (*WSMessage, error) {
	if len(raw) >= 1<<32 {
		return nil, fmt.Errorf("raw bytes too long")
	}

	if uint32(len(raw)) < MsgHeaderLen {
		return nil, fmt.Errorf("invalid format")
	}

	msgLen := binary.LittleEndian.Uint32(raw[0:4])

	msg := WSMessage{
		Version:    binary.LittleEndian.Uint16(raw[4:6]),
		Code:       binary.LittleEndian.Uint16(raw[6:8]),
		Sender:     binary.LittleEndian.Uint16(raw[8:10]),
		RawPayload: nil,
	}

	if int(msgLen) != len(raw) {
		return nil, fmt.Errorf("invalid length")
	}

	if msgLen > MsgHeaderLen {
		msg.RawPayload = raw[MsgHeaderLen:]
	}

	return &msg, nil
}

func Marshal(msg *WSMessage) ([]byte, error) {
	bufferLen := MsgHeaderLen
	if msg.RawPayload != nil {
		bufferLen += uint32(len(msg.RawPayload))
	}

	buffer := make([]byte, 0, bufferLen)
	buffer = binary.LittleEndian.AppendUint32(buffer, uint32(bufferLen))
	buffer = binary.LittleEndian.AppendUint16(buffer, CurrentMsgVersion)
	buffer = binary.LittleEndian.AppendUint16(buffer, msg.Code)
	buffer = binary.LittleEndian.AppendUint16(buffer, msg.Sender)

	if msg.RawPayload != nil {
		buffer = append(buffer, msg.RawPayload...)
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
