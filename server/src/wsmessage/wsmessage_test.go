package wsmessage

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestUnmarshal(t *testing.T) {
	raw := []byte{11, 0, 0, 0, 1, 0, 2, 0, 23, 0, 42}

	msg, err := Unmarshal(raw)
	if err != nil {
		t.Error(err)
	}

	if msg.Code != 2 {
		t.Errorf("unexpected msg code %d", msg.Code)
	}

	if msg.Length != 11 {
		t.Errorf("unexpected msg length %d", msg.Length)
	}

	if msg.Version != version1 {
		t.Errorf("unexpected msg version %d", msg.Version)
	}

	if msg.Sender != 23 {
		t.Errorf("unexpected sender %d", msg.Sender)
	}

	if !slices.Equal(msg.RawPayload, []byte{42}) {
		t.Errorf("unexpected payload %+v", msg.RawPayload)
	}
}

func TestMarshalNilPayload(t *testing.T) {
	raw, err := Marshal(2, 23, nil)
	if err != nil {
		t.Error(err)
	}

	if len(raw) != 10 {
		t.Errorf("unexpected len: %d", len(raw))
	}

	if !slices.Equal(raw, []byte{10, 0, 0, 0, 1, 0, 2, 0, 23, 0}) {
		t.Errorf("unexpected raw bytes %+v", raw)
	}
}

func TestMarshalPayload(t *testing.T) {
	raw, err := Marshal(3, 23, []byte{1, 2, 3, 4})
	if err != nil {
		t.Error(err)
	}

	if len(raw) != 14 {
		t.Errorf("unexpected len: %d", len(raw))
	}

	if !slices.Equal(raw, []byte{14, 0, 0, 0, 1, 0, 3, 0, 23, 0, 1, 2, 3, 4}) {
		t.Errorf("unexpected raw bytes %+v", raw)
	}
}

func TestParseMessageWithPayload(t *testing.T) {
	pbytes, err := Marshal(CodeWelcome, 23, []byte(`{"user_id": 23}`))
	if err != nil {
		t.Error(err)

	}
	outStruct := WelcomePayload{}

	err = ParseMessageWithPayload(websocket.TextMessage, pbytes, CodeWelcome, &outStruct)
	if err == nil || err.Error() != "unexpected message type" {
		t.Error("expected message type err")
	}

	err = ParseMessageWithPayload(websocket.BinaryMessage, []byte{0x23}, CodeGame, &outStruct)
	if err == nil || err.Error() != "malformed message" {
		t.Error("expected malformed err")
	}

	pBytesWrongVersion := make([]byte, len(pbytes))
	copy(pBytesWrongVersion, pbytes)
	pBytesWrongVersion[4] = 2
	err = ParseMessageWithPayload(websocket.BinaryMessage, pBytesWrongVersion, CodeGame, &outStruct)
	if err == nil || err.Error() != "unexpected message version" {
		t.Error("expected version error")
	}

	err = ParseMessageWithPayload(websocket.BinaryMessage, pbytes, CodeGame, &outStruct)
	if err == nil || err.Error() != "unexpected message type code" {
		t.Error("expected code err")
	}

	badPBytes, err := Marshal(CodeWelcome, 23, []byte(`{"badjson"}`))
	if err != nil {
		t.Error(err)
	}

	err = ParseMessageWithPayload(websocket.BinaryMessage, badPBytes, CodeWelcome, &outStruct)
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
	msg, _ := Marshal(CodeJoin, 23, []byte(`{"name": "ataboo", "room_code": ""}`))

	sb := strings.Builder{}
	for _, c := range msg {
		sb.WriteString(fmt.Sprintf("%x ", c))
	}

	t.Log(sb.String())

}
