package wsmessage

import (
	"slices"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	raw := []byte{9, 0, 0, 0, 1, 0, 2, 0, 42}

	msg, err := Unmarshal(raw)
	if err != nil {
		t.Error(err)
	}

	if msg.Code != 2 {
		t.Errorf("unexpected msg code %d", msg.Code)
	}

	if msg.Length != 9 {
		t.Errorf("unexpected msg length %d", msg.Length)
	}

	if msg.Version != version1 {
		t.Errorf("unexpected msg version %d", msg.Version)
	}

	if !slices.Equal(msg.RawPayload, []byte{42}) {
		t.Errorf("unexpected payload %+v", msg.RawPayload)
	}
}

func TestMarshalNilPayload(t *testing.T) {
	raw, err := Marshal(2, nil)
	if err != nil {
		t.Error(err)
	}

	if len(raw) != 8 {
		t.Errorf("unexpected len: %d", len(raw))
	}

	if !slices.Equal(raw, []byte{8, 0, 0, 0, 1, 0, 2, 0}) {
		t.Errorf("unexpected raw bytes %+v", raw)
	}
}

func TestMarshalPayload(t *testing.T) {
	raw, err := Marshal(3, []byte{1, 2, 3, 4})
	if err != nil {
		t.Error(err)
	}

	if len(raw) != 12 {
		t.Errorf("unexpected len: %d", len(raw))
	}

	if !slices.Equal(raw, []byte{12, 0, 0, 0, 1, 0, 3, 0, 1, 2, 3, 4}) {
		t.Errorf("unexpected raw bytes %+v", raw)
	}
}
