package wsmessage

import (
	"slices"
	"testing"
)

func TestMarshal(t *testing.T) {
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
