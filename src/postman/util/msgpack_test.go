package util

import (
	"testing"
)

func TestEncode(t *testing.T) {
	a := map[string]string{"a": "b"}
	b := map[string]string{}
	msg, _ := MsgEncode(a)
	MsgDecode(msg, b)
	t.Log(b)
}
