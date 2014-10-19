package tunnel

import (
	"testing"
)

var c = &Proto{}
var a = newCommand(c, "Test", map[string]string{"a": "b"})

type testSt struct {
	name string
	age  int
}

var cl = &Proto{}
var b = newCommand(c, "Demo", testSt{"vt", 23})

func TestCommandGenerate(t *testing.T) {
	if len(a.Id) != 5 {
		t.Fatalf("Command generate with wrong ID: %s", a.String())
	}
}
