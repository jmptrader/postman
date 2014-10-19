package tunnel

type Tunnel interface {
	Register(action string, instance func() interface{}, handler func(*Client, interface{}))
	Serve()
}

// create new tunnel
func New(config Config) Tunnel {
	return &Client{
		Config:          config,
		RequestChan:     make(chan interface{}, 12),
		actionMap:       map[string]*Action{},
		requestBlockMap: map[string]chan string{},
	}
}
