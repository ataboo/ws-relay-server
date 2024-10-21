package wsmessage

import (
	"os"
	"slices"
	"testing"

	"github.com/gorilla/websocket"
)

func TestUnmarshal(t *testing.T) {
	raw := []byte{13, 0, 0, 0, 1, 0, 2, 0, 23, 0, 24, 0, 42}

	msg, err := Unmarshal(raw)
	if err != nil {
		t.Error(err)
	}

	if msg.Code != 2 {
		t.Errorf("unexpected msg code %d", msg.Code)
	}

	if msg.Version != Version1 {
		t.Errorf("unexpected msg version %d", msg.Version)
	}

	if msg.Sender != 23 {
		t.Errorf("unexpected sender %d", msg.Sender)
	}

	if msg.PayloadId != 24 {
		t.Errorf("unexpected payload id %d", msg.PayloadId)
	}

	if !slices.Equal(msg.RawPayload, []byte{42}) {
		t.Errorf("unexpected payload %+v", msg.RawPayload)
	}
}

func TestMarshalNilPayload(t *testing.T) {
	msg, err := NewWsMessage(2, 23, 24, nil)
	if err != nil {
		t.Error(err)
	}

	raw, err := Marshal(msg)
	if err != nil {
		t.Error(err)
	}

	if len(raw) != 12 {
		t.Errorf("unexpected len: %d", len(raw))
	}

	if !slices.Equal(raw, []byte{12, 0, 0, 0, 1, 0, 2, 0, 23, 0, 24, 0}) {
		t.Errorf("unexpected raw bytes %+v", raw)
	}
}

func TestMarshalPayload(t *testing.T) {
	msg := WSMessage{
		Version:    Version1,
		Code:       42,
		Sender:     23,
		PayloadId:  24,
		RawPayload: []byte{1, 2, 3},
	}

	raw, err := Marshal(&msg)
	if err != nil {
		t.Error(err)
	}

	if len(raw) != 15 {
		t.Errorf("unexpected len: %d", len(raw))
	}

	if !slices.Equal(raw, []byte{15, 0, 0, 0, 1, 0, 42, 0, 23, 0, 24, 0, 1, 2, 3}) {
		t.Errorf("unexpected raw bytes %+v", raw)
	}
}

func TestParseMessageWithPayload(t *testing.T) {
	msg, err := NewWsMessage(CodeWelcome, 23, PldIdWelcome, WelcomePayload{UserId: 23})
	if err != nil {
		t.Error(err)
	}

	pbytes, err := Marshal(msg)
	if err != nil {
		t.Error(err)

	}
	outStruct := WelcomePayload{}

	err = ParseMessageWithPayload(websocket.TextMessage, pbytes, CodeWelcome, &outStruct)
	if err == nil || err.Error() != "unexpected message type" {
		t.Error("expected message type err")
	}

	err = ParseMessageWithPayload(websocket.BinaryMessage, []byte{0x23}, CodeBroadcast, &outStruct)
	if err == nil || err.Error() != "malformed message: invalid format" {
		t.Error("expected malformed err", err.Error())
	}

	pBytesWrongVersion := make([]byte, len(pbytes))
	copy(pBytesWrongVersion, pbytes)
	pBytesWrongVersion[4] = 2
	err = ParseMessageWithPayload(websocket.BinaryMessage, pBytesWrongVersion, CodeBroadcast, &outStruct)
	if err == nil || err.Error() != "unexpected message version" {
		t.Error("expected version error")
	}

	err = ParseMessageWithPayload(websocket.BinaryMessage, pbytes, CodeBroadcast, &outStruct)
	if err == nil || err.Error() != "unexpected message type code" {
		t.Error("expected code err", err)
	}

	msg, err = NewWsMessage(3, 23, PldIdWelcome, []byte(`{"badjson"}`))
	if err != nil {
		t.Error(err)
	}

	badPBytes, err := Marshal(msg)
	if err != nil {
		t.Error(err)
	}

	err = ParseMessageWithPayload(websocket.BinaryMessage, badPBytes, CodeBroadcast, &outStruct)
	if err == nil || err.Error() != "failed to parse payload" {
		t.Error("expected payload err")
	}

	err = ParseMessageWithPayload(websocket.BinaryMessage, pbytes, CodeWelcome, &outStruct)
	if err != nil {
		t.Error(err)
	}

	if outStruct.UserId != 23 {
		t.Error("unexpected payload")
	}
}

func TestFoo(t *testing.T) {
	msg, err := NewWsMessage(CodeJoin, 23, PldIdJoin, &JoinPayload{Name: "ataboo", RoomCode: "MSFBEU"})
	if err != nil {
		t.Error(err)
	}

	msgBytes, err := Marshal(msg)
	if err != nil {
		t.Error(err)
	}

	err = os.WriteFile("msg_bytes.bin", msgBytes, 0777)
	if err != nil {
		t.Error(err)
	}
}
