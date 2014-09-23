package tunnel

import (
	"testing"
)

var c *Client = &Client{}
var a *Command = newCommand(c, "Test", map[string]string{"a": "b"})

type testSt struct {
	name string
	age  int
}

var cl, _ = createClient()
var b *Command = newCommand(c, "Demo", testSt{"vt", 23})

func TestCommandGenerate(t *testing.T) {
	if len(a.Id) != 16 {
		t.Fatalf("Command generate with wrong ID: %s", a.String())
	}
	t.Log(a.String())
}

func TestCommandParse(t *testing.T) {
	cl.Register("Demo", func() interface{} {
		return testSt{}
	}, func(args interface{}) (string, error) {
		return args.(testSt).name, nil
	})
	_b := receiveCommand(cl, b.String())
	r, err := _b.Handler(_b.Args)
	if err != nil && r != "vt" {
		t.Log(_b)
		t.Fatal("fail to parse command")
	}
}
