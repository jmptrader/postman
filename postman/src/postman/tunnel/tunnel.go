package tunnel

type Tunnel interface {
	// auth to remove server
	Auth(str string)
	// set client authenticated
	SetAuthenticated()
	Request(action string, args interface{}) string
	// send request and will not return util get response from remote
	RequestBlock(action string, args interface{}) (string, error)
	Register(action string, instance func() interface{}, handler func(interface{}))
	Serve()
	// exit and close connection
	Close()
}

// create new tunnel
func New(config Config) Tunnel {
	return &Proto{
		config:          config,
		RequestChan:     make(chan interface{}, 12),
		actionMap:       map[string]*Action{},
		requestBlockMap: map[string]chan string{},
	}
}
