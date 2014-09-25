package tunnel

import (
	"postman/store"
)

type Tunnel interface {
	Register(action string, instance func() interface{}, handler func(interface{}))
	Serve()
}

type tunnel struct {
	config Config
	client Client
}

type Config struct {
}

func New(config Config) Tunnel {
	t := &tunnel{
		config: config,
	}
	t.client = createClient()
	return t
}

func (t *tunnel) Register(action string, instance func() interface{}, handler func(interface{})) {
	t.client.Register(action, instance, handler)
}

func (t *tunnel) Serve() {

}
