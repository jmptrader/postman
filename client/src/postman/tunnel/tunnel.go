package tunnel

type Tunnel interface {
	Register(action string, instance func() interface{}, handler func(*Client, interface{}))
	Serve()
}

// create new tunnel
func New(config Config) Tunnel {
	return &Client{
		config:      config,
		RequestChan: make(chan string, 12),
		actionMap:   map[string]*Action{},
	}
}
